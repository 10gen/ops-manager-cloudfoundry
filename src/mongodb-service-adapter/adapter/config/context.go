package config

type DocContext struct {
	ID                      string
	Key                     string
	AdminPassword           string
	Version                 string
	AutomationAgentPassword string
	CompatibilityVersion    string
	Nodes                   []string
	Cluster                 *Cluster
	Password                string
	RequireSSL              bool
}

type Cluster struct {
	Routers       []string
	ConfigServers []string
	Shards        [][]string
}
