package adapter_test

import (
	"fmt"
	"mongodb-service-adapter/adapter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

var _ = Describe("Binding", func() {

	var (
		bindingId              string
		deploymentTopology     bosh.BoshVMs
		manifest               bosh.BoshManifest
		requestParams          serviceadapter.RequestParameters
		secrets                serviceadapter.ManifestSecrets
		dnsAddresses           serviceadapter.DNSAddresses
		binder                 *adapter.Binder
		createBindingAction    serviceadapter.Binding
		deleteBindingError     error
		createBindingError     error
		getWithCredentialError error
		config                 *adapter.Config
		err                    error
		adminPassword          string
		addrs                  []string
		createBindingParams    serviceadapter.CreateBindingParams
		deleteBindingParams    serviceadapter.DeleteBindingParams
	)

	BeforeEach(func() {
		config, err = adapter.LoadConfig("../../mongodb-service-adapter/testdata/manifest.json")

		if err != nil {
			fmt.Print("Error opening manifest file ")
		}

		bindingId = "id1"
		adminPassword = "admin"
		addrs = []string{config.NodeAddresses + ":28000"}

		deploymentTopology = bosh.BoshVMs{
			"mongod_node": []string{
				config.NodeAddresses,
			},
		}

		manifest = bosh.BoshManifest{
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

		requestParams = serviceadapter.RequestParameters{}
		secrets = serviceadapter.ManifestSecrets{}
		dnsAddresses = serviceadapter.DNSAddresses{}
		binder = &adapter.Binder{}
		createBindingParams = serviceadapter.CreateBindingParams{
			BindingID:          bindingId,
			DeploymentTopology: deploymentTopology,
			Manifest:           manifest,
			RequestParams:      requestParams,
			Secrets:            secrets,
			DNSAddresses:       dnsAddresses,
		}
		deleteBindingParams = serviceadapter.DeleteBindingParams{
			BindingID:          bindingId,
			DeploymentTopology: deploymentTopology,
			Manifest:           manifest,
			RequestParams:      requestParams,
			Secrets:            secrets,
			DNSAddresses:       dnsAddresses,
		}
		createBindingAction, createBindingError = binder.CreateBinding(createBindingParams)
		deleteBindingError = binder.DeleteBinding(deleteBindingParams)
		_, getWithCredentialError = adapter.GetWithCredentials(addrs, adminPassword, false)
	})

	Describe("CreateBinding", func() {
		Context("when nothing is missing ", func() {
			It("calls CreateBinding without error ", func() {
				Expect(createBindingError).ToNot(HaveOccurred())
			})

			It("has credentials", func() {
				Expect(createBindingAction.Credentials != nil).To(BeTrue())
			})
		})

		Context("when bindingId is missing", func() {
			BeforeEach(func() {
				bindingId = ""
			})

			It("returns no error", func() {
				_, err = binder.CreateBinding(createBindingParams)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when deploymentTopology is missing", func() {
			BeforeEach(func() {
				createBindingParams.DeploymentTopology = nil
			})

			It("returns error", func() {
				_, err = binder.CreateBinding(createBindingParams)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when optional parameters : requestParams, secrets and dnsAddresses are missing", func() {
			BeforeEach(func() {
				createBindingParams.RequestParams = nil
				createBindingParams.Secrets = nil
				createBindingParams.DNSAddresses = nil
			})

			It("returns no error", func() {
				_, err = binder.CreateBinding(createBindingParams)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when plan is replica_set", func() {
			BeforeEach(func() {
				manifest = bosh.BoshManifest{
					Properties: map[string]interface{}{
						"mongo_ops": map[interface{}]interface{}{
							"url":            config.URL,
							"group_id":       config.GroupID,
							"admin_password": "admin",
							"username":       config.Username,
							"admin_api_key":  config.APIKey,
							"require_ssl":    false,
							"plan_id":        adapter.PlanReplicaSet,
							"routers":        0,
							"config_servers": 0,
							"replicas":       0,
						},
					},
				}
				createBindingParams.Manifest = manifest
			})

			It("contains replicaSet=pcf_repl", func() {
				createBindingAction, err = binder.CreateBinding(createBindingParams)
				Expect(createBindingAction.Credentials["uri"]).To(ContainSubstring("replicaSet=pcf_repl"))
			})
		})

		Context("when plan is sharded_cluster", func() {
			BeforeEach(func() {
				manifest = bosh.BoshManifest{
					Properties: map[string]interface{}{
						"mongo_ops": map[interface{}]interface{}{
							"url":            config.URL,
							"group_id":       config.GroupID,
							"admin_password": "admin",
							"username":       config.Username,
							"admin_api_key":  config.APIKey,
							"require_ssl":    false,
							"plan_id":        adapter.PlanShardedCluster,
							"routers":        1,
							"config_servers": 0,
							"replicas":       1,
						},
					},
				}
				createBindingParams.Manifest = manifest
			})

			It("returns no error ", func() {
				_, err = binder.CreateBinding(createBindingParams)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("DeleteBinding", func() {
		Context("when nothing is missing ", func() {
			It("calls deleteBinding without error ", func() {
				Expect(deleteBindingError).ToNot(HaveOccurred())
			})
		})

		Context("when bindingId is different ", func() {
			BeforeEach(func() {
				deleteBindingParams.BindingID = "qwerty"
			})

			It("returns error", func() {
				err = binder.DeleteBinding(deleteBindingParams)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when deploymentTopology is missing", func() {
			BeforeEach(func() {
				deleteBindingParams.DeploymentTopology = nil
			})

			It("returns error", func() {
				err = binder.DeleteBinding(deleteBindingParams)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetWithCredentials", func() {
		Context("when nothing is missing ", func() {
			It("calls GetWithCredentials without error ", func() {
				Expect(getWithCredentialError).ToNot(HaveOccurred())
			})
		})

		Context("when addrs is missing ", func() {
			BeforeEach(func() {
				addrs = nil
			})

			It("returns error", func() {
				_, err = adapter.GetWithCredentials(addrs, adminPassword, false)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when admin password is wrong ", func() {
			BeforeEach(func() {
				adminPassword = "wrongPassword"
			})

			It("returns error", func() {
				_, err = adapter.GetWithCredentials(addrs, adminPassword, false)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
