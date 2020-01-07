module smoke-tests

go 1.12

replace github.com/mongodb-labs/pcgc v0.0.3 => github.com/vasilevp/pcgc v0.0.5

replace mongodb-service-adapter => ../mongodb-service-adapter

require (
	github.com/onsi/ginkgo v1.10.3
	github.com/onsi/gomega v1.7.1
	github.com/pborman/uuid v1.2.0
	github.com/pivotal-cf-experimental/cf-test-helpers v0.0.0-20170428144005-e56b6ec41da9
	github.com/pivotal-cf/on-demand-services-sdk v0.35.0
	mongodb-service-adapter v0.0.0-00010101000000-000000000000
)
