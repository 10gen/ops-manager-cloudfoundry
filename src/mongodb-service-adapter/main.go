package main

import (
	"log"
	"os"

	"github.com/emiloserdov/mongodb-on-demand-release/src/mongodb-service-adapter/adapter"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

func main() {
	logger := log.New(os.Stderr, "[mongodb-service-adapter] ", log.LstdFlags)
	manifestGenerator := &adapter.ManifestGenerator{
		Logger: logger,
	}
	binder := &adapter.Binder{
		Logger: logger,
	}
	handler := serviceadapter.CommandLineHandler{
		ManifestGenerator:     manifestGenerator,
		Binder:                binder,
		DashboardURLGenerator: &adapter.DashboardURLGenerator{},
	}
	serviceadapter.HandleCLI(os.Args, handler)
}
