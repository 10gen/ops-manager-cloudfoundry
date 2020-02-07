package data

import ()

var (
	//AppPath path to ruby simple application
	AppPath = "../assets/cf-mongo-example-app"
	//CreateQuotaArgs quota for services
	CreateQuotaArgs = []string{
		"-m", "10G",
		"-r", "1000",
		"-s", "100",
		"--allow-paid-service-plans",
	}
	//PushArgs for pushing arguments
	PushArgs = []string{
		"-m", "256M",
		"-p", AppPath,
		"-s", "cflinuxfs3",
		"--no-start",
	}
)
