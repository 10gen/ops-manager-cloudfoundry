module mongodb-config-agent

go 1.12

replace mongodb-service-adapter => ../mongodb-service-adapter

require (
	github.com/mongodb-labs/pcgc v0.0.3
	mongodb-service-adapter v0.0.0-00010101000000-000000000000
)