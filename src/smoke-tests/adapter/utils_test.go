package adapter_test

import (
	"github.com/10gen/ops-manager-cloudfoundry/src/mongodb-service-adapter/adapter"
	"reflect"
	"testing"
)

func TestGenerateString(t *testing.T) {
	t.Parallel()

	res, err := adapter.GenerateString(10)

	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 10 {
		t.Fatal("String length is not correct ")
	}

}

func TestNodesToCluster(t *testing.T) {
	t.Parallel()

	nodes := []string{
		"192.168.1.10", "192.168.1.11", "192.168.1.12",
		"192.168.1.13", "192.168.1.14", "192.168.1.15",
		"192.168.1.16", "192.168.1.17", "192.168.1.18",
		"192.168.1.19", "192.168.1.20", "192.168.1.21",
	}

	routers := 2
	configServers := 2
	shards := 4

	res, err := adapter.NodesToCluster(nodes, routers, configServers, shards)

	if err != nil {
		t.Fatal(err)
	}

	if routers != len(res.Routers) {
		t.Fatal("Routers not equal to ", routers)
	}

	if configServers != len(res.ConfigServers) {
		t.Fatal("Config Servers not equal to ", configServers)
	}

	if count := ((len(nodes) - routers - configServers) / shards); count != len(res.Shards) {
		t.Fatal("Shards not equal to ", count)
	}

	want := &adapter.Cluster{
		Routers:       []string{"192.168.1.10", "192.168.1.11"},
		ConfigServers: []string{"192.168.1.12", "192.168.1.13"},
		Shards: [][]string{
			{"192.168.1.14", "192.168.1.15", "192.168.1.16", "192.168.1.17"},
			{"192.168.1.18", "192.168.1.19", "192.168.1.20", "192.168.1.21"},
		},
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Cluster = %#v, want %#v", res, want)
	}
}
