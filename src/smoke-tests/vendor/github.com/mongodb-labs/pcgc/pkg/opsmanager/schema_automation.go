package opsmanager

// AutomationConfig represents a cluster definition within an automation config object
// NOTE: this struct is mutable
type AutomationConfig struct {
	Auth               Auth                     `json:"auth,omitempty"`
	LDAP               map[string]interface{}   `json:"ldap,omitempty"`
	Processes          []*Process               `json:"processes,omitempty"`
	ReplicaSets        []ReplicaSet             `json:"replicaSets,omitempty"`
	Roles              []map[string]interface{} `json:"roles,omitempty"`
	MonitoringVersions []*AgentVersion          `json:"monitoringVersions,omitempty"`
	BackupVersions     []*AgentVersion          `json:"backupVersions,omitempty"`
	MongoSQLDs         []map[string]interface{} `json:"mongosqlds,omitempty"`
	MongoDBVersions    []*MongoDBVersion        `json:"mongoDbVersions,omitempty"`
	AgentVersion       map[string]interface{}   `json:"agentVersion,omitempty"`
	Balancer           map[string]interface{}   `json:"balancer,omitempty"`
	CPSModules         []map[string]interface{} `json:"cpsModules,omitempty"`
	IndexConfigs       []map[string]interface{} `json:"indexConfigs,omitempty"`
	Kerberos           map[string]interface{}   `json:"kerberos,omitempty"`
	MongoTs            []map[string]interface{} `json:"mongots,omitempty"`
	Options            *Options                 `json:"options"`
	SSL                *SSL                     `json:"ssl,omitempty"`
	Version            int                      `json:"version,omitempty"`
	Sharding           []Sharding               `json:"sharding,omitempty"`
	UIBaseURL          string                   `json:"uiBaseUrl,omitempty"`
}

// AgentVersion agent versions struct
type AgentVersion struct {
	Name      string     `json:"name,omitempty"`
	Hostname  string     `json:"hostname"`
	LogPath   string     `json:"logPath,omitempty"`
	LogRotate *LogRotate `json:"logRotate,omitempty"`
}

// SSL ssl config properties
type SSL struct {
	AutoPEMKeyFilePath    string `json:"autoPEMKeyFilePath"`
	CAFilePath            string `json:"CAFilePath"`
	ClientCertificateMode string `json:"clientCertificateMode"`
}

// Auth authentication config
type Auth struct {
	AutoUser                 string        `json:"autoUser"`
	AutoPwd                  string        `json:"autoPwd"`
	DeploymentAuthMechanisms []string      `json:"deploymentAuthMechanisms"`
	Key                      string        `json:"key"`
	Keyfile                  string        `json:"keyfile"`
	KeyfileWindows           string        `json:"keyfileWindows"`
	Disabled                 bool          `json:"disabled"`
	UsersDeleted             []interface{} `json:"usersDeleted"`
	UsersWanted              []UserWanted  `json:"usersWanted"`
	AutoAuthMechanism        string        `json:"autoAuthMechanism"`
}

// UserWanted :shrug:
type UserWanted struct {
	DB      string `json:"db"`
	Roles   []Role `json:"roles"`
	User    string `json:"user"`
	InitPwd string `json:"initPwd"`
}

// Role user role
type Role struct {
	DB   string `json:"db"`
	Role string `json:"role"`
}

// Member configs
type Member struct {
	ID          int    `json:"_id"`
	ArbiterOnly bool   `json:"arbiterOnly"`
	Hidden      bool   `json:"hidden"`
	Host        string `json:"host"`
	Priority    int    `json:"priority"`
	SlaveDelay  int    `json:"slaveDelay"`
	Votes       int    `json:"votes"`
}

// ReplicaSet configs
type ReplicaSet struct {
	ID              string   `json:"_id"`
	ProtocolVersion int      `json:"protocolVersion,omitempty"`
	Members         []Member `json:"members"`
}

// Options configs
type Options struct {
	DownloadBase string `json:"downloadBase"`
}

// Sharding configs
type Sharding struct {
	Shards              []Shard       `json:"shards"`
	Name                string        `json:"name"`
	ConfigServer        []interface{} `json:"configServer"`
	ConfigServerReplica string        `json:"configServerReplica"`
	Collections         []interface{} `json:"collections"`
}

// Shard configs
type Shard struct {
	Tags []interface{} `json:"tags"`
	ID   string        `json:"_id"`
	Rs   string        `json:"rs"`
}
