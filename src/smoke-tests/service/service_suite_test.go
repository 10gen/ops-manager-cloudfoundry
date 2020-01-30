package service_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/cf-test-helpers/services"

	"smoke-tests/retry"
	"smoke-tests/service/configuration"
	"smoke-tests/service/reporter"
)

var (
	configPath        = os.Getenv("CONFIG_PATH")
	cfTestConfig      = loadCFTestConfig(configPath)
	mongodbConfig     = loadMongodbTestConfig(configPath)
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
