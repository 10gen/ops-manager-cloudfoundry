package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"mongodb-service-adapter/adapter"
	"mongodb-service-adapter/adapter/config"
)

var (
	configFilePath string
)

func main() {
	flag.StringVar(&configFilePath, "config", "", "Location of the config file")
	flag.Parse()

	cfg, err := LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Error loading config file: %s", err)
	}

	logger := log.New(os.Stderr, "[mongodb-config-agent] ", log.LstdFlags)
	omClient := adapter.OMClient{Url: cfg.URL, Username: cfg.Username, ApiKey: cfg.APIKey}

	automation, err := omClient.Client().GetAutomationConfig(cfg.GroupID)
	if err != nil {
		logger.Fatal(err)
	}

	nodes := strings.Split(cfg.NodeAddresses, ",")
	ctx := &config.DocContext{
		ID:                      cfg.ID,
		Key:                     cfg.AuthKey,
		AdminPassword:           cfg.AdminPassword,
		AutomationAgentPassword: cfg.AutomationAgentPassword,
		Nodes:                   nodes,
		Version:                 cfg.EngineVersion,
		RequireSSL:              cfg.RequireSSL,
	}

	if cfg.PlanID == adapter.PlanShardedCluster {
		var err error
		ctx.Cluster, err = adapter.NodesToCluster(nodes, cfg.Routers, cfg.ConfigServers, cfg.Replicas)
		if err != nil {
			logger.Fatal(err)
		}
	}

	if ctx.Password == "" {
		var err error
		ctx.Password, err = adapter.GenerateString(32)
		if err != nil {
			logger.Fatal(err)
		}
	}

	if strings.HasPrefix(ctx.Version, "3.4") {
		ctx.CompatibilityVersion = "3.4"
	} else if strings.HasPrefix(ctx.Version, "3.6") {
		ctx.CompatibilityVersion = "3.6"
	} else if strings.HasPrefix(ctx.Version, "4.0") {
		ctx.CompatibilityVersion = "4.0"
	}

	automation, err = omClient.Client().UpdateAutomationConfig(cfg.GroupID, automation)
	if err != nil {
		logger.Fatal(err)
	}

	doc, err := json.Marshal(automation)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(string(doc))

	err = omClient.Client().UpdateMonitoringConfig(cfg.GroupID, config.GetMonitoringAgentConfiguration(ctx))
	if err != nil {
		logger.Fatal(err)
	}

	err = omClient.Client().UpdateBackupConfig(cfg.GroupID, config.GetBackupAgentConfiguration(ctx))
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("Configured group %s", cfg.GroupID)

	for {
		logger.Printf("Checking group %s", cfg.GroupID)
		groupHosts, err := omClient.Client().GetHosts(cfg.GroupID)
		if err != nil {
			logger.Fatal(err)
		}
		if groupHosts.TotalCount == 0 {
			logger.Fatalf("Host count for %s is 0...", cfg.GroupID)
		}
		time.Sleep(30 * time.Second)
	}
}
