package adapter_test

import (
	"fmt"
	"github.com/10gen/ops-manager-cloudfoundry/src/mongodb-service-adapter/adapter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

var _ = Describe("Url", func() {

	var (
		err      error
		exp      adapter.DashboardURLGenerator
		plan     serviceadapter.Plan
		manifest bosh.BoshManifest
	)

	BeforeEach(func() {

		config, err := adapter.LoadConfig("../../mongodb-service-adapter/testdata/manifest.json")

		if err != nil {
			fmt.Print("Error opening manifest file ")
		}

		manifest = bosh.BoshManifest{
			Properties: map[string]interface{}{
				"mongo_ops": map[interface{}]interface{}{
					"url":      config.URL,
					"group_id": config.GroupID,
				},
			},
		}

		exp = adapter.DashboardURLGenerator{}
		plan = serviceadapter.Plan{}

	})

	Describe("DashboardUrl", func() {

		When("When nothing is missing", func() {

			It("calls DashboardUrl without error", func() {
				_, err = exp.DashboardUrl("id1", plan, manifest)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
