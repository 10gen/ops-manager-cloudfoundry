package configuration

import (
	"fmt"
	"strings"

	"github.com/pborman/uuid"
)

//ServiceParameters supported parameter list for creating\updating services from CF
type ServiceParameters struct {
	Instance               InstanceNames
	MarketPlaceServiceName string
	PlanName               string
	orgID                  string
	BackupEnable           string
	SSLEnable              string
	MongoDBVersion         string
	Replicas               string
	Shards                 string
	ConfigServers          string
	Mongos                 string
	BoshDNSDisable         string //TODO depricated?
}

//InstanceNames instance coordinates (names, key)
type InstanceNames struct {
	AppName           string
	ServiceTestName   string
	SecurityGroupName string
	ServiceKeyName    string
	Plan              string
}

//PrintParameters is printing parameters used in tests
func (sp ServiceParameters) PrintParameters() string {
	return fmt.Sprintf("%s (Backup: %s, SSL: %s, version: %s)", strings.ToUpper(sp.PlanName), sp.BackupEnable, sp.SSLEnable, sp.MongoDBVersion)
}

//GenerateInstanceNames create random name for instance
func (i *InstanceNames) GenerateInstanceNames() {
	i.AppName = RandomName()
	i.ServiceTestName = RandomName()
	i.SecurityGroupName = RandomName()
	i.ServiceKeyName = RandomName()
}

//RandomName gen new name
func RandomName() string {
	return uuid.NewRandom().String()
}
