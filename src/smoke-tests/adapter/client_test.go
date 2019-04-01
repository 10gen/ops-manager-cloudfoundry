package adapter_test

import (
	adapter "../../mongodb-service-adapter/adapter"
	"encoding/json"
	"testing"
)

var latestVersion = "4.0.7"

var GetConfig = func() *adapter.Config {
	config, err := adapter.LoadConfig("../../mongodb-service-adapter/testdata/manifest.json")

	if err != nil {
		return nil
	}

	return config
}

func TestLoadDoc(t *testing.T) {
	t.Parallel()
	c := &adapter.OMClient{}

	for p, ctx := range map[string]*adapter.DocContext{
		adapter.PlanStandalone: {
			ID:            "d3a98gf1",
			Key:           "key",
			AdminPassword: "pwd",
			Nodes:         []string{"192.168.1.1", "192.168.1.2", "192.168.1.3", "192.168.1.4"},
			Version:       "3.2.11",
		},

		adapter.PlanShardedCluster: {
			ID:            "d3a98gf1",
			Key:           "key",
			AdminPassword: "pwd",
			Version:       "3.4.1",
			Cluster: &adapter.Cluster{
				Routers:       []string{"192.168.1.1", "192.168.1.2"},
				ConfigServers: []string{"192.168.1.3", "192.168.1.4"},
				Shards: [][]string{
					{"192.168.0.10", "192.168.0.11", "192.168.0.12"},
					{"192.168.0.13", "192.168.0.14", "192.168.0.15"},
				},
			},
		},

		adapter.PlanReplicaSet: {
			ID:            "d3a98gf1",
			Key:           "key",
			AdminPassword: "pwd",
			Nodes:         []string{"192.168.1.1", "192.168.1.2", "192.168.1.3", "192.168.1.4"},
			Version:       "3.2.11",
		},
	} {
		t.Run(string(p), func(t *testing.T) {
			s, err := c.LoadDoc(p, ctx)
			if err != nil {
				t.Fatal(err)
			}

			// validate json output
			if err := json.Unmarshal([]byte(s), &map[string]interface{}{}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetGroupByName(t *testing.T) {

	c := GetOMClient()

	_, err := c.GetGroupByName("Project 2")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateGroup(t *testing.T) {

	c := GetOMClient()

	request := adapter.GroupCreateRequest{
		Name: "Project 2",
		Tags: []string{"tag1", "tag2"},
	}

	_, err := c.CreateGroup("abcdefgh123id", request)

	if err != nil {
		t.Fatal(err)
	}

}

func TestCreateGroupAPIKey(t *testing.T) {

	c := GetOMClient()

	_, err := c.CreateGroupAPIKey(GetConfig().GroupID)

	if err != nil {
		t.Fatal(err)
	}

}

func TestUpdateGroup(t *testing.T) {

	c := GetOMClient()

	request := adapter.GroupUpdateRequest{
		Tags: []string{
			"tag1",
			"tag2",
		},
	}

	_, err := c.UpdateGroup(GetConfig().GroupID, request)

	if err != nil {
		t.Fatal(err)
	}

}

func TestGetGroupAuthAgentPassword(t *testing.T) {

	c := GetOMClient()

	_, err := c.GetGroupAuthAgentPassword(GetConfig().GroupID)

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetGroup(t *testing.T) {

	c := GetOMClient()

	_, err := c.GetGroup(GetConfig().GroupID)

	if err != nil {
		t.Fatal(err)
	}

}

func TestDeleteGroup(t *testing.T) {

	c := GetOMClient()

	err := c.DeleteGroup("test-group")

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetGroupHosts(t *testing.T) {

	c := GetOMClient()

	_, err := c.GetGroupHosts(GetConfig().GroupID)

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetGroupHostnames(t *testing.T) {

	c := GetOMClient()

	_, err := c.GetGroupHostnames(GetConfig().GroupID, "standalone")

	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigureGroup(t *testing.T) {

	c := GetOMClient()

	ctx := &adapter.DocContext{
		ID:                      "id1",
		Key:                     GetConfig().AuthKey,
		AdminPassword:           "admin",
		AutomationAgentPassword: "admin",
		Nodes:                   []string{GetConfig().NodeAddresses},
		Version:                 "4.0.7",
		RequireSSL:              false,
	}

	doc, err := c.LoadDoc(adapter.PlanStandalone, ctx)

	if err != nil {
		t.Fatal(err)
	}

	err = c.ConfigureGroup(doc, GetConfig().GroupID)

	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigureMonitoringAgent(t *testing.T) {

	c := GetOMClient()

	ctx := &adapter.DocContext{
		ID:                      "id1",
		Key:                     GetConfig().AuthKey,
		AdminPassword:           "admin",
		AutomationAgentPassword: "admin",
		Nodes:                   []string{GetConfig().NodeAddresses},
		Version:                 "4.0.7",
		RequireSSL:              false,
	}

	monitoringAgentDoc, err := c.LoadDoc(adapter.MonitoringAgentConfiguration, ctx)

	if err != nil {
		t.Fatal(err)
	}

	err = c.ConfigureMonitoringAgent(monitoringAgentDoc, GetConfig().GroupID)

	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigureBackupAgent(t *testing.T) {

	c := GetOMClient()

	ctx := &adapter.DocContext{
		ID:                      "id1",
		Key:                     GetConfig().AuthKey,
		AdminPassword:           "admin",
		AutomationAgentPassword: "admin",
		Nodes:                   []string{GetConfig().NodeAddresses},
		Version:                 "4.0.7",
		RequireSSL:              false,
	}

	backupAgentDoc, err := c.LoadDoc(adapter.BackupAgentConfiguration, ctx)

	if err != nil {
		t.Fatal(err)
	}

	err = c.ConfigureBackupAgent(backupAgentDoc, GetConfig().GroupID)

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAvailableVersions(t *testing.T) {

	c := GetOMClient()

	_, err := c.GetAvailableVersions(GetConfig().GroupID)

	if err != nil {
		t.Fatal()
	}
}

func TestGetLatestVersion(t *testing.T) {

	c := GetOMClient()

	_, err := c.GetLatestVersion(GetConfig().GroupID)

	if err != nil {
		t.Fatal()
	}

}

func TestValidateVersion(t *testing.T) {

	c := GetOMClient()

	_, err := c.ValidateVersion(GetConfig().GroupID, latestVersion)

	if err != nil {
		t.Fatal()
	}
}

func TestValidateVersionManifest(t *testing.T) {

	c := GetOMClient()

	_, err := c.ValidateVersionManifest(latestVersion)

	if err != nil {
		t.Fatal(err)
	}
}

func GetOMClient() *adapter.OMClient {

	config := GetConfig()

	c := &adapter.OMClient{
		Url:      config.URL,
		Username: config.Username,
		ApiKey:   config.APIKey,
	}

	return c
}
