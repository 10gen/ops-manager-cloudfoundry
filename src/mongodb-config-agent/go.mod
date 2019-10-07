module mongodb-config-agent

go 1.12

replace mongodb-service-adapter => ../mongodb-service-adapter

replace github.com/mongodb-labs/pcgc => ../pcgc

require (
	github.com/drewolson/testflight v1.0.0 // indirect
	github.com/onsi/ginkgo v1.10.2 // indirect
	github.com/onsi/gomega v1.7.0 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	mongodb-service-adapter v0.0.0-00010101000000-000000000000
)
