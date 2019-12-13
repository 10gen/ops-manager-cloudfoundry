module mongodb-config-agent

go 1.12

replace mongodb-service-adapter => ../mongodb-service-adapter

replace github.com/mongodb-labs/pcgc v0.0.3 => /home/exploser/go/src/github.com/vasilevp/pcgc

require mongodb-service-adapter v0.0.0-00010101000000-000000000000
