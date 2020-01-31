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

	omClient := adapter.NewPCGCClient(cfg.URL, cfg.Username, cfg.APIKey)

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
		ctx.Password = adapter.GenerateString(32, true)
	}

	switch c := ctx.Version[:3]; c {
	case "3.4", "3.6", "4.0":
		ctx.CompatibilityVersion = c
	}

	ac, err := omClient.GetAutomationConfig(cfg.GroupID)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Println("Original AutomationConfig:")
	printConfig(logger, ac)

	automation := (*config.AutomationConfig)(&ac)
	switch cfg.PlanID {
	case adapter.PlanStandalone:
		automation.ToStandalone(ctx)

	case adapter.PlanReplicaSet:
		automation.ToReplicaSet(ctx)

	case adapter.PlanShardedCluster:
		automation.ToShardedCluster(ctx)

	default:
		logger.Fatalf("unknown plan ID %q", cfg.PlanID)
	}

	logger.Println("Modified AutomationConfig:")
	printConfig(logger, ac)

	_, err = omClient.UpdateAutomationConfig(cfg.GroupID, ac)
	if err != nil {
		logger.Fatal(err)
	}

	err = omClient.UpdateMonitoringConfig(cfg.GroupID, config.GetMonitoringAgentConfiguration(ctx))
	if err != nil {
		logger.Fatal(err)
	}

	err = omClient.UpdateBackupConfig(cfg.GroupID, config.GetBackupAgentConfiguration(ctx))
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("Configured group %s", cfg.GroupID)

	for {
		logger.Printf("Checking group %s", cfg.GroupID)
		groupHosts, err := omClient.GetHosts(cfg.GroupID)
		if err != nil {
			logger.Fatal(err)
		}
		if groupHosts.TotalCount == 0 {
			logger.Fatalf("Host count for %s is 0...", cfg.GroupID)
		}
		time.Sleep(30 * time.Second)
	}
}

func printConfig(logger *log.Logger, cfg interface{}) {
	data, err := json.Marshal(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Print(string(data))
}
