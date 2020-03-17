package adapter

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/mongodb-labs/pcgc/pkg/opsmanager"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	"github.com/pkg/errors"
)

const (
	StemcellAlias           = "mongodb-stemcell"
	MongodInstanceGroupName = "mongod_node"
	MongodJobName           = "mongod_node"
	AliasesJobName          = "mongodb-dns-aliases"
	SyslogJobName           = "syslog_forwarder"
	BoshDNSEnableJobName    = "bosh-dns-enable"
	ConfigAgentJobName      = "mongodb_config_agent"
	CleanupErrandJobName    = "cleanup_service"
	PostSetupErrandJobName  = "post_setup"
	LifecycleErrandType     = "errand"
)

type ManifestGenerator struct {
	Logger *log.Logger
}

func (m *ManifestGenerator) logf(msg string, v ...interface{}) {
	if m.Logger != nil {
		m.Logger.Printf(msg, v...)
	}
}

func (m ManifestGenerator) GenerateManifest(params serviceadapter.GenerateManifestParams) (serviceadapter.GenerateManifestOutput, error) {
	m.logf("request params: %#v", params.RequestParams)

	arbitraryParams := params.RequestParams.ArbitraryParams()

	mongoOps := params.Plan.Properties["mongo_ops"].(map[string]interface{})
	syslogProps := params.Plan.Properties["syslog"].(map[string]interface{})

	username := mongoOps["username"].(string)
	apiKey := mongoOps["api_key"].(string)
	boshDNSDisable := mongoOps["bosh_dns_disable"].(bool)

	// trim trailing slash
	url := mongoOps["url"].(string)
	url = strings.TrimRight(url, "/")

	oc := NewPCGCClient(url, username, apiKey)
	oc404 := NewPCGCClient(url, username, apiKey, http.StatusNotFound)

	var previousMongoProperties map[interface{}]interface{}

	if params.PreviousManifest != nil {
		previousMongoProperties = mongoPlanProperties(*params.PreviousManifest)
	}

	adminPassword := passwordForMongoServer(previousMongoProperties)
	id := idForMongoServer(previousMongoProperties)
	group, err := groupForMongoServer(id, oc, oc404, previousMongoProperties, arbitraryParams)
	if err != nil {
		return serviceadapter.GenerateManifestOutput{}, errors.Wrap(err, "cannot create new group")
	}
	m.logf("created group %s", group.ID)

	releases := []bosh.Release{}
	for _, release := range params.ServiceDeployment.Releases {
		releases = append(releases, bosh.Release{
			Name:    release.Name,
			Version: release.Version,
		})
	}

	mongodInstanceGroup := findInstanceGroup(params.Plan, MongodInstanceGroupName)
	if mongodInstanceGroup == nil {
		return serviceadapter.GenerateManifestOutput{}, errors.Errorf("no definition found for instance group %q", MongodInstanceGroupName)
	}

	mongodJobNames := []string{MongodJobName}
	configAgentJobNames := []string{ConfigAgentJobName, CleanupErrandJobName, PostSetupErrandJobName}
	addonsJobNames := []string{BoshDNSEnableJobName}

	if syslogProps["address"].(string) != "" {
		mongodJobNames = append(mongodJobNames, SyslogJobName)
		configAgentJobNames = append(configAgentJobNames, SyslogJobName)
	}

	mongodJobs, err := gatherJobs(params.ServiceDeployment.Releases, mongodJobNames)
	if err != nil {
		return serviceadapter.GenerateManifestOutput{}, errors.Wrapf(err, "cannot gather jobs %q", mongodJobNames)
	}
	mongodJobs[0].AddSharedProvidesLink(MongodJobName)

	configAgentJobs, err := gatherJobs(params.ServiceDeployment.Releases, configAgentJobNames)
	if err != nil {
		return serviceadapter.GenerateManifestOutput{}, errors.Wrapf(err, "cannot gather jobs %q", configAgentJobNames)
	}

	addonsJobs, err := gatherJobs(params.ServiceDeployment.Releases, addonsJobNames)
	if err != nil {
		return serviceadapter.GenerateManifestOutput{}, errors.Wrapf(err, "cannot gather jobs %q", addonsJobNames)
	}
	if boshDNSDisable {
		addonsJobs = []bosh.Job{}
	}

	mongodNetworks := []bosh.Network{}
	for _, network := range mongodInstanceGroup.Networks {
		mongodNetworks = append(mongodNetworks, bosh.Network{Name: network})
	}
	if len(mongodNetworks) == 0 {
		return serviceadapter.GenerateManifestOutput{}, errors.Errorf("no networks definition found for instance group %q", MongodInstanceGroupName)
	}

	var engineVersion string
	version := getArbitraryParam("version", "engine_version", arbitraryParams, previousMongoProperties)
	if version == nil {
		engineVersion, err = getLatestVersion(oc, group.ID)
		if err != nil {
			return serviceadapter.GenerateManifestOutput{},
				errors.Wrapf(
					err,
					"unable to find the latest MongoDB version from the MongoDB Ops Manager API. Please contact your system administrator to ensure versions are available in the Version Manager for project %q in MongoDB Ops Manager. If your MongoDB Ops Manager is running in Local Mode, then after validating versions are available, please indicate a specific MongoDB version using 'version’ paramater when calling 'create-service'",
					group.Name,
				)
		}
	} else {
		skipVersionValidation := false
		e := getArbitraryParam("skip_version_check", "skip_version_check", arbitraryParams, previousMongoProperties)
		if e != nil {
			skipVersionValidation = e.(bool)
		}
		if !skipVersionValidation {
			engineVersion, err = validateVersionManifest(version.(string))
			if err != nil {
				return serviceadapter.GenerateManifestOutput{}, errors.Wrap(err, "cannot validate version manifest")
			}
		} else {
			engineVersion = version.(string)
		}
	}

	// sharded_cluster parameters
	replicas := 0
	shards := 0
	routers := 0
	configServers := 0

	// total number of instances
	//
	// standalone:      always one
	// replica_set:     number of replicas
	// sharded_cluster: shards*replicas + config_servers + mongos
	instances := mongodInstanceGroup.Instances

	planID := params.Plan.Properties["id"].(string)
	switch planID {
	case PlanStandalone:
		// ok
	case PlanReplicaSet:
		r := getArbitraryParam("replicas", "replicas", arbitraryParams, previousMongoProperties)
		if r != nil {
			instances = r.(int)
		}
		replicas = instances
	case PlanShardedCluster:
		shards = 2
		s := getArbitraryParam("shards", "shards", arbitraryParams, previousMongoProperties)
		if s != nil {
			shards = s.(int)
		}

		replicas = 3
		r := getArbitraryParam("replicas", "replicas", arbitraryParams, previousMongoProperties)
		if r != nil {
			replicas = r.(int)
		}

		configServers = 3
		c := getArbitraryParam("config_servers", "config_servers", arbitraryParams, previousMongoProperties)
		if c != nil {
			configServers = c.(int)
		}

		routers = 2
		r = getArbitraryParam("mongos", "routers", arbitraryParams, previousMongoProperties)
		if r != nil {
			routers = r.(int)
		}

		instances = routers + configServers + shards*replicas
	default:
		return serviceadapter.GenerateManifestOutput{}, errors.Errorf("unknown plan: %q", planID)
	}

	authKey := authKeyForMongoServer(previousMongoProperties)
	backupEnabled := false
	if planID != PlanStandalone {
		e := getArbitraryParam("backup_enabled", "backup_enabled", arbitraryParams, previousMongoProperties)
		if e != nil {
			backupEnabled = e.(bool)
		} else {
			backupEnabled = mongoOps["backup_enabled"].(bool)
		}
	}
	requireSSL := false
	if !boshDNSDisable {
		e := getArbitraryParam("ssl_enabled", "ssl_enabled", arbitraryParams, previousMongoProperties)
		if e != nil {
			requireSSL = e.(bool)
		} else {
			requireSSL = mongoOps["ssl_enabled"].(bool)
		}
	}

	agentPassword := ""
	if params.PreviousManifest == nil {
		agentPassword = GenerateString(32, true)
	}

	// if manifest updates we should use password from previous manifest
	if agentPassword == "" && params.PreviousManifest != nil {
		cfg, err := oc.GetAutomationConfig(group.ID)
		if err != nil {
			return serviceadapter.GenerateManifestOutput{}, errors.Wrapf(err, "cannote get automation config for group %q", group.ID)
		}

		agentPassword = cfg.Auth.AutoPwd
	}

	caCert := ""
	if mongoOps["ssl_ca_cert"] != "" {
		caCert = mongoOps["ssl_ca_cert"].(string)
	}
	sslPem := ""
	if mongoOps["ssl_pem"] != "" {
		sslPem = mongoOps["ssl_pem"].(string)
	}

	updateBlock := &bosh.Update{
		Canaries:        1,
		CanaryWatchTime: "3000-180000",
		UpdateWatchTime: "3000-180000",
		MaxInFlight:     1,
	}

	manifest := bosh.BoshManifest{
		Name:     params.ServiceDeployment.DeploymentName,
		Releases: releases,
		Addons: []bosh.Addon{
			{
				Name: "mongodb-dns-helpers",
				Jobs: addonsJobs,
			},
		},
		Stemcells: []bosh.Stemcell{
			{
				Alias:   StemcellAlias,
				OS:      params.ServiceDeployment.Stemcells[0].OS,
				Version: params.ServiceDeployment.Stemcells[0].Version,
			},
		},
		InstanceGroups: []bosh.InstanceGroup{
			{
				Name:               MongodInstanceGroupName,
				Instances:          instances,
				Jobs:               mongodJobs,
				VMType:             mongodInstanceGroup.VMType,
				VMExtensions:       mongodInstanceGroup.VMExtensions,
				Stemcell:           StemcellAlias,
				PersistentDiskType: mongodInstanceGroup.PersistentDiskType,
				AZs:                mongodInstanceGroup.AZs,
				Networks:           mongodNetworks,
				Env: map[string]interface{}{
					"persistent_disk_fs": "xfs",
				},
				Properties: map[string]interface{}{
					"bpm": map[string]interface{}{
						"enabled": boolPointer(true),
					},
				},
			},
			{
				Name:         "mongodb-config-agent",
				Instances:    1,
				Jobs:         configAgentJobs,
				VMType:       mongodInstanceGroup.VMType,
				VMExtensions: mongodInstanceGroup.VMExtensions,
				Stemcell:     StemcellAlias,
				AZs:          mongodInstanceGroup.AZs,
				Networks:     mongodNetworks,

				// See mongodb_config_agent job spec
				Properties: map[string]interface{}{
					"bpm": map[string]interface{}{
						"enabled": boolPointer(true),
					},
					"mongo_ops": map[string]interface{}{
						"id":               id,
						"url":              url,
						"agent_api_key":    group.AgentAPIKey,
						"api_key":          apiKey,
						"auth_key":         authKey,
						"username":         username,
						"auth_pwd":         agentPassword,
						"group_id":         group.ID,
						"plan_id":          planID,
						"admin_password":   adminPassword,
						"engine_version":   engineVersion,
						"routers":          routers,
						"config_servers":   configServers,
						"replicas":         replicas,
						"shards":           shards,
						"backup_enabled":   backupEnabled,
						"require_ssl":      requireSSL,
						"ssl_ca_cert":      caCert,
						"ssl_pem":          sslPem,
						"bosh_dns_disable": boshDNSDisable,
					},
				},
			},
		},
		Update: updateBlock,
		Features: bosh.BoshFeatures{
			UseDNSAddresses: boolPointer(true),
		},
		Properties: map[string]interface{}{
			"mongo_ops": map[string]interface{}{
				"url":            url,
				"api_key":        group.AgentAPIKey,
				"admin_api_key":  apiKey,
				"group_id":       group.ID,
				"username":       username,
				"admin_password": adminPassword,
				"require_ssl":    requireSSL,

				// options needed for binding
				"plan_id":        planID,
				"routers":        routers,
				"config_servers": configServers,
				"replicas":       replicas,
			},
			"syslog": map[string]interface{}{
				"address":        syslogProps["address"],
				"port":           syslogProps["port"],
				"transport":      syslogProps["transport"],
				"tls_enabled":    syslogProps["tls_enabled"],
				"permitted_peer": syslogProps["permitted_peer"],
				"ca_cert":        syslogProps["ca_cert"],
			},
		},
	}

	return serviceadapter.GenerateManifestOutput{
		Manifest:          manifest,
		ODBManagedSecrets: serviceadapter.ODBManagedSecrets{},
	}, nil
}

