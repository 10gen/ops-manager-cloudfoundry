package adapter_test

import (
	"../../mongodb-service-adapter/adapter"
	"fmt"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	"testing"
)

func TestDashboardUrl(t *testing.T) {
	t.Parallel()

	config, err := adapter.LoadConfig("../../mongodb-service-adapter/testdata/manifest.json")

	if err != nil {
		fmt.Print("Error opening manifest file ")
	}

	
	manifest := bosh.BoshManifest{
		Properties: map[string]interface{}{
			"mongo_ops": map[interface{}]interface{}{
				"url":     config.URL,
				"group_id": config.GroupID,
			},
		},
	}

	exp := &adapter.DashboardURLGenerator{}
	plan := serviceadapter.Plan{}
	
	_, err = exp.DashboardUrl("id1", plan, manifest)

	if err != nil {
		t.Errorf("Error getting DashboardUrl ")
	}
}
