package adapter_test

import (
	"fmt"
	"github.com/10gen/ops-manager-cloudfoundry/src/mongodb-service-adapter/adapter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

var _ = Describe("Manifest", func() {

	var (
		serviceDeployment serviceadapter.ServiceDeployment
		plan              serviceadapter.Plan
		requestParams     serviceadapter.RequestParameters
		previousManifest  *bosh.BoshManifest
		previousPlan      *serviceadapter.Plan
		config            *adapter.Config
		mGenerator        *adapter.ManifestGenerator
		err               error
	)

	BeforeEach(func() {
		config, err = adapter.LoadConfig("../../mongodb-service-adapter/testdata/manifest.json")

		if err != nil {
			fmt.Println("Error opening manifest file ", err)
			return
		}

		serviceDeployment = serviceadapter.ServiceDeployment{
			Releases: serviceadapter.ServiceReleases{
				serviceadapter.ServiceRelease{
					Name:    "standalone",
					Version: "4.0.7",
					Jobs: []string{adapter.MongodJobName, adapter.BPMJobName,
						adapter.SyslogJobName, adapter.ConfigAgentJobName,
						adapter.CleanupErrandJobName, adapter.PostSetupErrandJobName,
						adapter.BoshDNSEnableJobName},
				},
			},
			DeploymentName: "deploy_name",
			Stemcell: serviceadapter.Stemcell{
				OS:      "Ubuntu",
				Version: "16.X",
			},
		}

		plan = serviceadapter.Plan{
			Properties: serviceadapter.Properties{
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
					"port":           "28000",
					"transport":      "tls",
					"tls_enabled":    false,
					"permitted_peer": 1,
					"ca_cert":        "ca_cert",
				},
				"id": "standalone",
			},
			InstanceGroups: []serviceadapter.InstanceGroup{
				{
					Name:     adapter.MongodJobName,
					Networks: []string{"network1"},
				},
			},
		}

		requestParams = serviceadapter.RequestParameters{
			"parameters": map[string]interface{}{
				"project_name": "Project 0",
				"orgId":        config.OrgID,
			},
		}

		previousManifest = &bosh.BoshManifest{
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

		previousPlan = &serviceadapter.Plan{}
		mGenerator = &adapter.ManifestGenerator{}
		_, err = mGenerator.GenerateManifest(serviceDeployment, plan, requestParams, previousManifest, previousPlan, nil)

	})

	Describe("GenerateManifest", func() {

		It("runs without error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if plan instance group name is different ", func() {
			plan.InstanceGroups = []serviceadapter.InstanceGroup{
				{
					Name:     "wrong-name",
					Networks: []string{"network1"},
				},
			}

			_, err = mGenerator.GenerateManifest(serviceDeployment, plan, requestParams, previousManifest, previousPlan, nil)
			Expect(err).To(HaveOccurred())
		})

		It("returns error when no jobs exists ", func() {
			serviceDeployment.Releases[0].Jobs = []string{}
			_, err = mGenerator.GenerateManifest(serviceDeployment, plan, requestParams, previousManifest, previousPlan, nil)

			Expect(err).To(HaveOccurred())
		})

		It("returns error when bpm job not exists ", func() {
			serviceDeployment.Releases[0].Jobs = []string{adapter.MongodJobName}
			_, err = mGenerator.GenerateManifest(serviceDeployment, plan, requestParams, previousManifest, previousPlan, nil)

			Expect(err.Error()).To(ContainSubstring("no release provided for job 'bpm'"))
		})

		It("returns error when syslog_forwarder job not exists ", func() {
			serviceDeployment.Releases[0].Jobs = []string{adapter.MongodJobName,
				adapter.BPMJobName}
			_, err = mGenerator.GenerateManifest(serviceDeployment, plan, requestParams, previousManifest, previousPlan, nil)

			Expect(err.Error()).To(ContainSubstring("no release provided for job 'syslog_forwarder'"))
		})

		It("returns error when mongodb_config_agent job not exists ", func() {
			serviceDeployment.Releases[0].Jobs = []string{adapter.MongodJobName,
				adapter.BPMJobName, adapter.SyslogJobName}
			_, err = mGenerator.GenerateManifest(serviceDeployment, plan, requestParams, previousManifest, previousPlan, nil)

			Expect(err.Error()).To(ContainSubstring("no release provided for job 'mongodb_config_agent'"))
		})

		It("returns error when cleanup_service job not exists ", func() {
			serviceDeployment.Releases[0].Jobs = []string{adapter.MongodJobName,
				adapter.BPMJobName, adapter.SyslogJobName, adapter.ConfigAgentJobName, adapter.PostSetupErrandJobName}
			_, err = mGenerator.GenerateManifest(serviceDeployment, plan, requestParams, previousManifest, previousPlan, nil)

			Expect(err.Error()).To(ContainSubstring("no release provided for job 'cleanup_service'"))
		})

		It("returns error when bosh-dns-enable job not exists ", func() {
			serviceDeployment.Releases[0].Jobs = []string{adapter.MongodJobName,
				adapter.BPMJobName, adapter.CleanupErrandJobName, adapter.SyslogJobName, adapter.ConfigAgentJobName, adapter.PostSetupErrandJobName}
			_, err = mGenerator.GenerateManifest(serviceDeployment, plan, requestParams, previousManifest, previousPlan, nil)

			Expect(err.Error()).To(ContainSubstring("no release provided for job 'bosh-dns-enable'"))
		})

		It("returns error when no networks definition exists ", func() {
			serviceDeployment.Releases[0].Jobs = []string{adapter.MongodJobName,
				adapter.BPMJobName, adapter.BoshDNSEnableJobName, adapter.CleanupErrandJobName, adapter.SyslogJobName, adapter.ConfigAgentJobName,
				adapter.PostSetupErrandJobName}
			plan.InstanceGroups[0].Networks = []string{}

			_, err = mGenerator.GenerateManifest(serviceDeployment, plan, requestParams, previousManifest, previousPlan, nil)
			Expect(err.Error()).To(ContainSubstring("no networks definition found"))
		})

		It("returns error when previousManifest instance group id not exists ", func() {

			previousManifest.InstanceGroups[1].Properties["mongo_ops"] = map[interface{}]interface{}{
				"admin_password": "admin",
				"id":             "standalone",
				"group_id":       "",
				"agent_api_key":  config.APIKey,
				"auth_key":       config.AuthKey,
			}

			_, err = mGenerator.GenerateManifest(serviceDeployment, plan, requestParams, previousManifest, previousPlan, nil)
			Expect(err).To(HaveOccurred())
		})

		It("returns error when wrong plan id is provided", func() {

			plan.Properties["id"] = "wrong_plan"

			_, err = mGenerator.GenerateManifest(serviceDeployment, plan, requestParams, previousManifest, previousPlan, nil)
			Expect(err.Error()).To(ContainSubstring("unknown plan"))
		})

		It("returns unable to find the latest MongoDB version from the MongoDB", func() {

			previousManifest.InstanceGroups[1].Properties = map[string]interface{}{
				"mongo_ops": map[interface{}]interface{}{
					"admin_password": "admin",
					"id":             "standalone",
					"group_id":       "0123456789",
					"agent_api_key":  config.APIKey,
					"auth_key":       config.AuthKey,
				},
			}

			_, err = mGenerator.GenerateManifest(serviceDeployment, plan, requestParams, previousManifest, previousPlan, nil)
			Expect(err.Error()).To(ContainSubstring("unable to find the latest MongoDB version from the MongoDB Ops Manager API"))
		})
	})

})
