package opsmanager

// NetSSL defines SSL parameters for Net
type NetSSL struct {
	Mode       string `json:"mode"`
	PEMKeyFile string `json:"PEMKeyFile"`
}

// Net part of the internal Process struct
type Net struct {
	Port int     `json:"port,omitempty"`
	SSL  *NetSSL `json:"ssl,omitempty"`
}

// StorageArg part of the internal Process struct
type StorageArg struct {
	DBPath string `json:"dbPath,omitempty"`
}

// ReplicationArg is part of the internal Process struct
type ReplicationArg struct {
	ReplSetName string `json:"replSetName"`
}

// ShardingArg is part of the internal Process struct
type ShardingArg struct {
	ClusterRole string `json:"clusterRole"`
}

// SystemLog part of the internal Process struct
type SystemLog struct {
	Destination string `json:"destination,omitempty"`
	Path        string `json:"path,omitempty"`
}

// Args26 part of the internal Process struct
type Args26 struct {
	NET         *Net            `json:"net,omitempty"`
	Storage     *StorageArg     `json:"storage,omitempty"`
	SystemLog   *SystemLog      `json:"systemLog,omitempty"`
	Replication *ReplicationArg `json:"replication,omitempty"`
	Sharding    *ShardingArg    `json:"sharding,omitempty"`
}

// LogRotate part of the internal Process struct
type LogRotate struct {
	SizeThresholdMB  float64 `json:"sizeThresholdMB,omitempty"`
	TimeThresholdHrs int     `json:"timeThresholdHrs,omitempty"`
}

// Process represents a single process in a deployment
type Process struct {
	Name                        string     `json:"name,omitempty"`
	ProcessType                 string     `json:"processType,omitempty"`
	Version                     string     `json:"version,omitempty"`
	AuthSchemaVersion           int        `json:"authSchemaVersion,omitempty"`
	FeatureCompatibilityVersion string     `json:"featureCompatibilityVersion,omitempty"`
	Disabled                    bool       `json:"disabled,omitempty"`
	ManualMode                  bool       `json:"manualMode,omitempty"`
	Hostname                    string     `json:"hostname,omitempty"`
	Args26                      *Args26    `json:"args2_6,omitempty"`
	LogRotate                   *LogRotate `json:"logRotate,omitempty"`
	Plan                        []string   `json:"plan,omitempty"`
	LastGoalVersionAchieved     int        `json:"lastGoalVersionAchieved,omitempty"`
	Cluster                     string     `json:"cluster,omitempty"`
}
