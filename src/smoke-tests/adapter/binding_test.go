package adapter_test

import (
	adapter "../../mongodb-service-adapter/adapter"
	"fmt"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	"testing"
)

func TestGetWithCredentials(t *testing.T) {

	config, err := adapter.LoadConfig("../../mongodb-service-adapter/testdata/manifest.json")

	if err != nil {
		fmt.Print("Error opening manifest file ")
	}

	addrs := []string{config.NodeAddresses}

	_, err = adapter.GetWithCredentials(addrs, "admin", false)

	if err != nil {
		t.Fatal(err)
	}

}

func TestCreateBinding(t *testing.T) {

	config, err := adapter.LoadConfig("../../mongodb-service-adapter/testdata/manifest.json")

	if err != nil {
		fmt.Print("Error opening manifest file ")
	}

	deploymentTopology := bosh.BoshVMs{
		"mongod_node": []string{
			config.NodeAddresses,
		},
	}

	manifest := bosh.BoshManifest{
		Properties: map[string]interface{}{
			"mongo_ops": map[interface{}]interface{}{
				"url":            config.URL,
				"group_id":       config.GroupID,
				"admin_password": "admin",
				"username":       config.Username,
				"admin_api_key":  config.APIKey,
				"require_ssl":    false,
				"plan_id":        "standalone",
				"routers":        0,
				"config_servers": 0,
				"replicas":       0,
			},
		},
	}

	requestParams := serviceadapter.RequestParameters{}
	secrets := serviceadapter.ManifestSecrets{}
	dnsAddresses := serviceadapter.DNSAddresses{}
	binder := &adapter.Binder{}

	_, err = binder.CreateBinding("binding1", deploymentTopology, manifest, requestParams, secrets, dnsAddresses)

	if err != nil {
		t.Fatal(err)
	}

}

func TestDeleteBinding(t *testing.T) {

	config, err := adapter.LoadConfig("../../mongodb-service-adapter/testdata/manifest.json")

	if err != nil {
		fmt.Print("Error opening manifest file ")
	}

	deploymentTopology := bosh.BoshVMs{
		"mongod_node": []string{
			config.NodeAddresses,
		},
	}

	manifest := bosh.BoshManifest{
		Properties: map[string]interface{}{
			"mongo_ops": map[interface{}]interface{}{
				"url":            config.URL,
				"group_id":       config.GroupID,
				"admin_password": "admin",
				"username":       config.Username,
				"admin_api_key":  config.APIKey,
				"require_ssl":    false,
				"plan_id":        "standalone",
				"routers":        0,
				"config_servers": 0,
				"replicas":       0,
			},
		},
	}

	requestParams := serviceadapter.RequestParameters{}
	secrets := serviceadapter.ManifestSecrets{}

	binder := &adapter.Binder{}

	err = binder.DeleteBinding("binding1", deploymentTopology, manifest, requestParams, secrets)

	if err != nil {
		t.Fatal(err)
	}

}
