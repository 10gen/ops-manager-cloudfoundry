package opsmanager

// Link holder for links returned by the API
type Link struct {
	HREF string `json:"href"`
	Rel  string `json:"rel"`
}
