package config

import (
	"fmt"

	"github.com/mongodb-labs/pcgc/pkg/opsmanager"
)

type AutomationConfig opsmanager.AutomationConfig

func (r *AutomationConfig) ToStandalone(ctx *DocContext) {
	r.Options.DownloadBase = "/var/lib/mongodb-mms-automation"

	if ctx.RequireSSL {
		r.SSL = getDefaultSSL()
	}

	r.MongoDBVersions, r.BackupVersions, r.MonitoringVersions = getDefaultVersions(ctx.Version, ctx.Nodes[0])

	r.Processes = []*opsmanager.Process{
		{
			Args26: &opsmanager.Args26{
				NET:       getdefaultNet(ctx.RequireSSL),
				SystemLog: getDefaultSystemLog(),
				Storage: &opsmanager.StorageArg{
					DBPath: "/var/vcap/store/mongodb-data",
				},
			},
			Hostname:                    ctx.Nodes[0],
			LogRotate:                   getDefaultLogRotate(),
			Name:                        ctx.Nodes[0],
			ProcessType:                 "mongod",
			Version:                     ctx.Version,
			FeatureCompatibilityVersion: ctx.CompatibilityVersion,
			AuthSchemaVersion:           5,
		},
	}

	r.ReplicaSets = []opsmanager.ReplicaSet{}
	r.Roles = []map[string]interface{}{}
	r.Sharding = []opsmanager.Sharding{}

	r.Auth.AutoUser = "mms-automation"
	r.Auth.AutoPwd = ctx.AutomationAgentPassword
	r.Auth.DeploymentAuthMechanisms = []string{"MONGODB-CR"}
	r.Auth.Key = ctx.Key
	r.Auth.Keyfile = "/var/vcap/store/mongod_node/mongo_om.key"
	// r.Auth.KeyfileWindows = ctx.KeyfileWindows
	r.Auth.Disabled = false
	r.Auth.UsersDeleted = []interface{}{}
	r.Auth.AutoAuthMechanism = "MONGODB-CR"
	r.Auth.UsersWanted = getDefaultUsers(ctx)
}

func (r *AutomationConfig) ToShardedCluster(ctx *DocContext) {
	r.Options.DownloadBase = "/var/lib/mongodb-mms-automation"
	if ctx.RequireSSL {
		r.SSL = getDefaultSSL()
	}

	r.MongoDBVersions, r.BackupVersions, r.MonitoringVersions = getDefaultVersions(ctx.Version, ctx.Cluster.ConfigServers[0])

	protocolVersion := 0
	if ctx.CompatibilityVersion == "4.0" {
		protocolVersion = 1
	}

	r.ReplicaSets = []opsmanager.ReplicaSet{{
		ID:              ctx.ID + "_config",
		ProtocolVersion: protocolVersion,
	}}

	for _, node := range ctx.Cluster.Routers {
		r.Processes = append(r.Processes, &opsmanager.Process{
			Args26: &opsmanager.Args26{
				NET:       getdefaultNet(ctx.RequireSSL),
				SystemLog: getDefaultSystemLog(),
				Storage: &opsmanager.StorageArg{
					DBPath: "/var/vcap/store/mongodb-data",
				},
			},
			Name:                        node,
			Hostname:                    node,
			LogRotate:                   getDefaultLogRotate(),
			Version:                     ctx.Version,
			FeatureCompatibilityVersion: ctx.CompatibilityVersion,
			AuthSchemaVersion:           5,
			ProcessType:                 "mongos",
			Cluster:                     ctx.ID + "_cluster",
		})
	}

	for i, node := range ctx.Cluster.ConfigServers {
		r.Processes = append(r.Processes, &opsmanager.Process{
			Args26: &opsmanager.Args26{
				NET:       getdefaultNet(ctx.RequireSSL),
				SystemLog: getDefaultSystemLog(),
				Storage: &opsmanager.StorageArg{
					DBPath: "/var/vcap/store/mongodb-data",
				},
				Replication: &opsmanager.ReplicationArg{
					ReplSetName: ctx.ID + "_config",
				},
				Sharding: &opsmanager.ShardingArg{
					ClusterRole: "configsvr",
				},
			},
			Name:                        node,
			Hostname:                    node,
			LogRotate:                   getDefaultLogRotate(),
			Version:                     ctx.Version,
			FeatureCompatibilityVersion: ctx.CompatibilityVersion,
			AuthSchemaVersion:           5,
			ProcessType:                 "mongos",
		})

		r.ReplicaSets[0].Members = append(r.ReplicaSets[0].Members, opsmanager.Member{
			ID:          i,
			ArbiterOnly: false,
			Hidden:      false,
			Host:        node,
			Priority:    1,
			SlaveDelay:  0,
			Votes:       1,
		})
	}

	for ii, shard := range ctx.Cluster.Shards {
		rs := opsmanager.ReplicaSet{
			ID:              fmt.Sprintf("%s_shard_%d", ctx.ID, ii),
			ProtocolVersion: protocolVersion,
		}

		for i, node := range shard {
			r.Processes = append(r.Processes, &opsmanager.Process{
				Args26: &opsmanager.Args26{
					NET:       getdefaultNet(ctx.RequireSSL),
					SystemLog: getDefaultSystemLog(),
					Storage: &opsmanager.StorageArg{
						DBPath: "/var/vcap/store/mongodb-data",
					},
					Replication: &opsmanager.ReplicationArg{
						ReplSetName: fmt.Sprintf("%s_shard_%d", ctx.ID, ii),
					},
				},
				Name:                        node,
				Hostname:                    node,
				LogRotate:                   getDefaultLogRotate(),
				Version:                     ctx.Version,
				FeatureCompatibilityVersion: ctx.CompatibilityVersion,
				AuthSchemaVersion:           5,
				ProcessType:                 "mongod",
			})

			rs.Members = append(rs.Members, opsmanager.Member{
				ID:          i,
				ArbiterOnly: false,
				Hidden:      false,
				Host:        node,
				Priority:    1,
				SlaveDelay:  0,
				Votes:       1,
			})
		}

		r.ReplicaSets = append(r.ReplicaSets, rs)
	}

	r.Sharding = []opsmanager.Sharding{{
		Name:                ctx.ID + "_cluster",
		ConfigServer:        []interface{}{},
		ConfigServerReplica: ctx.ID + "_config",
		Collections:         []interface{}{},
	}}

	for i := range ctx.Cluster.Shards {
		id := fmt.Sprintf("%s_shard_%d", ctx.ID, i)
		r.Sharding[0].Shards = append(r.Sharding[0].Shards, opsmanager.Shard{
			ID:   id,
			Rs:   id,
			Tags: []interface{}{},
		})
	}

	r.Auth.AutoUser = "mms-automation"
	r.Auth.AutoPwd = ctx.AutomationAgentPassword
	r.Auth.DeploymentAuthMechanisms = []string{"MONGODB-CR"}
	r.Auth.Key = ctx.Key
	r.Auth.Keyfile = "/var/vcap/store/mongod_node/mongo_om.key"
	// r.Auth.KeyfileWindows = ctx.KeyfileWindows
	r.Auth.Disabled = false
	r.Auth.UsersDeleted = []interface{}{}
	r.Auth.AutoAuthMechanism = "MONGODB-CR"
	r.Auth.UsersWanted = getDefaultUsers(ctx)
}

