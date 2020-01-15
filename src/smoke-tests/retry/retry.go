package retry

import (
	"fmt"
	"math"
	"regexp"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega/gexec"
)

type Check struct {
	sessionProvider sessionProvider
	sessionTimeout  time.Duration
	failHandler     failHandler
	backoff         Backoff
	maxRetries      int
}

func Session(sp sessionProvider) *Check {
	return &Check{
		sessionProvider: sp,
		sessionTimeout:  time.Second,
		failHandler:     ginkgo.Fail,
		backoff:         None(time.Second),
		maxRetries:      10,
	}
}

func (rc *Check) WithFailHandler(handler failHandler) *Check {
	rc.failHandler = handler
	return rc
}

func (rc *Check) AndFailHandler(handler failHandler) *Check {
	return rc.WithFailHandler(handler)
}

func (rc *Check) WithSessionTimeout(timeout time.Duration) *Check {
	rc.sessionTimeout = timeout
	return rc
}

func (rc *Check) AndSessionTimeout(timeout time.Duration) *Check {
	return rc.WithSessionTimeout(timeout)
}

func (rc *Check) WithMaxRetries(max int) *Check {
	rc.maxRetries = max
	return rc
}

func (rc *Check) AndMaxRetries(max int) *Check {
	return rc.WithMaxRetries(max)
}

func (rc *Check) WithBackoff(b Backoff) *Check {
	rc.backoff = b
	return rc
}

func (rc *Check) AndBackoff(b Backoff) *Check {
	return rc.WithBackoff(b)
}

func (rc *Check) Until(c Condition, msg ...string) {
	if rc.check(c) {
		return
	}

	if len(msg) == 0 {
		msg = []string{fmt.Sprintf("Exceeded %d retries", rc.maxRetries)}
	}

	rc.failHandler(msg[0])
}

func (rc *Check) UntilAny(c []Condition, msg ...string) {
	if len(c) < 1 {
		rc.failHandler("Provide at least one condition to match")
		return
	}

	if rc.checkAny(c...) {
		return
	}

	if len(msg) == 0 {
		msg = []string{fmt.Sprintf("Exceeded %d retries", rc.maxRetries)}
	}

	rc.failHandler(msg[0])
}

func (rc *Check) UntilAll(c []Condition, msg ...string) {
	if len(c) < 1 {
		rc.failHandler("Provide at least one condition to match")
		return
	}

	if rc.checkAll(c...) {
		return
	}

	if len(msg) == 0 {
		msg = []string{fmt.Sprintf("Exceeded %d retries", rc.maxRetries)}
	}

	rc.failHandler(msg[0])
}

func (rc *Check) check(c Condition) bool {
	for retry := 0; retry <= rc.maxRetries; retry++ {
		time.Sleep(rc.backoff(uint(retry)))

		session := rc.sessionProvider().Wait(rc.sessionTimeout)

		if c(session) {
			return true
		}
	}

	return false
}

func (rc *Check) checkAny(conditions ...Condition) bool {
	for retry := 0; retry <= rc.maxRetries; retry++ {
		time.Sleep(rc.backoff(uint(retry)))

		session := rc.sessionProvider().Wait(rc.sessionTimeout)

		for _, condition := range conditions {
			if condition(session) {
				return true
			}
		}
	}
	return false
}

func (rc *Check) checkAll(conditions ...Condition) bool {
RetryLoop:
	for retry := 0; retry <= rc.maxRetries; retry++ {
		time.Sleep(rc.backoff(uint(retry)))

		session := rc.sessionProvider().Wait(rc.sessionTimeout)

		for _, condition := range conditions {
			if !condition(session) {
				continue RetryLoop
			}
		}
		return true
	}
	return false
}

type Condition func(session *gexec.Session) bool

func Succeeds(session *gexec.Session) bool {
	return session.ExitCode() == 0
}

func MatchesOutput(regex *regexp.Regexp) Condition {
	return func(session *gexec.Session) bool {
		return regex.Match(session.Out.Contents())
	}
}

func MatchesErrorOutput(regex *regexp.Regexp) Condition {
	return func(session *gexec.Session) bool {
		return regex.Match(session.Err.Contents())
	}
}

func MatchesStdOrErrorOutput(regex *regexp.Regexp) Condition {
	return func(session *gexec.Session) bool {
		return regex.Match(session.Out.Contents()) || regex.Match(session.Err.Contents())
	}
}

type Backoff func(retryCount uint) time.Duration

func None(baseline time.Duration) Backoff {
	return func(retryCount uint) time.Duration {
		if retryCount == 0 {
			return 0
		}

		return baseline
	}
}

func Linear(baseline time.Duration) Backoff {
	return func(retryCount uint) time.Duration {
		return time.Duration(retryCount) * baseline
	}
}

func Exponential(baseline time.Duration) Backoff {
	return func(retryCount uint) time.Duration {
		if retryCount == 0 {
			return 0
		}

		return time.Duration(math.Pow(2, float64(retryCount))) * baseline
	}
}

type sessionProvider func() *gexec.Session

type failHandler func(string, ...int)