func getLatestVersion(oc opsmanager.Client, groupID string) (string, error) {
	cfg, err := oc.GetAutomationConfig(groupID)
	if err != nil {
		return "", errors.Wrapf(err, "cannot get automation config for group %q", groupID)
	}
	if len(cfg.MongoDBVersions) == 0 {
		return "", errors.Errorf("no MongoDB versions for group %q", groupID)
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

func findInstanceGroup(plan serviceadapter.Plan, jobName string) *serviceadapter.InstanceGroup {
	for _, instanceGroup := range plan.InstanceGroups {
		if instanceGroup.Name == jobName {
			return &instanceGroup
		}
	}
	return nil
}

func gatherJobs(releases serviceadapter.ServiceReleases, requiredJobs []string) ([]bosh.Job, error) {
	jobs := []bosh.Job{}
	for _, requiredJob := range requiredJobs {
		release, err := findReleaseForJob(releases, requiredJob)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot find release for job %q", requiredJob)
		}

		job := bosh.Job{
			Name:     requiredJob,
			Release:  release.Name,
			Consumes: map[string]interface{}{},
			Provides: map[string]bosh.ProvidesLink{},
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func mongoPlanProperties(manifest bosh.BoshManifest) map[interface{}]interface{} {
	return manifest.InstanceGroups[1].Properties["mongo_ops"].(map[interface{}]interface{})
}

func passwordForMongoServer(previousManifestProperties map[interface{}]interface{}) string {
	if previousManifestProperties != nil {
		pass, ok := previousManifestProperties["admin_password"].(string)
		if ok {
			return pass
		}
	}

	return GenerateString(20, true)
}

func idForMongoServer(previousManifestProperties map[interface{}]interface{}) string {
	if previousManifestProperties != nil {
		id, ok := previousManifestProperties["id"].(string)
		if ok {
			return id
		}
	}

	return GenerateString(8, true)
}

func authKeyForMongoServer(previousManifestProperties map[interface{}]interface{}) string {
	if previousManifestProperties != nil {
		key, ok := previousManifestProperties["auth_key"].(string)
		if ok {
			return key
		}
	}

	return GenerateString(8, false)
}

func groupForMongoServer(
	mongoID string,
	oc opsmanager.Client,
	oc404 opsmanager.Client,
	previousMongoProperties map[interface{}]interface{},
	arbitraryParams map[string]interface{},
) (opsmanager.ProjectResponse, error) {
	name := ""
	orgID := ""
	if p, found := arbitraryParams["projectName"]; found {
		name, _ = p.(string)
	}

	if p, found := arbitraryParams["orgId"]; found {
		orgID, _ = p.(string)
	}

	// TODO: investigate the use of tags
	// tags := []string{}
	// if p := planProperties["mongo_ops"].(map[string]interface{})["tags"]; p != nil {
	// 	t := p.([]interface{})
	// 	for _, tag := range t {
	// 		tags = append(tags, tag.(map[string]interface{})["tag_name"].(string))
	// 	}
	// }

	if previousMongoProperties != nil {
		previousID, ok := previousMongoProperties["group_id"].(string)
		if !ok {
			return opsmanager.ProjectResponse{},
				errors.Errorf(
					"previous group_id is not a string: %q (%T)",
					previousMongoProperties["group_id"],
					previousMongoProperties["group_id"],
				)
		}

		group, err := oc.GetProjectByID(previousID)
		// AgentAPIKey is empty for PATCH and GET requests in OM 3.6, taking the value from previous manifest instead
		group.AgentAPIKey = previousMongoProperties["agent_api_key"].(string)
		return group, errors.Wrapf(err, "cannot get project %q", previousID)
	}

	// CreateGroup(mongoID, req)
	log.Println(fmt.Sprintf("CreateGroup: id %q, name %q", mongoID, name))

	if name == "" {
		name = fmt.Sprintf("PCF_%s", strings.ToUpper(mongoID))
	}

	group, err := oc404.GetProjectByName(name)
	if err != nil {
		log.Printf("CreateGroup GetGroupByName with name: %q, error: %v", name, err)
		return group, errors.Wrapf(err, "cannot get project by name %q", name)
	}

	if group.Name == name {
		log.Printf("Continue with existing group %q", group.ID)
		apiKey, err := oc.CreateAgentAPIKEY(group.ID, "MongoDB On-Demand broker generated Agent API Key")
		if err != nil {
			log.Printf("CreateGroup CreateGroupAPIKey group.ID: %s, error: %v", group.ID, err)
			return group, errors.Wrapf(err, "cannot create agent API key for group %q", name)
		}
		group.AgentAPIKey = apiKey.Key
		return group, nil
	}

	resp, err := oc.CreateOneProject(name, orgID)
	if err != nil {
		log.Printf("CreateGroup CreateOneProject: name %q, orgID %q error: %v", name, orgID, err)
		return group, errors.Wrapf(err, "cannot create project %q in org %q", name, orgID)
	}

	return resp, nil
}

func findReleaseForJob(releases serviceadapter.ServiceReleases, requiredJob string) (serviceadapter.ServiceRelease, error) {
	releasesThatProvideRequiredJob := serviceadapter.ServiceReleases{}

	for _, release := range releases {
		for _, providedJob := range release.Jobs {
			if providedJob == requiredJob {
				releasesThatProvideRequiredJob = append(releasesThatProvideRequiredJob, release)
			}
		}
	}

	if len(releasesThatProvideRequiredJob) == 0 {
		return serviceadapter.ServiceRelease{}, errors.Errorf("no release provided for job %q", requiredJob)
	}

	if len(releasesThatProvideRequiredJob) > 1 {
		releaseNames := []string{}
		for _, release := range releasesThatProvideRequiredJob {
			releaseNames = append(releaseNames, release.Name)
		}

		return serviceadapter.ServiceRelease{}, errors.Errorf("job %q defined in multiple releases: %s", requiredJob, strings.Join(releaseNames, ", "))
	}

	return releasesThatProvideRequiredJob[0], nil
}

func getArbitraryParam(propName string, manifestName string, arbitraryParams map[string]interface{}, previousMongoProperties map[interface{}]interface{}) interface{} {
	var prop interface{}
	var found bool
	if prop, found = arbitraryParams[propName]; found {
		goto found
	}
	if prop, found = previousMongoProperties[manifestName]; found {
		goto found
	}
	return nil
found:
	// we are interested only in string, bool and int properties, though json conversion returns float64 for all integer properties
	if p, ok := prop.(float64); ok {
		return int(p)
	}
	return prop
}

func boolPointer(b bool) *bool {
	return &b
}
