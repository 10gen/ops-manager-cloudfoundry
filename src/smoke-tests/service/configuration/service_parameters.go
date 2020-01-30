package configuration

import (
	"fmt"
	"strings"
)

//ServiceParameters supported parameter list for creating\updating services from CF
type ServiceParameters struct {
	ServiceName    string //what? rename it. to serviteTestName maybe
	PlanName       string
	orgID          string
	BackupEnable   string
	SSLEnable      string
	MongoDBVersion string
	Replicas       string
	Shards         string
	ConfigServers  string
	Mongos         string
	BoshDNSDisable string //says depricated?
}

//PrintParameters is printing parameters used in tests
func (sp ServiceParameters) PrintParameters() string {
	return fmt.Sprintf("%s (Backup: %s, SSL: %s, version: %s)", strings.ToUpper(sp.PlanName), sp.BackupEnable, sp.SSLEnable, sp.MongoDBVersion)
}
