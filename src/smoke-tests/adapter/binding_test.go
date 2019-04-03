package adapter_test

import (
	adapter "../../mongodb-service-adapter/adapter"
	"fmt"
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

		createBindingAction, createBindingError = binder.CreateBinding(bindingId, deploymentTopology, manifest, requestParams, secrets, dnsAddresses)
		deleteBindingError = binder.DeleteBinding(bindingId, deploymentTopology, manifest, requestParams, secrets)
		_, getWithCredentialError = adapter.GetWithCredentials(addrs, adminPassword, false)
	})

	Describe("CreateBinding", func() {
		Context("when nothing is missing ", func() {
			It("calls CreateBinding without error ", func() {
				Expect(createBindingError).ToNot(HaveOccurred())
			})

			It("has credentials", func() {
				fmt.Println(createBindingAction.Credentials)
				Expect(createBindingAction.Credentials != nil).To(BeTrue())
			})
		})

		Context("when bindingId is missing", func() {
			BeforeEach(func() {
				bindingId = ""
			})

			It("returns no error", func() {
				_, err = binder.CreateBinding(bindingId, deploymentTopology, manifest, requestParams, secrets, dnsAddresses)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when deploymentTopology is missing", func() {
			BeforeEach(func() {
				deploymentTopology = nil
			})

			It("returns error", func() {
				_, err = binder.CreateBinding(bindingId, deploymentTopology, manifest, requestParams, secrets, dnsAddresses)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when optional parameters : requestParams, secrets and dnsAddresses are missing", func() {
			BeforeEach(func() {
				requestParams = nil
				secrets = nil
				dnsAddresses = nil
			})

			It("returns no error", func() {
				_, err = binder.CreateBinding(bindingId, deploymentTopology, manifest, requestParams, secrets, dnsAddresses)
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
			})

			It("contains replicaSet=pcf_repl", func() {
				createBindingAction, err = binder.CreateBinding(bindingId, deploymentTopology, manifest, requestParams, secrets, dnsAddresses)
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
							"routers":        0,
							"config_servers": 0,
							"replicas":       1,
						},
					},
				}
			})

			It("returns error ", func() {
				_, err = binder.CreateBinding(bindingId, deploymentTopology, manifest, requestParams, secrets, dnsAddresses)
				Expect(err).To(HaveOccurred())
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
				bindingId = "qwerty"
			})

			It("returns error", func() {
				err = binder.DeleteBinding(bindingId, deploymentTopology, manifest, requestParams, secrets)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when deploymentTopology is missing", func() {
			BeforeEach(func() {
				deploymentTopology = nil
			})

			It("returns error", func() {
				err = binder.DeleteBinding(bindingId, deploymentTopology, manifest, requestParams, secrets)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetWitCredentials", func() {
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
