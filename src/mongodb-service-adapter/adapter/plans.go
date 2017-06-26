package adapter

import (
	"reflect"
	"text/template"
)

const (
	PlanStandalone     = "standalone"
	PlanReplicaSet     = "replica_set"
	PlanShardedCluster = "sharded_cluster"
)

var plans = map[string]*template.Template{}

func init() {
	funcs := template.FuncMap{
		"last": func(a interface{}, x int) bool {
			return reflect.ValueOf(a).Len()-1 == x
		},
	}

	var err error
	for k, s := range plansRaw {
		plans[k], err = template.New(string(k)).Funcs(funcs).Parse(s)
		if err != nil {
			panic(err)
		}
	}
}

var plansRaw = map[string]string{
	PlanStandalone: `{
    "options": {
        "downloadBase": "/var/lib/mongodb-mms-automation"
    },
    "mongoDbVersions": [
        {"name": "{{.Version}}"}
    ],
    "backupVersions": [{
        "hostname": "{{index .Nodes 0}}",
        "logPath": "/var/vcap/sys/log/mongod_node/backup-agent.log",
        "logRotate": {
            "sizeThresholdMB": 1000,
            "timeThresholdHrs": 24
        }
    }],
    "monitoringVersions": [
        {{range $i, $node := .Nodes}}{{if $i}},{{end}}{
            "hostname": "{{$node}}",
            "logPath": "/var/vcap/sys/log/mongod_node/monitoring-agent.log",
            "logRotate": {
                "sizeThresholdMB": 1000,
                "timeThresholdHrs": 24
            }
        }{{end}}
    ],
    "processes": [{
        "args2_6": {
            "net": {
                "port": 28000
            },
            "storage": {
                "dbPath": "/var/vcap/store/mongodb-data"
            },
            "systemLog": {
                "destination": "file",
                "path": "/var/vcap/sys/log/mongod_node/mongodb.log"
            }
        },
        "hostname": "{{index .Nodes 0}}",
        "logRotate": {
            "sizeThresholdMB": 1000,
            "timeThresholdHrs": 24
        },
        "name": "{{index .Nodes 0}}",
        "processType": "mongod",
        "version": "{{.Version}}",
        "authSchemaVersion": 5
    }],
    "replicaSets": [],
    "roles": [],
    "sharding": [],

    "auth": {
        "autoUser": "mms-automation",
        "autoPwd": "{{.Password}}",
        "deploymentAuthMechanisms": [
            "SCRAM-SHA-1"
        ],
        "key": "{{.Key}}",
        "keyfile": "/var/vcap/jobs/mongod_node/config/mongo_om.key",
        "disabled": false,
        "usersDeleted": [],
        "usersWanted": [
            {
                "db": "admin",
                "roles": [
                    {
                        "db": "admin",
                        "role": "clusterMonitor"
                    }
                ],
                "user": "mms-monitoring-agent",
                "initPwd": "{{.Password}}"
            },
            {
                "db": "admin",
                "roles": [
                    {
                        "db": "admin",
                        "role": "clusterAdmin"
                    },
                    {
                        "db": "admin",
                        "role": "readAnyDatabase"
                    },
                    {
                        "db": "admin",
                        "role": "userAdminAnyDatabase"
                    },
                    {
                        "db": "local",
                        "role": "readWrite"
                    },
                    {
                        "db": "admin",
                        "role": "readWrite"
                    }
                ],
                "user": "mms-backup-agent",
                "initPwd": "{{.Password}}"
            },
            {
               "db": "admin" ,
               "user": "admin" ,
               "roles": [
                 {
                     "db": "admin",
                     "role": "clusterAdmin"
                 },
                 {
                     "db": "admin",
                     "role": "readAnyDatabase"
                 },
                 {
                     "db": "admin",
                     "role": "userAdminAnyDatabase"
                 },
                 {
                     "db": "local",
                     "role": "readWrite"
                 },
                 {
                     "db": "admin",
                     "role": "readWrite"
                 }
               ],
               "initPwd": "{{.AdminPassword}}"
            }
        ],
        "autoAuthMechanism": "SCRAM-SHA-1"
    }
}`,

	PlanShardedCluster: `{
    "options": {
        "downloadBase": "/var/lib/mongodb-mms-automation"
    },
    "mongoDbVersions": [
        {"name": "{{.Version}}"}
    ],
    "backupVersions": [
    ],
    "monitoringVersions": [
        {{range $i, $node := .Nodes}}{{if $i}},{{end}}{
            "hostname": "{{$node}}",
            "logPath": "/var/vcap/sys/log/mongod_node/monitoring-agent.log",
            "logRotate": {
                "sizeThresholdMB": 1000,
                "timeThresholdHrs": 24
            }
        }{{end}}
    ],
    "processes": [
      {{range $i, $node := .Cluster.Routers}}{
          "args2_6": {
              "net": {
                  "port": 28000
              },
              "systemLog": {
                  "destination": "file",
                  "path": "/var/vcap/sys/log/mongod_node/mongodb.log"
              }
          },
          "name": "{{$node}}",
          "hostname": "{{$node}}",
          "logRotate": {
              "sizeThresholdMB": 1000,
              "timeThresholdHrs": 24
          },
          "version": "{{$.Version}}",
          "authSchemaVersion": 5,
          "processType": "mongos",
          "cluster": "{{$.ID}}_cluster"
      },{{end}}

      {{range $i, $node := .Cluster.ConfigServers}}{
          "args2_6": {
              "net": {
                  "port": 28000
              },
              "replication": {
                  "replSetName": "{{$.ID}}_config"
              },
              "sharding": {
                "clusterRole": "configsvr"
              },
              "storage": {
                  "dbPath": "/var/vcap/store/mongodb-data"
              },
              "systemLog": {
                  "destination": "file",
                  "path": "/var/vcap/sys/log/mongod_node/mongodb.log"
              }
          },
          "name": "{{$node}}",
          "hostname": "{{$node}}",
          "logRotate": {
              "sizeThresholdMB": 1000,
              "timeThresholdHrs": 24
          },
          "version": "{{$.Version}}",
          "authSchemaVersion": 5,
          "processType": "mongod"
      }{{if last $.Cluster.ConfigServers $i}}{{else}},{{end}}{{end}}

      {{range $ii, $shard := .Cluster.Shards}}
          {{range $i, $node := $shard}},{
              "args2_6": {
                  "net": {
                      "port": 28000
                  },
                  "replication": {
                      "replSetName": "{{$.ID}}_shard_{{$ii}}"
                  },
                  "storage": {
                       "dbPath": "/var/vcap/store/mongodb-data"
                  },
                  "systemLog": {
                      "destination": "file",
                      "path": "/var/vcap/sys/log/mongod_node/mongodb.log"
                  }
              },
              "name": "{{$node}}",
              "hostname": "{{$node}}",
              "logRotate": {
                  "sizeThresholdMB": 1000,
                  "timeThresholdHrs": 24
              },
              "version": "{{$.Version}}",
              "authSchemaVersion": 5,
              "processType": "mongod"
          }{{end}}
      {{end}}
    ],

    "replicaSets": [{
        "_id": "{{$.ID}}_config",
        "members": [
            {{range $i, $node := .Cluster.ConfigServers}}{{if $i}},{{end}}{
                "_id": {{$i}},
                "arbiterOnly": false,
                "hidden": false,
                "host": "{{$node}}",
                "priority": 1,
                "slaveDelay": 0,
                "votes": 1
            }{{end}}
        ]
    }
    {{range $i, $shard := .Cluster.Shards}},{
        "_id": "{{$.ID}}_shard_{{$i}}",
        "members": [{{range $i, $node := $shard}}
            {{if $i}},{{end}}{
                "_id": {{$i}},
                "arbiterOnly": false,
                "hidden": false,
                "host": "{{$node}}",
                "priority": 1,
                "slaveDelay": 0,
                "votes": 1
            }
            {{end}}
        ]
    }{{end}}],

    "sharding": [{
        "shards": [
             {{range $i, $shard := .Cluster.Shards}}{{if $i}},{{end}}{
                 "tags": [],
                 "_id": "{{$.ID}}_shard_{{$i}}",
                 "rs": "{{$.ID}}_shard_{{$i}}"
             }{{end}}
        ],
        "name": "{{.ID}}_cluster",
        "configServer": [],
        "configServerReplica": "{{.ID}}_config",
        "collections": []
    }],

    "auth": {
        "autoUser": "mms-automation",
        "autoPwd": "{{.Password}}",
        "deploymentAuthMechanisms": [
            "SCRAM-SHA-1"
        ],
        "key": "{{.Key}}",
        "keyfile": "/var/vcap/jobs/mongod_node/config/mongo_om.key",
        "disabled": false,
        "usersDeleted": [],
        "usersWanted": [
            {
                "db": "admin",
                "roles": [
                    {
                        "db": "admin",
                        "role": "clusterMonitor"
                    }
                ],
                "user": "mms-monitoring-agent",
                "initPwd": "{{.Password}}"
            },
            {
                "db": "admin",
                "roles": [
                    {
                        "db": "admin",
                        "role": "clusterAdmin"
                    },
                    {
                        "db": "admin",
                        "role": "readAnyDatabase"
                    },
                    {
                        "db": "admin",
                        "role": "userAdminAnyDatabase"
                    },
                    {
                        "db": "local",
                        "role": "readWrite"
                    },
                    {
                        "db": "admin",
                        "role": "readWrite"
                    }
                ],
                "user": "mms-backup-agent",
                "initPwd": "{{.Password}}"
            },
            {
               "db": "admin" ,
               "user": "admin" ,
               "roles": [
                 {
                     "db": "admin",
                     "role": "clusterAdmin"
                 },
                 {
                     "db": "admin",
                     "role": "readAnyDatabase"
                 },
                 {
                     "db": "admin",
                     "role": "userAdminAnyDatabase"
                 },
                 {
                     "db": "local",
                     "role": "readWrite"
                 },
                 {
                     "db": "admin",
                     "role": "readWrite"
                 }
               ],
               "initPwd": "{{.AdminPassword}}"
            }
        ],
        "autoAuthMechanism": "SCRAM-SHA-1"
    }
}`,

	PlanReplicaSet: `{
    "options": {
        "downloadBase": "/var/lib/mongodb-mms-automation"
    },
    "mongoDbVersions": [
        {"name": "{{.Version}}"}
    ],
    "backupVersions": [{
        "hostname": "{{index .Nodes 0}}",
        "logPath": "/var/vcap/sys/log/mongod_node/backup-agent.log",
        "logRotate": {
            "sizeThresholdMB": 1000,
            "timeThresholdHrs": 24
        }
    }],
    "monitoringVersions": [
        {{range $i, $node := .Nodes}}{{if $i}},{{end}}{
            "hostname": "{{$node}}",
            "logPath": "/var/vcap/sys/log/mongod_node/monitoring-agent.log",
            "logRotate": {
                "sizeThresholdMB": 1000,
                "timeThresholdHrs": 24
            }
        }{{end}}
    ],
    "processes": [{{range $i, $node := .Nodes}}
      {{if $i}},{{end}}{
        "args2_6": {
            "net": {
                "port": 28000
            },
            "replication": {
                "replSetName": "pcf_repl"
            },
            "storage": {
                "dbPath": "/var/vcap/store/mongodb-data"
            },
            "systemLog": {
                "destination": "file",
                "path": "/var/vcap/sys/log/mongod_node/mongodb.log"
            }
        },
        "hostname": "{{$node}}",
        "logRotate": {
            "sizeThresholdMB": 1000,
            "timeThresholdHrs": 24
        },
        "name": "{{$node}}",
        "processType": "mongod",
        "version": "{{$.Version}}",
        "authSchemaVersion": 5
    }
    {{end}}
  ],
    "replicaSets": [{
        "_id": "pcf_repl",
        "members": [
          {{range $i, $node := .Nodes}}
          {{if $i}},{{end}}{
            "_id": {{$i}},
            "host": "{{$node}}"{{if last $.Nodes $i}},
            "arbiterOnly": true,
            "priority": 0
            {{end}}
          }
          {{end}}
        ]
    }],
    "roles": [],
    "sharding": [],

    "auth": {
        "autoUser": "mms-automation",
        "autoPwd": "{{.Password}}",
        "deploymentAuthMechanisms": [
            "SCRAM-SHA-1"
        ],
        "key": "{{.Key}}",
        "keyfile": "/var/vcap/jobs/mongod_node/config/mongo_om.key",
        "disabled": false,
        "usersDeleted": [],
        "usersWanted": [
            {
                "db": "admin",
                "roles": [
                    {
                        "db": "admin",
                        "role": "clusterMonitor"
                    }
                ],
                "user": "mms-monitoring-agent",
                "initPwd": "{{.Password}}"
            },
            {
                "db": "admin",
                "roles": [
                    {
                        "db": "admin",
                        "role": "clusterAdmin"
                    },
                    {
                        "db": "admin",
                        "role": "readAnyDatabase"
                    },
                    {
                        "db": "admin",
                        "role": "userAdminAnyDatabase"
                    },
                    {
                        "db": "local",
                        "role": "readWrite"
                    },
                    {
                        "db": "admin",
                        "role": "readWrite"
                    }
                ],
                "user": "mms-backup-agent",
                "initPwd": "{{.Password}}"
            },
            {
               "db": "admin" ,
               "user": "admin" ,
               "roles": [
                 {
                     "db": "admin",
                     "role": "clusterAdmin"
                 },
                 {
                     "db": "admin",
                     "role": "readAnyDatabase"
                 },
                 {
                     "db": "admin",
                     "role": "userAdminAnyDatabase"
                 },
                 {
                     "db": "local",
                     "role": "readWrite"
                 },
                 {
                     "db": "admin",
                     "role": "readWrite"
                 }
               ],
               "initPwd": "{{.AdminPassword}}"
            }
        ],
        "autoAuthMechanism": "SCRAM-SHA-1"
    }
}`,
}
