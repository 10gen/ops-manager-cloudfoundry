package config

import "github.com/mongodb-labs/pcgc/pkg/opsmanager"

func getDefaultLogRotate() *opsmanager.LogRotate {
	return &opsmanager.LogRotate{
		SizeThresholdMB:  1000,
		TimeThresholdHrs: 24,
	}
}

func getDefaultSSL() *opsmanager.SSL {
	return &opsmanager.SSL{
		AutoPEMKeyFilePath:    "/var/vcap/jobs/mongod_node/config/server.pem",
		CAFilePath:            "/var/vcap/jobs/mongod_node/config/cacert.pem",
		ClientCertificateMode: "OPTIONAL",
	}
}

func getDefaultSystemLog() *opsmanager.SystemLog {
	return &opsmanager.SystemLog{
		Destination: "file",
		Path:        "/var/vcap/sys/log/mongod_node/mongodb.log",
	}
}

func getdefaultNet(ssl bool) *opsmanager.Net {
	r := opsmanager.Net{
		Port: 28000,
	}
	if ssl {
		r.SSL = &opsmanager.NetSSL{
			Mode:       "requireSSL",
			PEMKeyFile: "/var/vcap/jobs/mongod_node/config/server.pem",
		}
	}
	return &r
}

func getDefaultUsers(ctx *DocContext) []opsmanager.UserWanted {
	return []opsmanager.UserWanted{
		{
			DB:      "admin",
			User:    "mms-monitoring-agent",
			InitPwd: ctx.Password,
			Roles: []opsmanager.Role{
				{
					DB:   "admin",
					Role: "clusterMonitor",
				},
			},
		},
		{
			DB:      "admin",
			User:    "mms-backup-agent",
			InitPwd: ctx.Password,
			Roles: []opsmanager.Role{
				{
					DB:   "admin",
					Role: "clusterAdmin",
				},
				{
					DB:   "admin",
					Role: "readAnyDatabase",
				},
				{
					DB:   "admin",
					Role: "userAdminAnyDatabase",
				},
				{
					DB:   "local",
					Role: "readWrite",
				},
				{
					DB:   "admin",
					Role: "readWrite",
				},
			},
		},
		{
			DB:      "admin",
			User:    "admin",
			InitPwd: ctx.AdminPassword,
			Roles: []opsmanager.Role{
				{
					DB:   "admin",
					Role: "clusterAdmin",
				},
				{
					DB:   "admin",
					Role: "readAnyDatabase",
				},
				{
					DB:   "admin",
					Role: "userAdminAnyDatabase",
				},
				{
					DB:   "local",
					Role: "readWrite",
				},
				{
					DB:   "admin",
					Role: "readWrite",
				},
			},
		},
	}
}

func getDefaultVersions(version, name string) ([]*opsmanager.MongoDBVersion, []*opsmanager.AgentVersion, []*opsmanager.AgentVersion) {
	mongoDBVersions := []*opsmanager.MongoDBVersion{{
		Name: version,
	}}

	backupVersions := []*opsmanager.AgentVersion{{
		Hostname:  name,
		LogPath:   "/var/vcap/sys/log/mongod_node/backup-agent.log",
		LogRotate: getDefaultLogRotate(),
	}}

	monitoringVersions := []*opsmanager.AgentVersion{{
		Hostname:  name,
		LogPath:   "/var/vcap/sys/log/mongod_node/monitoring-agent.log",
		LogRotate: getDefaultLogRotate(),
	}}

	return mongoDBVersions, backupVersions, monitoringVersions
}
