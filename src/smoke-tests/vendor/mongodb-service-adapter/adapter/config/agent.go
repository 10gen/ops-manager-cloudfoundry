package config

import "github.com/mongodb-labs/pcgc/pkg/opsmanager"

var defaultAgentConfig = opsmanager.AgentAttributes{
	SSLPEMKeyfile: "/var/vcap/jobs/mongod_node/config/server.pem",
	LogRotate:     getDefaultLogRotate(),
}

func GetMonitoringAgentConfiguration(ctx *DocContext) opsmanager.AgentAttributes {
	r := defaultAgentConfig
	r.LogPath = "/var/vcap/sys/log/mongod_node/monitoring-agent.log"

	if !ctx.RequireSSL {
		r.SSLPEMKeyfile = ""
	}

	return r
}

func GetBackupAgentConfiguration(ctx *DocContext) opsmanager.AgentAttributes {
	r := defaultAgentConfig
	r.LogPath = "/var/vcap/sys/log/mongod_node/backup-agent.log"

	if !ctx.RequireSSL {
		r.SSLPEMKeyfile = ""
	}

	return r
}
