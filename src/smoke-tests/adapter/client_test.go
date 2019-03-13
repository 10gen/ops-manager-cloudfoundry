package client_test

import (
	"encoding/json"
	client "github.com/10gen/ops-manager-cloudfoundry/src/mongodb-service-adapter/adapter"
	"os"
	"testing"
)

var LatestVersion = "3.6.9"

func TestLoadDoc(t *testing.T) {
	t.Parallel()
	c := &client.OMClient{}

	for p, ctx := range map[string]*client.DocContext{
		client.PlanStandalone: {
			ID:            "d3a98gf1",
			Key:           "key",
			AdminPassword: "pwd",
			Nodes:         []string{"192.168.1.1", "192.168.1.2", "192.168.1.3", "192.168.1.4"},
			Version:       "3.2.11",
		},

		client.PlanShardedCluster: {
			ID:            "d3a98gf1",
			Key:           "key",
			AdminPassword: "pwd",
			Version:       "3.4.1",
			Cluster: &client.Cluster{
				Routers:       []string{"192.168.1.1", "192.168.1.2"},
				ConfigServers: []string{"192.168.1.3", "192.168.1.4"},
				Shards: [][]string{
					{"192.168.0.10", "192.168.0.11", "192.168.0.12"},
					{"192.168.0.13", "192.168.0.14", "192.168.0.15"},
				},
			},
		},

		client.PlanReplicaSet: {
			ID:            "d3a98gf1",
			Key:           "key",
			AdminPassword: "pwd",
			Nodes:         []string{"192.168.1.1", "192.168.1.2", "192.168.1.3", "192.168.1.4"},
			Version:       "3.2.11",
		},
	} {
		t.Run(string(p), func(t *testing.T) {
			s, err := c.LoadDoc(p, ctx)
			if err != nil {
				t.Fatal(err)
			}

			// validate json output
			if err := json.Unmarshal([]byte(s), &map[string]interface{}{}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetGroupByName(t *testing.T) {

	c := &client.OMClient{
		Url:      os.Getenv("Url"),
		Username: os.Getenv("Username"),
		ApiKey:   os.Getenv("ApiKey"),
	}

	_, err := c.GetGroupByName("test-group")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateGroup(t *testing.T) {
	t.Parallel()

	c := &client.OMClient{
		Url:      os.Getenv("Url"),
		Username: os.Getenv("Username"),
		ApiKey:   os.Getenv("ApiKey"),
	}

	request := client.GroupCreateRequest{
		Name: "test-group1",
	}

	_, err := c.CreateGroup(os.Getenv("GroupId")+"id", request)

	if err != nil {
		t.Fatal(err)
	}

}

func TestCreateGroupAPIKey(t *testing.T) {

	c := &client.OMClient{
		Url:      os.Getenv("Url"),
		Username: os.Getenv("Username"),
		ApiKey:   os.Getenv("ApiKey"),
	}

	_, err := c.CreateGroupAPIKey(os.Getenv("GroupId"))

	if err != nil {
		t.Fatal(err)
	}

}

func TestUpdateGroup(t *testing.T) {

	c := &client.OMClient{
		Url:      os.Getenv("Url"),
		Username: os.Getenv("Username"),
		ApiKey:   os.Getenv("ApiKey"),
	}

	request := client.GroupUpdateRequest{
		Tags: []string{
			"tag1",
			"tag2",
		},
	}

	_, err := c.UpdateGroup(os.Getenv("GroupId"), request)

	if err != nil {
		t.Fatal(err)
	}

}

func TestGetGroupAuthAgentPassword(t *testing.T) {

	c := &client.OMClient{
		Url:      os.Getenv("Url"),
		Username: os.Getenv("Username"),
		ApiKey:   os.Getenv("ApiKey"),
	}

	_, err := c.GetGroupAuthAgentPassword(os.Getenv("GroupId"))

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetGroup(t *testing.T) {

	c := &client.OMClient{
		Url:      os.Getenv("Url"),
		Username: os.Getenv("Username"),
		ApiKey:   os.Getenv("ApiKey"),
	}

	_, err := c.GetGroup("test-group-id")

	if err != nil {
		t.Fatal(err)
	}

}

func TestDeleteGroup(t *testing.T) {
	t.Parallel()

	c := &client.OMClient{
		Url:      os.Getenv("Url"),
		Username: os.Getenv("Username"),
		ApiKey:   os.Getenv("ApiKey"),
	}

	err := c.DeleteGroup("test-group-id")

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetGroupHosts(t *testing.T) {

	c := &client.OMClient{
		Url:      os.Getenv("Url"),
		Username: os.Getenv("Username"),
		ApiKey:   os.Getenv("ApiKey"),
	}

	_, err := c.GetGroupHosts(os.Getenv("GroupId"))

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetGroupHostnames(t *testing.T) {

	c := &client.OMClient{
		Url:      os.Getenv("Url"),
		Username: os.Getenv("Username"),
		ApiKey:   os.Getenv("ApiKey"),
	}

	_, err := c.GetGroupHostnames(os.Getenv("GroupId"), "test-plan-id")

	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigureGroup(t *testing.T) {
	t.Parallel()

	c := &client.OMClient{
		Url:      os.Getenv("Url"),
		Username: os.Getenv("Username"),
		ApiKey:   os.Getenv("ApiKey"),
	}

	err := c.ConfigureGroup("configure-doc", os.Getenv("GroupId"))

	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigureMonitoringAgent(t *testing.T) {
	t.Parallel()

	c := &client.OMClient{
		Url:      os.Getenv("Url"),
		Username: os.Getenv("Username"),
		ApiKey:   os.Getenv("ApiKey"),
	}

	err := c.ConfigureMonitoringAgent("configure-doc", os.Getenv("GroupId"))

	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigureBackupAgent(t *testing.T) {
	t.Parallel()

	c := &client.OMClient{
		Url:      os.Getenv("Url"),
		Username: os.Getenv("Username"),
		ApiKey:   os.Getenv("ApiKey"),
	}

	err := c.ConfigureBackupAgent("configure-doc", os.Getenv("GroupId"))

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAvailableVersions(t *testing.T) {
	t.Parallel()

	c := &client.OMClient{
		Url:      os.Getenv("Url"),
		Username: os.Getenv("Username"),
		ApiKey:   os.Getenv("ApiKey"),
	}

	_, err := c.GetAvailableVersions(os.Getenv("GroupId"))

	if err != nil {
		t.Fatal()
	}
}

func TestGetLatestVersion(t *testing.T) {
	t.Parallel()

	c := &client.OMClient{
		Url:      os.Getenv("Url"),
		Username: os.Getenv("Username"),
		ApiKey:   os.Getenv("ApiKey"),
	}

	latestVersion, err := c.GetLatestVersion(os.Getenv("GroupId"))

	if err != nil {
		t.Fatal()
	}

	LatestVersion = latestVersion
}

func TestValidateVersion(t *testing.T) {
	t.Parallel()

	c := &client.OMClient{
		Url:      os.Getenv("Url"),
		Username: os.Getenv("Username"),
		ApiKey:   os.Getenv("ApiKey"),
	}

	_, err := c.ValidateVersion(os.Getenv("GroupId"), LatestVersion)

	if err != nil {
		t.Fatal()
	}
}

func TestValidateVersionManifest(t *testing.T) {
	t.Parallel()

	c := &client.OMClient{
		Url:      os.Getenv("Url"),
		Username: os.Getenv("Username"),
		ApiKey:   os.Getenv("ApiKey"),
	}

	_, err := c.ValidateVersionManifest(LatestVersion)

	if err != nil {
		t.Fatal()
	}
}
