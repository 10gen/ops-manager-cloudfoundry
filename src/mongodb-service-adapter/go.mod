module mongodb-service-adapter

go 1.12

replace github.com/mongodb-labs/pcgc => ../pcgc

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible // indirect
	github.com/drewolson/testflight v1.0.0 // indirect
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/mongodb-labs/pcgc v0.0.0-00010101000000-000000000000
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pivotal-cf/brokerapi v6.0.0+incompatible // indirect
	github.com/pivotal-cf/on-demand-services-sdk v0.31.1
	github.com/pkg/errors v0.8.1 // indirect
	github.com/stretchr/testify v1.3.0 // indirect
	github.com/tidwall/gjson v1.3.2
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
)
