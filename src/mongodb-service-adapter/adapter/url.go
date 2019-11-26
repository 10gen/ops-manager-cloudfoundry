package adapter

import (
	"fmt"

	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

type DashboardURLGenerator struct{}

var _ serviceadapter.DashboardUrlGenerator = new(DashboardURLGenerator)

func (d *DashboardURLGenerator) DashboardUrl(params serviceadapter.DashboardUrlParams) (serviceadapter.DashboardUrl, error) { // nolint
	properties := params.Manifest.Properties["mongo_ops"].(map[interface{}]interface{})
	url := properties["url"].(string)
	groupID := properties["group_id"].(string)

	return serviceadapter.DashboardUrl{
		DashboardUrl: fmt.Sprintf("%s/v2/%s", url, groupID),
	}, nil
}
