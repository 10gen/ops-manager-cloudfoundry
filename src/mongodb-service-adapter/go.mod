module mongodb-service-adapter

go 1.12

replace github.com/mongodb-labs/pcgc v0.0.3 => github.com/vasilevp/pcgc v0.0.5

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible // indirect
	github.com/mongodb-labs/pcgc v0.0.3
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pivotal-cf/on-demand-services-sdk v0.35.0
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/tidwall/gjson v1.3.2
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)
