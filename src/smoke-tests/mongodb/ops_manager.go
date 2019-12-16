package mongodb

import (

)

type OpsManager struct {
	BaseURL string
}

type Service struct {
	Name          string
	Service       string
	Tags          string
	Plan          string
	Description   string
	Documentation string
	Dashboard     string
	ServiceBroker string
	Status        string
	Message       string
	Started       string
	Update        string
}
