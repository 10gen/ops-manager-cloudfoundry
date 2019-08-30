module smoke-tests

replace github.com/pivotal-cf/on-demand-services-sdk => ./vendor/github.com/pivotal-cf/on-demand-services-sdk

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible // indirect
	github.com/10gen/ops-manager-cloudfoundry v0.0.0-20190827134617-d3bda43b69ca
	github.com/gorilla/mux v1.7.3
	github.com/onsi/ginkgo v1.5.0
	github.com/onsi/gomega v1.4.0
	github.com/pborman/uuid v0.0.0-20170612153648-e790cca94e6c
	github.com/pivotal-cf-experimental/cf-test-helpers v0.0.0-20170428144005-e56b6ec41da9
	github.com/pivotal-cf/brokerapi v5.1.0+incompatible
	github.com/pivotal-cf/on-demand-services-sdk v0.31.0
	github.com/pkg/errors v0.8.1
	github.com/tidwall/gjson v1.3.2
	github.com/tidwall/match v1.0.1
	github.com/tidwall/pretty v1.0.0
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)
