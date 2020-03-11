package configuration

//OpsManagerCred used for connection to opsmanager
type OpsManagerCred struct {
	URL        string `json:"url"`
	UserName   string `json:"username"`
	UserAPIKey string `json:"api_key"`
}
