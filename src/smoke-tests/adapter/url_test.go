package adapter_test

import (
	"../../mongodb-service-adapter/adapter"
	"fmt"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	"os"
	"testing"
)

func TestDashboardUrl(t *testing.T) {
	t.Parallel()

	manifest := bosh.BoshManifest{
		Properties: map[string]interface{}{
			"mongo_ops": map[interface{}]interface{}{
				"url":      os.Getenv("Url"),
				"group_id": os.Getenv("GroupId"),
			},
		},
	}

	exp := &adapter.DashboardURLGenerator{}
	plan := serviceadapter.Plan{}
	res, err := exp.DashboardUrl("id1", plan, manifest)
	if err != nil {
		t.Errorf("Error getting DashboardUrl ")
	}
	fmt.Println("DashboardUrl : ", res)
}
