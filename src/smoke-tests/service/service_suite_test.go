package service_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"smoke-tests/cf"
	prepare "smoke-tests/service/configuration"
	"smoke-tests/service/reporter"
)

var (
	configPath        = os.Getenv("CONFIG_PATH")
	cfTestConfig      = cf.LoadCFTestConfig(configPath)
	mongodbConfig     = prepare.LoadMongodbTestConfig(configPath)
	smokeTestReporter *reporter.SmokeTestReport
)

func TestService(t *testing.T) {
	smokeTestReporter = new(reporter.SmokeTestReport)

	reporter := []Reporter{
		Reporter(smokeTestReporter),
	}

	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "MongoDB Smoke Tests", reporter)
}
