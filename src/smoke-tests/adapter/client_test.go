package adapter_test

// TODO: remove this file?
/*
import (
	"encoding/json"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"mongodb-service-adapter/adapter"
)

var _ = Describe("Client", func() {

	var (
		err           error
		latestVersion string
		config        *adapter.Config
		p             string
		ctx           *adapter.DocContext
		c             *adapter.OMClient
		s             string
	)

	BeforeEach(func() {

		latestVersion = "4.0.7"
		config = mongodb.GetConfig()

		c = &adapter.OMClient{
			Url:      config.URL,
			Username: config.Username,
			ApiKey:   config.APIKey,
		}

	})

	Describe("LoadDoc", func() {
		Context("when plan is standalone", func() {
			BeforeEach(func() {
				p = adapter.PlanStandalone
				ctx = &adapter.DocContext{
					ID:            "d3a98gf1",
					Key:           "key",
					AdminPassword: "pwd",
					Nodes:         []string{"192.168.1.1", "192.168.1.2", "192.168.1.3", "192.168.1.4"},
					Version:       "3.2.11",
				}

				c = &adapter.OMClient{}

				s, err = c.LoadDoc(p, ctx)
			})

			It("runs without error ", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("validates json output", func() {
				err = json.Unmarshal([]byte(s), &map[string]interface{}{})
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when plan is replica_set", func() {
			BeforeEach(func() {
				p = adapter.PlanReplicaSet
				ctx = &adapter.DocContext{
					ID:            "d3a98gf1",
					Key:           "key",
					AdminPassword: "pwd",
					Nodes:         []string{"192.168.1.1", "192.168.1.2", "192.168.1.3", "192.168.1.4"},
					Version:       "3.2.11",
				}

				c = &adapter.OMClient{}

				s, err = c.LoadDoc(p, ctx)
			})

			It("runs without error ", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("validates json output", func() {
				err = json.Unmarshal([]byte(s), &map[string]interface{}{})
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when plan is sharded cluster", func() {
			BeforeEach(func() {
				p = adapter.PlanShardedCluster
				ctx = &adapter.DocContext{
					ID:            "d3a98gf1",
					Key:           "key",
					AdminPassword: "pwd",
					Version:       "3.4.1",
					Cluster: &adapter.Cluster{
						Routers:       []string{"192.168.1.1", "192.168.1.2"},
						ConfigServers: []string{"192.168.1.3", "192.168.1.4"},
						Shards: [][]string{
							{"192.168.0.10", "192.168.0.11", "192.168.0.12"},
							{"192.168.0.13", "192.168.0.14", "192.168.0.15"},
						},
					},
				}

				c = &adapter.OMClient{}

				s, err = c.LoadDoc(p, ctx)
			})

			It("runs without error ", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("validates json output", func() {
				err = json.Unmarshal([]byte(s), &map[string]interface{}{})
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when plan name is wrong", func() {
			BeforeEach(func() {
				p = "wrongPlan"
				ctx = &adapter.DocContext{}
				c = &adapter.OMClient{}
				s, err = c.LoadDoc(p, ctx)
			})

			It("returns error plan not found ", func() {
				Expect(err.Error()).To(ContainSubstring("not found"))
			})
		})
	})

	Describe("GetGroupByName", func() {

		It("returns error when group name is not defined ", func() {
			name, _ := adapter.GeneratePassword(5)
			group, _ := c.GetGroupByName("Why" + name)
			Expect(group.Name == "").To(BeTrue())
		})

		It("returns no error", func() {
			_, err := c.GetGroupByName("Project 0")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("CreateGroup", func() {
		It("returns no error when GroupCreateRequest is nil ", func() {
			id, _ := adapter.GeneratePassword(10)
			request := adapter.GroupCreateRequest{}
			_, err = c.CreateGroup(id, request)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error when group name not exists", func() {
			id, _ := adapter.GeneratePassword(10)
			request := adapter.GroupCreateRequest{
				Name: "PCF" + id[0:(len(id)/2)],
			}
			_, err := c.CreateGroup(id, request)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns no error ", func() {
			id := "some-id"
			request := adapter.GroupCreateRequest{
				Name: "Project 0",
			}
			_, err := c.CreateGroup(id, request)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("CreateGroupAPIKey", func() {

		It("returns error given wrong id but still creates key", func() {
			_, err = c.CreateGroupAPIKey("wrong-group-id")
			Expect(err).To(HaveOccurred())
		})

		It("returns no error", func() {
			_, err := c.CreateGroupAPIKey(config.GroupID)
			Expect(err).ToNot(HaveOccurred())
		})

	})

	Describe("UpdateGroup", func() {

		It("returns no error if GroupUpdateRequest is nil ", func() {
			request := adapter.GroupUpdateRequest{}
			_, err := c.UpdateGroup(config.GroupID, request)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when GroupUpdateRequest is correct", func() {

			var (
				request adapter.GroupUpdateRequest
			)

			BeforeEach(func() {
				request = adapter.GroupUpdateRequest{
					Tags: []string{
						"tag1",
						"tag2",
					},
				}
			})

			It("returns empty group if id is wrong", func() {
				group, _ := c.UpdateGroup("wrong-id", request)
				Expect(reflect.DeepEqual(group, adapter.Group{})).To(BeTrue())
			})

			It("returns no error ", func() {
				_, err = c.UpdateGroup(config.GroupID, request)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("GetGroupAuthAgentPassword", func() {

		It("returns empty if groupID is wrong", func() {
			pwd, _ := c.GetGroupAuthAgentPassword("wrong-id")
			Expect(reflect.DeepEqual(pwd, "")).To(BeTrue())
		})

		It("returns no error", func() {
			_, err = c.GetGroupAuthAgentPassword(config.GroupID)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("GetGroup", func() {

		It("return empty given wrong groupID", func() {
			group, _ := c.GetGroup("wrong-id")
			Expect(reflect.DeepEqual(group, adapter.Group{})).To(BeTrue())
		})

		It("returns no error", func() {
			_, err := c.GetGroup(config.GroupID)
			Expect(err).ToNot(HaveOccurred())
		})

	})

	Describe("DeleteGroup", func() {
		It("returns no error ", func() {
			err = c.DeleteGroup("test-group-id")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("GetGroupHosts", func() {

		It("returns no error", func() {
			_, err = c.GetGroupHosts(config.GroupID)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns 0 given wrong id ", func() {
			groupHosts, _ := c.GetGroupHosts("wrong-group-id")
			Expect(reflect.DeepEqual(groupHosts.TotalCount, 0)).To(BeTrue())
		})
	})

	Describe("GetGroupHostnames", func() {

		It("returns no error", func() {
			_, err = c.GetGroupHostnames(config.GroupID, adapter.PlanStandalone)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns servers if planID is sharded_cluster", func() {
			_, err := c.GetGroupHostnames(config.GroupID, adapter.PlanShardedCluster)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("GetKeyfileWindows", func() {
		It("returns no error", func() {
			_, err := c.GetKeyfileWindows(config.GroupID)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("ConfigureGroup", func() {

		It("returns no error ", func() {

			kfw, err := c.GetKeyfileWindows(config.GroupID)

			ctx := &adapter.DocContext{
				ID:                      "id1",
				Key:                     mongodb.GetConfig().AuthKey,
				AdminPassword:           "admin",
				AutomationAgentPassword: "admin",
				Nodes:                   []string{mongodb.GetConfig().NodeAddresses},
				Version:                 latestVersion,
				RequireSSL:              false,
				KeyfileWindows:          kfw,
			}

			doc, err := c.LoadDoc(adapter.PlanStandalone, ctx)

			Expect(err).ToNot(HaveOccurred())

			err = c.ConfigureGroup(doc, config.GroupID)

			Expect(err).ToNot(HaveOccurred())

		})
	})

	Describe("ConfigureMonitoringAgent", func() {

		It("returns no error ", func() {
			ctx := &adapter.DocContext{
				ID:                      "id1",
				Key:                     mongodb.GetConfig().AuthKey,
				AdminPassword:           "admin",
				AutomationAgentPassword: "admin",
				Nodes:                   []string{mongodb.GetConfig().NodeAddresses},
				Version:                 latestVersion,
				RequireSSL:              false,
			}

			monitoringAgentDoc, err := c.LoadDoc(adapter.MonitoringAgentConfiguration, ctx)

			Expect(err).ToNot(HaveOccurred())

			err = c.ConfigureMonitoringAgent(monitoringAgentDoc, config.GroupID)

			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("ConfigureBackupAgent", func() {

		It("returns no error ", func() {

			ctx := &adapter.DocContext{
				ID:                      "id1",
				Key:                     mongodb.GetConfig().AuthKey,
				AdminPassword:           "admin",
				AutomationAgentPassword: "admin",
				Nodes:                   []string{mongodb.GetConfig().NodeAddresses},
				Version:                 "4.0.7",
				RequireSSL:              false,
			}

			backupAgentDoc, err := c.LoadDoc(adapter.BackupAgentConfiguration, ctx)

			Expect(err).ToNot(HaveOccurred())

			err = c.ConfigureBackupAgent(backupAgentDoc, mongodb.GetConfig().GroupID)

			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("GetAvailableVersions", func() {

		It("returns no error ", func() {

			_, err := c.GetAvailableVersions(config.GroupID)

			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("GetLatestVersion", func() {

		It("returns error given wrong group id ", func() {

			_, err := c.GetLatestVersion("wrong-group-id")

			Expect(err).To(HaveOccurred())

		})

		It("returns no error ", func() {

			_, err := c.GetLatestVersion(config.GroupID)

			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("ValidateVersion", func() {

		It("returns no error ", func() {

			_, err := c.ValidateVersion(config.GroupID, latestVersion)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error given wrong version ", func() {

			_, err := c.ValidateVersion(config.GroupID, "wrong-version")
			Expect(err).To(HaveOccurred())
		})

	})

	Describe("ValidateVersionManifest", func() {

		It("returns no error", func() {

			_, err := c.ValidateVersionManifest("4.0.2")
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error when version is not available", func() {

			_, err := c.ValidateVersionManifest("wrong-version")
			Expect(err).To(HaveOccurred())
		})
	})

})
*/
