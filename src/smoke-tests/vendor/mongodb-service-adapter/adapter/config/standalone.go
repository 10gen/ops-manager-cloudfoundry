package config

const reference = `{
    "options": {
        "downloadBase": "/var/lib/mongodb-mms-automation"
    },
    {{ if .RequireSSL }}
    "ssl" : {
				"autoPEMKeyFilePath": "/var/vcap/jobs/mongod_node/config/server.pem",
        "CAFilePath": "/var/vcap/jobs/mongod_node/config/cacert.pem",
        "clientCertificateMode": "OPTIONAL"
    },
    {{end}}
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
    "monitoringVersions": [{
        "hostname": "{{index .Nodes 0}}",
        "logPath": "/var/vcap/sys/log/mongod_node/monitoring-agent.log",
        "logRotate": {
            "sizeThresholdMB": 1000,
            "timeThresholdHrs": 24
        }
    }],
    "processes": [{
        "args2_6": {
            "net": {
                {{ if .RequireSSL }}
                "ssl": {
                    "mode": "requireSSL",
                    "PEMKeyFile": "/var/vcap/jobs/mongod_node/config/server.pem"
                },
                {{ end }}
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
        {{if ne .CompatibilityVersion ""}}
            "featureCompatibilityVersion": "{{.CompatibilityVersion}}",
        {{end}}
        "authSchemaVersion": 5
    }],
    "replicaSets": [],
    "roles": [],
    "sharding": [],

    "auth": {
        "autoUser": "mms-automation",   
		"autoPwd": "{{.AutomationAgentPassword}}",
        "deploymentAuthMechanisms": [
            "MONGODB-CR"
        ],
        "key": "{{.Key}}",
        "keyfile": "/var/vcap/store/mongod_node/mongo_om.key",
        {{if ne .KeyfileWindows ""}}
            "keyfileWindows": "{{.KeyfileWindows}}",
        {{end}}
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
        "autoAuthMechanism": "MONGODB-CR"
    }
}`
