package adapter

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"

	"mongodb-service-adapter/adapter/config"
)

// GenerateString generates a random string or panics
// if something goes wrong.
func GenerateString(l int) (string, error) {
	b := make([]byte, l)
	for i := l; i != 0; {
		n, err := rand.Read(b)
		if err != nil {
			return "", err
		}
		if n == 0 {
			return "", errors.New("couldn't read from crypto/rand")
		}

		i -= n
	}

	return fmt.Sprintf("%x", b)[:l], nil
}

// TODO: validate input
// NodesToCluster transforms a nodes list into cluster configuration object.
func NodesToCluster(nodes []string, routers, configServers, replicas int) (*config.Cluster, error) {
	// nodes have to be ordered because
	// bosh provides them in random order
	sort.Slice(nodes, func(i, j int) bool {
		return addrn(nodes[i]) < addrn(nodes[j])
	})

	c := &config.Cluster{
		Routers:       nodes[:routers],
		ConfigServers: nodes[routers : routers+configServers],
	}

	nodes = nodes[routers+configServers:]
	c.Shards = make([][]string, 0, len(nodes)/replicas)
	for i := 0; i < len(nodes)/replicas; i++ {
		c.Shards = append(c.Shards, make([]string, 0, replicas))
		for j := 0; j < replicas; j++ {
			c.Shards[i] = append(c.Shards[i], nodes[i*replicas+j])
		}
	}
	return c, nil
}

func addrn(addr string) int {
	if !strings.Contains(addr, ":") {
		addr = addr + ":0"
	}

	n := 0
	a, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}

	for _, b := range a.IP {
		n += int(b)
	}
	return n
}
