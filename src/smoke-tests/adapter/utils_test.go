package adapter_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"mongodb-service-adapter/adapter"
	"mongodb-service-adapter/adapter/config"
)

var _ = Describe("Utils", func() {

	var (
		nodes         []string
		routers       int
		configServers int
		shards        int
		cluster       *config.Cluster
		err           error
	)

	BeforeEach(func() {

	})

	Describe("GeneratePassword", func() {

		It("returns without error", func() {
			_, err := adapter.GeneratePassword(10)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return length equal to 10", func() {
			res, _ := adapter.GeneratePassword(10)
			Expect(len(res) == 10).To(BeTrue())
		})
	})

	Describe("NodesToCluster", func() {

		BeforeEach(func() {

			nodes = []string{
				"192.168.1.10", "192.168.1.11", "192.168.1.12",
				"192.168.1.13", "192.168.1.14", "192.168.1.15",
				"192.168.1.16", "192.168.1.17", "192.168.1.18",
				"192.168.1.19", "192.168.1.20", "192.168.1.21",
			}

			routers = 2
			configServers = 2
			shards = 4

			cluster, err = adapter.NodesToCluster(nodes, routers, configServers, shards)
		})

		It("returns without error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("Routers should equal to "+string(routers), func() {
			Expect(len(cluster.Routers) == routers).To(BeTrue())
		})

		It("Config servers should equal to "+string(configServers), func() {
			Expect(len(cluster.ConfigServers) == configServers).To(BeTrue())
		})

		It("Shards should equal to "+string(shards), func() {
			count := ((len(nodes) - routers - configServers) / shards)
			Expect(len(cluster.Shards) == count).To(BeTrue())
		})

		It("Cluster should match ", func() {

			want := &config.Cluster{
				Routers:       []string{"192.168.1.10", "192.168.1.11"},
				ConfigServers: []string{"192.168.1.12", "192.168.1.13"},
				Shards: [][]string{
					{"192.168.1.14", "192.168.1.15", "192.168.1.16", "192.168.1.17"},
					{"192.168.1.18", "192.168.1.19", "192.168.1.20", "192.168.1.21"},
				},
			}

			Expect(reflect.DeepEqual(cluster, want)).To(BeTrue())

		})
	})
})
