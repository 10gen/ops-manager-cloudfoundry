module mongodb-config-agent

go 1.12

replace mongodb-service-adapter => ../mongodb-service-adapter

replace github.com/mongodb-labs/pcgc v0.0.3 => github.com/vasilevp/pcgc v0.0.8

require (
	github.com/pkg/errors v0.8.1
	golang.org/x/mod v0.2.0
	mongodb-service-adapter v0.0.0-00010101000000-000000000000
)
