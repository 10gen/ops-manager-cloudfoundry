package adapter_test

import (
	adapter "../../mongodb-service-adapter/adapter"
	"fmt"
	bosh "github.com/pivotal-cf/on-demand-services-sdk/bosh"
	sa "github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	"testing"
)

func TestGenerateManifest(t *testing.T) {

	config, err := adapter.LoadConfig("../../mongodb-service-adapter/testdata/manifest.json")

	if err != nil {
		fmt.Print("Error opening manifest file ")
	}

	serviceDeployment := sa.ServiceDeployment{
		Releases: sa.ServiceReleases{
			sa.ServiceRelease{
				Name:    "standalone",
				Version: "4.0.6",
				Jobs: []string{adapter.MongodJobName, adapter.BPMJobName,
					adapter.SyslogJobName, adapter.ConfigAgentJobName,
					adapter.CleanupErrandJobName, adapter.ConfigureBackupsErrandJobName,
					adapter.BoshDNSEnableJobName},
			},
		},
		DeploymentName: "deploy_name",
		Stemcell: sa.Stemcell{
			OS:      "Ubuntu",
			Version: "16.X",
		},
	}

	plan := sa.Plan{
		Properties: sa.Properties{
			"mongo_ops": map[string]interface{}{
				"username":         config.Username,
				"api_key":          config.APIKey,
				"bosh_dns_disable": true,
				"url":              config.URL,
				"backup_enabled":   true,
				"ssl_enabled":      false,
				"ssl_ca_cert":      "ca_cert",
				"ssl_pem":          "ssl_pem",
				"tags":             nil,
			},
			"syslog": map[string]interface{}{
				"address":        config.NodeAddresses,
				"port":           "27017",
				"transport":      "tls",
				"tls_enabled":    false,
				"permitted_peer": 1,
				"ca_cert":        "ca_cert",
			},
			"id": "standalone",
		},
		InstanceGroups: []sa.InstanceGroup{
			{
				Name:     adapter.MongodJobName,
				Networks: []string{"network1"},
			},
		},
	}

	requestParams := sa.RequestParameters{
		"parameters": map[string]interface{}{
			"project_name": "Project 2",
			"orgId":        config.OrgID,
		},
	}

	previousManifest := bosh.BoshManifest{
		InstanceGroups: []bosh.InstanceGroup{
			{
				Properties: map[string]interface{}{
					"mongo_ops": map[interface{}]interface{}{
						"admin_password": "admin",
						"id":             "standalone",
					},
				},
			},
			{
				Properties: map[string]interface{}{
					"mongo_ops": map[interface{}]interface{}{
						"admin_password": "admin",
						"id":             "standalone",
						"group_id":       config.GroupID,
						"agent_api_key":  config.APIKey,
						"auth_key":       config.AuthKey,
					},
				},
			},
		},
	}

	previousPlan := sa.Plan{}
	mGenerator := &adapter.ManifestGenerator{}
	_, err = mGenerator.GenerateManifest(serviceDeployment, plan, requestParams, &previousManifest, &previousPlan)

	if err != nil {
		t.Fatal(err)
	}
}
