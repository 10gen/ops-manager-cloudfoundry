package adapter

import (
	"fmt"
	"log"
	"strings"

	"github.com/mongodb-labs/pcgc/pkg/opsmanager"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
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

	oc := &OMClient{URL: url, Username: username, APIKey: apiKey}

	var previousMongoProperties map[interface{}]interface{}

	if params.PreviousManifest != nil {
		previousMongoProperties = mongoPlanProperties(*params.PreviousManifest)
	}

	adminPassword, err := passwordForMongoServer(previousMongoProperties)
	if err != nil {
		return serviceadapter.GenerateManifestOutput{}, err
	}

	id, err := idForMongoServer(previousMongoProperties)
	if err != nil {
		return serviceadapter.GenerateManifestOutput{}, err
	}

	group, err := groupForMongoServer(id, oc, params.Plan.Properties, previousMongoProperties, arbitraryParams)
	if err != nil {
		return serviceadapter.GenerateManifestOutput{}, fmt.Errorf("could not create new group (%s)", err.Error())
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
		return serviceadapter.GenerateManifestOutput{}, fmt.Errorf("no definition found for instance group '%s'", MongodInstanceGroupName)
	}

	mongodJobs, err := gatherJobs(params.ServiceDeployment.Releases, []string{MongodJobName})
	if err != nil {
		return serviceadapter.GenerateManifestOutput{}, err
	}
	mongodJobs[0].AddSharedProvidesLink(MongodJobName)
	if syslogProps["address"].(string) != "" {
		mongodJobs, err = gatherJobs(params.ServiceDeployment.Releases, []string{MongodJobName, SyslogJobName})
		if err != nil {
			return serviceadapter.GenerateManifestOutput{}, err
		}
		mongodJobs[0].AddSharedProvidesLink(MongodJobName)
	}

	configAgentJobs, err := gatherJobs(params.ServiceDeployment.Releases, []string{ConfigAgentJobName, CleanupErrandJobName, PostSetupErrandJobName})
	if err != nil {
		return serviceadapter.GenerateManifestOutput{}, err
	}
	if syslogProps["address"].(string) != "" {
		configAgentJobs, err = gatherJobs(params.ServiceDeployment.Releases, []string{ConfigAgentJobName, CleanupErrandJobName, SyslogJobName, PostSetupErrandJobName})
		if err != nil {
			return serviceadapter.GenerateManifestOutput{}, err
		}
	}

	addonsJobs, err := gatherJobs(params.ServiceDeployment.Releases, []string{BoshDNSEnableJobName})
	if err != nil {
		return serviceadapter.GenerateManifestOutput{}, err
	}
	if boshDNSDisable {
		addonsJobs = []bosh.Job{}
	}

	mongodNetworks := []bosh.Network{}
	for _, network := range mongodInstanceGroup.Networks {
		mongodNetworks = append(mongodNetworks, bosh.Network{Name: network})
	}
	if len(mongodNetworks) == 0 {
		return serviceadapter.GenerateManifestOutput{}, fmt.Errorf("no networks definition found for instance group '%s'", MongodInstanceGroupName)
	}

	var engineVersion string
	version := getArbitraryParam("version", "engine_version", arbitraryParams, previousMongoProperties)
	if version == nil {
		engineVersion, err = oc.GetLatestVersion(group.ID)
		if err != nil {
			return serviceadapter.GenerateManifestOutput{}, fmt.Errorf("unable to find the latest MongoDB version from the MongoDB Ops Manager API. Please contact your system administrator to ensure versions are available in the Version Manager for project '%s' in MongoDB Ops Manager. If your MongoDB Ops Manager is running in Local Mode, then after validating versions are available, please indicate a specific MongoDB version using 'version’ paramater when calling 'create-service'", group.Name)
		}
	} else {
		skipVersionValidation := false
		e := getArbitraryParam("skip_version_check", "skip_version_check", arbitraryParams, previousMongoProperties)
		if e != nil {
			skipVersionValidation = e.(bool)
		}
		if !skipVersionValidation {
			engineVersion, err = oc.ValidateVersionManifest(version.(string))
			if err != nil {
				return serviceadapter.GenerateManifestOutput{}, err
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
		return serviceadapter.GenerateManifestOutput{}, fmt.Errorf("unknown plan: %s", planID)
	}
	authKey, err := authKeyForMongoServer(previousMongoProperties)
	if err != nil {
		return serviceadapter.GenerateManifestOutput{}, err
	}
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

	// TODO: what is this? there's no 'autoPwd' field in ProjectResponse!
	agentPassword := ""
	if params.PreviousManifest == nil {
		var err error
		agentPassword, err = GenerateString(32)
		if err != nil {
			panic(err)
		}
	}

	// if manifest updates we should use password from previous manifest
	if agentPassword == "" && params.PreviousManifest != nil {
		cfg, err := oc.Client().GetAutomationConfig(group.ID)
		if err != nil {
			panic(err)
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
			return nil, err
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

func passwordForMongoServer(previousManifestProperties map[interface{}]interface{}) (string, error) {
	if previousManifestProperties != nil {
		return previousManifestProperties["admin_password"].(string), nil
	}

	return GenerateString(20)
}

func idForMongoServer(previousManifestProperties map[interface{}]interface{}) (string, error) {
	if previousManifestProperties != nil {
		return previousManifestProperties["id"].(string), nil
	}

	return GenerateString(8)
}

func authKeyForMongoServer(previousManifestProperties map[interface{}]interface{}) (string, error) {
	if previousManifestProperties != nil {
		return previousManifestProperties["auth_key"].(string), nil
	}

	return GenerateString(8)
}

func groupForMongoServer(mongoID string, oc *OMClient,
	planProperties map[string]interface{},
	previousMongoProperties map[interface{}]interface{},
	arbitraryParams map[string]interface{}) (opsmanager.ProjectResponse, error) {
	req := GroupCreateRequest{}
	if name, found := arbitraryParams["projectName"]; found {
		req.Name = name.(string)
	}
	if orgID, found := arbitraryParams["orgId"]; found {
		req.OrgID = orgID.(string)
	}
	tags := planProperties["mongo_ops"].(map[string]interface{})["tags"]
	if tags != nil {
		t := tags.([]interface{})
		for _, tag := range t {
			req.Tags = append(req.Tags, tag.(map[string]interface{})["tag_name"].(string))
		}
	}

	if previousMongoProperties != nil {
		group, err := oc.Client().SetProjectTags(previousMongoProperties["group_id"].(string), req.Tags)
		if err != nil {
			return opsmanager.ProjectResponse{}, err
		}
		// AgentAPIKey is empty for PATCH and GET requests in OM 3.6, taking the value from previous manifest instead
		group.AgentAPIKey = previousMongoProperties["agent_api_key"].(string)
		return group, nil
	}

	return oc.CreateGroup(mongoID, req)
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
		return serviceadapter.ServiceRelease{}, fmt.Errorf("no release provided for job '%s'", requiredJob)
	}

	if len(releasesThatProvideRequiredJob) > 1 {
		releaseNames := []string{}
		for _, release := range releasesThatProvideRequiredJob {
			releaseNames = append(releaseNames, release.Name)
		}

		return serviceadapter.ServiceRelease{}, fmt.Errorf("job '%s' defined in multiple releases: %s", requiredJob, strings.Join(releaseNames, ", "))
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
