package cf

import (
	"github.com/pivotal-cf-experimental/cf-test-helpers/runner"
	"github.com/onsi/gomega/gexec"
)

var Cf = func(args ...string) *gexec.Session {
	return runner.Run("cf", args...)
}
