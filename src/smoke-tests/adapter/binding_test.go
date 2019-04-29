package adapter_test

import (
	adapter "../../mongodb-service-adapter/adapter"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	"os"
	"testing"
)

func TestGetWithCredentials(t *testing.T) {
	addrs := []string{os.Getenv("Url")}

	_, err := adapter.GetWithCredentials(addrs, "admin", false)

	if err != nil {
		t.Fatal(err)
	}

	_, err = adapter.GetWithCredentials(addrs, "admin", true)

	if err != nil {
		t.Fatal(err)
	}

}

func TestCreateBinding(t *testing.T) {

	deploymentTopology := bosh.BoshVMs{
		"mongod_node": []string{
			os.Getenv("mongod_node"),
		},
	}

	manifest := bosh.BoshManifest{
		Properties: map[string]interface{}{
			"mongo_ops": map[interface{}]interface{}{
				"url":            os.Getenv("Url"),
				"group_id":       os.Getenv("GroupId"),
				"admin_password": "admin",
				"username":       os.Getenv("Username"),
				"admin_api_key":  os.Getenv("ApiKey"),
				"require_ssl":    true,
				"plan_id":        "sharded_cluster",
				"routers":        0,
				"config_servers": 1,
				"replicas":       1,
			},
		},
	}

	requestParams := serviceadapter.RequestParameters{}
	secrets := serviceadapter.ManifestSecrets{}
	dnsAddresses := serviceadapter.DNSAddresses{}
	binder := &adapter.Binder{}
	_, err := binder.CreateBinding("binding1", deploymentTopology, manifest, requestParams, secrets, dnsAddresses)

	if err != nil {
		t.Fatal(err)
	}

}

func TestDeleteBinding(t *testing.T) {

	deploymentTopology := bosh.BoshVMs{
		"mongod_node": []string{
			os.Getenv("mongod_node"),
		},
	}

	manifest := bosh.BoshManifest{
		Properties: map[string]interface{}{
			"mongo_ops": map[interface{}]interface{}{
				"url":            os.Getenv("Url"),
				"group_id":       os.Getenv("GroupId"),
				"admin_password": "admin",
				"username":       os.Getenv("Username"),
				"admin_api_key":  os.Getenv("ApiKey"),
				"require_ssl":    true,
				"plan_id":        "sharded_cluster",
				"routers":        0,
				"config_servers": 1,
				"replicas":       1,
			},
		},
	}

	requestParams := serviceadapter.RequestParameters{}
	secrets := serviceadapter.ManifestSecrets{}

	binder := &adapter.Binder{}

	err := binder.DeleteBinding("binding1", deploymentTopology, manifest, requestParams, secrets)

	if err != nil {
		t.Fatal(err)
	}

}
