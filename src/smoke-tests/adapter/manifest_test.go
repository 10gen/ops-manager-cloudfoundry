package adapter_test

import (
	adapter "../../mongodb-service-adapter/adapter"
	bosh "github.com/pivotal-cf/on-demand-services-sdk/bosh"
	sa "github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	"os"
	"testing"
)

func TestGenerateManifest(t *testing.T) {

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
				"username":         os.Getenv("Username"),
				"api_key":          os.Getenv("ApiKey"),
				"bosh_dns_disable": true,
				"url":              os.Getenv("Url"),
				"backup_enabled":   true,
				"ssl_enabled":      false,
				"ssl_ca_cert":      "ca_cert",
				"ssl_pem":          "ssl_pem",
				"tags":             nil,
			},
			"syslog": map[string]interface{}{
				"address":        os.Getenv("address"),
				"port":           os.Getenv("port"),
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
			"project_name": "test-group",
			"orgId":        os.Getenv("OrgId"),
		},
	}

	previousManifest := bosh.BoshManifest{
		InstanceGroups: []bosh.InstanceGroup{
			{
				Properties: map[string]interface{}{
					"mongo_ops": map[interface{}]interface{}{
						"admin_password": os.Getenv("admin_password"),
						"id":             "standalone",
					},
				},
			},
			{
				Properties: map[string]interface{}{
					"mongo_ops": map[interface{}]interface{}{
						"admin_password": os.Getenv("admin_password"),
						"id":             "standalone",
						"group_id":       os.Getenv("GroupId"),
						"agent_api_key":  os.Getenv("ApiKey"),
						"auth_key":       os.Getenv("auth_key"),
					},
				},
			},
		},
	}

	previousPlan := sa.Plan{}
	mGenerator := &adapter.ManifestGenerator{}
	_, err := mGenerator.GenerateManifest(serviceDeployment, plan, requestParams, &previousManifest, &previousPlan)

	if err != nil {
		t.Fatal(err)
	}
}
