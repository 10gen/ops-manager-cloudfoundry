package adapter

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/mongodb-labs/pcgc/pkg/httpclient"
	"github.com/mongodb-labs/pcgc/pkg/opsmanager"
	"github.com/tidwall/gjson"
)

type OMClient struct {
	Url      string
	Username string
	ApiKey   string
}

type Automation struct {
	MongoDbVersions []MongoDbVersionsType
}

type MongoDbVersionsType struct {
	Name string
}

type GroupCreateRequest struct {
	Name  string   `json:"name"`
	OrgId string   `json:"orgId,omitempty"`
	Tags  []string `json:"tags"`
}

type GroupUpdateRequest struct {
	Tags []string `json:"tags"`
}

type GroupHosts struct {
	TotalCount int `json:"totalCount"`
}

var versionsManifest = []string{"/var/vcap/packages/versions/versions.json", "../../mongodb_versions/versions.json"}

func (oc *OMClient) Client() opsmanager.Client {
	r := httpclient.NewURLResolverWithPrefix(oc.Url, "api/public/v1.0")
	return opsmanager.NewClientWithDigestAuth(r, oc.Username, oc.ApiKey)
}

func (oc *OMClient) CreateGroup(id string, request GroupCreateRequest) (opsmanager.ProjectResponse, error) {
	log.Println(fmt.Sprintf("CreateGroup in id : %s ,request : %+v", id, request))

	if request.Name == "" {
		request.Name = fmt.Sprintf("PCF_%s", strings.ToUpper(id))
	}

	group, err := oc.Client().GetProjectByName(request.Name)
	if err != nil {
		log.Printf("CreateGroup GetGroupByName with request.Name: %s, error: %v", request.Name, err)
		return group, err
	}

	if group.Name == request.Name {
		log.Printf("Continue with existing group %q", group.ID)
		apiKey, err := oc.Client().CreateAgentAPIKEY(group.ID, "MongoDB On-Demand broker generated Agent API Key")
		if err != nil {
			log.Printf("CreateGroup CreateGroupAPIKey group.ID: %s, error: %v", group.ID, err)
			return group, err
		}
		group.AgentAPIKey = apiKey.Key
		return group, nil
	}

	resp, err := oc.Client().CreateOneProject(request.Name, request.OrgId)
	if err != nil {
		log.Printf("CreateGroup CreateOneProject, request: %+v, error: %v", request, err)
		return group, err
	}

	return resp, nil
}

func (oc *OMClient) GetLatestVersion(groupID string) (string, error) {
	cfg, err := oc.Client().GetAutomationConfig(groupID)
	if err != nil || len(cfg.MongoDBVersions) == 0 {
		return "", fmt.Errorf("unable to find the latest MongoDB version from the MongoDB Ops Manager API. Please contact your system administrator to ensure versions are available in the Version Manager for group '%q' in MongoDB Ops Manager. If your MongoDB Ops Manager is running in Local Mode, then after validating versions are available, please indicate a specific MongoDB version using 'version’ paramater when calling 'create-service'", groupID)
	}

	versions := make([]string, len(cfg.MongoDBVersions))
	n := 0
	for _, i := range cfg.MongoDBVersions {
		if !strings.HasSuffix(i.Name, "ent") {
			versions[n] = i.Name
			n++
		}
	}
	versions = versions[:n]
	latestVersion := versions[len(versions)-1]

	return latestVersion, nil
}

func (oc *OMClient) ValidateVersion(groupID string, version string) (string, error) {
	cfg, err := oc.Client().GetAutomationConfig(groupID)
	if err != nil {
		return "", err
	}

	for _, v := range cfg.MongoDBVersions {
		if v.Name == version {
			return version, nil
		}
	}

	return "", errors.New("failed to find expected version, got " + version)
}

func (oc *OMClient) ValidateVersionManifest(version string) (string, error) {
	b, err := ioutil.ReadFile(versionsManifest[0])
	if err != nil {
		b, err = ioutil.ReadFile(versionsManifest[1])
		if err != nil {
			return "", err
		}
	}

	v := gjson.GetBytes(b, fmt.Sprintf(`versions.#[name="%s"].name`, version))
	log.Printf("Using %q version of MongoDB", v.String())
	if v.String() == "" {
		return "", errors.New("failed to find expected version, continue with provided versions ")
	}

	return version, nil
}
