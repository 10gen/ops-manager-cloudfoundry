package retry

import (
	"fmt"
	"math"
	"regexp"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega/gexec"
)

type retryCheck struct {
	sessionProvider sessionProvider
	sessionTimeout  time.Duration
	failHandler     failHandler
	backoff         Backoff
	maxRetries      int
}

func Session(sp sessionProvider) *retryCheck {
	return &retryCheck{
		sessionProvider: sp,
		sessionTimeout:  time.Second,
		failHandler:     ginkgo.Fail,
		backoff:         None(time.Second),
		maxRetries:      10,
	}
}

func (rc *retryCheck) WithFailHandler(handler failHandler) *retryCheck {
	rc.failHandler = handler
	return rc
}

func (rc *retryCheck) AndFailHandler(handler failHandler) *retryCheck {
	return rc.WithFailHandler(handler)
}

func (rc *retryCheck) WithSessionTimeout(timeout time.Duration) *retryCheck {
	rc.sessionTimeout = timeout
	return rc
}

func (rc *retryCheck) AndSessionTimeout(timeout time.Duration) *retryCheck {
	return rc.WithSessionTimeout(timeout)
}

func (rc *retryCheck) WithMaxRetries(max int) *retryCheck {
	rc.maxRetries = max
	return rc
}

func (rc *retryCheck) AndMaxRetries(max int) *retryCheck {
	return rc.WithMaxRetries(max)
}

func (rc *retryCheck) WithBackoff(b Backoff) *retryCheck {
	rc.backoff = b
	return rc
}

func (rc *retryCheck) AndBackoff(b Backoff) *retryCheck {
	return rc.WithBackoff(b)
}

func (rc *retryCheck) Until(c Condition, msg ...string) {
	if rc.check(c) {
		return
	}

	if len(msg) == 0 {
		msg = []string{fmt.Sprintf("Exceeded %d retries", rc.maxRetries)}
	}

	rc.failHandler(msg[0])
}

func (rc *retryCheck) UntilAny(c []Condition, msg ...string) {
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

func (rc *retryCheck) UntilAll(c []Condition, msg ...string) {
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

func (rc *retryCheck) check(c Condition) bool {
	for retry := 0; retry <= rc.maxRetries; retry++ {
		time.Sleep(rc.backoff(uint(retry)))

		session := rc.sessionProvider().Wait(rc.sessionTimeout)

		if c(session) {
			return true
		}
	}

	return false
}

func (rc *retryCheck) checkAny(conditions ...Condition) bool {
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

func (rc *retryCheck) checkAll(conditions ...Condition) bool {
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