func (r *AutomationConfig) ToReplicaSet(ctx *DocContext) {
	r.Options.DownloadBase = "/var/lib/mongodb-mms-automation"

	if ctx.RequireSSL {
		r.SSL = getDefaultSSL()
	}

	r.MongoDBVersions, r.BackupVersions, r.MonitoringVersions = getDefaultVersions(ctx.Version, ctx.Nodes[0])

	protocolVersion := 0
	if ctx.CompatibilityVersion == "4.0" {
		protocolVersion = 1
	}
	r.ReplicaSets = []opsmanager.ReplicaSet{{
		ID:              "pcf_repl",
		ProtocolVersion: protocolVersion,
	}}

	for i, node := range ctx.Nodes {
		r.Processes = append(r.Processes, &opsmanager.Process{
			Args26: &opsmanager.Args26{
				NET:       getdefaultNet(ctx.RequireSSL),
				SystemLog: getDefaultSystemLog(),
				Storage: &opsmanager.StorageArg{
					DBPath: "/var/vcap/store/mongodb-data",
				},
				Replication: &opsmanager.ReplicationArg{
					ReplSetName: "pcf_repl",
				},
			},
			Name:                        node,
			Hostname:                    node,
			LogRotate:                   getDefaultLogRotate(),
			ProcessType:                 "mongod",
			Version:                     ctx.Version,
			FeatureCompatibilityVersion: ctx.CompatibilityVersion,
			AuthSchemaVersion:           5,
		})

		r.ReplicaSets[0].Members = append(r.ReplicaSets[0].Members, opsmanager.Member{
			ID:   i,
			Host: node,
		})
	}

	r.Roles = []map[string]interface{}{}
	r.Sharding = []opsmanager.Sharding{}

	r.Auth.AutoUser = "mms-automation"
	r.Auth.AutoPwd = ctx.AutomationAgentPassword
	r.Auth.DeploymentAuthMechanisms = []string{"MONGODB-CR"}
	r.Auth.Key = ctx.Key
	r.Auth.Keyfile = "/var/vcap/store/mongod_node/mongo_om.key"
	// r.Auth.KeyfileWindows = ctx.KeyfileWindows
	r.Auth.Disabled = false
	r.Auth.UsersDeleted = []interface{}{}
	r.Auth.AutoAuthMechanism = "MONGODB-CR"
	r.Auth.UsersWanted = getDefaultUsers(ctx)
}
