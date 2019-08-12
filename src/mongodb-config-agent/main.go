package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"mongodb-service-adapter/adapter"
)

var (
	configFilePath string
)

func main() {
	flag.StringVar(&configFilePath, "config", "", "Location of the config file")
	flag.Parse()

	config, err := LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Error loading config file: %s", err)
	}

	logger := log.New(os.Stderr, "[mongodb-config-agent] ", log.LstdFlags)
	omClient := adapter.OMClient{Url: config.URL, Username: config.Username, ApiKey: config.APIKey}

	kfw, err := omClient.GetKeyfileWindows(config.GroupID)
	if err != nil {
		logger.Fatal(err)
	}

	nodes := strings.Split(config.NodeAddresses, ",")
	ctx := &adapter.DocContext{
		ID:                      config.ID,
		Key:                     config.AuthKey,
		AdminPassword:           config.AdminPassword,
		AutomationAgentPassword: config.AutomationAgentPassword,
		Nodes:                   nodes,
		Version:                 config.EngineVersion,
		RequireSSL:              config.RequireSSL,
		KeyfileWindows:          kfw,
	}

	if config.PlanID == adapter.PlanShardedCluster {
		var err error
		ctx.Cluster, err = adapter.NodesToCluster(nodes, config.Routers, config.ConfigServers, config.Replicas)
		if err != nil {
			logger.Fatal(err)
		}
	}

	logger.Printf("%+v", nodes)
	doc, err := omClient.LoadDoc(config.PlanID, ctx)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(doc)

	monitoringAgentDoc, err := omClient.LoadDoc(adapter.MonitoringAgentConfiguration, ctx)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(monitoringAgentDoc)

	backupAgentDoc, err := omClient.LoadDoc(adapter.MonitoringAgentConfiguration, ctx)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(backupAgentDoc)

	err = omClient.ConfigureGroup(doc, config.GroupID)
	if err != nil {
		logger.Fatal(err)
	}

	err = omClient.ConfigureMonitoringAgent(monitoringAgentDoc, config.GroupID)
	if err != nil {
		logger.Fatal(err)
	}

	err = omClient.ConfigureBackupAgent(backupAgentDoc, config.GroupID)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("Configured group %s", config.GroupID)

	for {
		logger.Printf("Checking group %s", config.GroupID)
		groupHosts, err := omClient.GetGroupHosts(config.GroupID)
		if err != nil {
			logger.Fatal(err)
		}
		if groupHosts.TotalCount == 0 {
			logger.Fatalf("Host count for %s is 0...", config.GroupID)
		}
		time.Sleep(30 * time.Second)
	}
}
