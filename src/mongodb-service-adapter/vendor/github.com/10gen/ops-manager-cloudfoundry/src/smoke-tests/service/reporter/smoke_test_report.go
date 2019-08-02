package reporter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/types"
)

type Step struct {
	Description string
	Result      string
	Task        func()
}

func (step *Step) Perform() {
	step.Result = "FAILED"
	step.Task()
	step.Result = "PASSED"
}

func NewStep(description string, task func()) *Step {
	return &Step{
		Description: description,
		Result:      "DIDN'T RUN",
		Task:        task,
	}
}

type failure struct {
	title   string
	message string
}

type SmokeTestReport struct {
	testCount        int
	failures         []failure
	beforeSuitesteps []*Step
	afterSuiteSteps  []*Step
	specSteps        []*Step
}

func (report *SmokeTestReport) RegisterBeforeSuiteSteps(steps []*Step) {
	report.beforeSuitesteps = append(report.beforeSuitesteps, steps...)
}

func (report *SmokeTestReport) RegisterAfterSuiteSteps(steps []*Step) {
	report.afterSuiteSteps = append(report.afterSuiteSteps, steps...)
}

func (report *SmokeTestReport) RegisterSpecSteps(steps []*Step) {
	report.specSteps = append(report.specSteps, steps...)
}

func (report *SmokeTestReport) ClearSpecSteps() {
	report.specSteps = []*Step{}
}

func (report *SmokeTestReport) SpecSuiteWillBegin(
	config config.GinkgoConfigType,
	summary *types.SuiteSummary,
) {
	if ginkgo.GinkgoParallelNode() != 1 {
		return
	}
	report.printMessageTitle("Beginning test suite setup")
}

func (report *SmokeTestReport) BeforeSuiteDidRun(summary *types.SetupSummary) {
	if ginkgo.GinkgoParallelNode() != 1 {
		return
	}
	if summary.State == types.SpecStateFailed ||
		summary.State == types.SpecStatePanicked ||
		summary.State == types.SpecStateTimedOut {

		report.failures = append(report.failures, failure{
			title:   "Suite setup",
			message: summary.Failure.Message,
		})
	}
	report.printMessageTitle("Finished test suite setup")

	fmt.Println("Smoke Test Suite Setup Results:")
	count := len(report.beforeSuitesteps)
	for i, step := range report.beforeSuitesteps {
		fmt.Printf("[%d/%d] %s: %s\n", i+1, count, step.Description, step.Result)
	}
	fmt.Println()
}

func (report *SmokeTestReport) SpecWillRun(summary *types.SpecSummary) {
	report.testCount++

	title := report.getTitleFromComponents(summary)
	message := fmt.Sprintf("START %d. %s", report.testCount, title)
	report.printMessageTitle(message)
}

func (report *SmokeTestReport) SpecDidComplete(summary *types.SpecSummary) {
	if summary.Failed() {
		report.failures = append(report.failures, failure{
			title:   summary.ComponentTexts[len(summary.ComponentTexts)-1],
			message: summary.Failure.Message,
		})
	}
	title := report.getTitleFromComponents(summary)
	message := fmt.Sprintf("END %d. %s", report.testCount, title)
	report.printMessageTitle(message)

	fmt.Println("Smoke Test plan Results:")
	count := len(report.specSteps)
	for i, step := range report.specSteps {
		fmt.Printf("[%d/%d] %s: %s\n", i+1, count, step.Description, step.Result)
	}
	fmt.Println()
}

func (report *SmokeTestReport) AfterSuiteDidRun(summary *types.SetupSummary) {
	if ginkgo.GinkgoParallelNode() != 1 {
		return
	}
	report.printMessageTitle("Finished suite teardown")

	fmt.Println("Smoke Test Suite Teardown Results:")
	count := len(report.afterSuiteSteps)
	for i, step := range report.afterSuiteSteps {
		fmt.Printf("[%d/%d] %s: %s\n", i+1, count, step.Description, step.Result)
	}
	fmt.Println()
}

func (report *SmokeTestReport) SpecSuiteDidEnd(summary *types.SuiteSummary) {
	if ginkgo.GinkgoParallelNode() != 1 {
		return
	}
	matchJSON, err := regexp.Compile(`{"FailReason":\s"(.*)"}`)
	if err != nil {
		fmt.Printf("\nSkipping \"Summarising failure reasons\": %s\n", err.Error())
		return
	}

	if summary.NumberOfFailedSpecs > 0 {
		report.printMessageTitle("Summarising Failures")

		for _, failure := range report.failures {
			fmt.Printf("\n%s\n", failure.title)

			failMessage := matchJSON.FindStringSubmatch(failure.message)
			if failMessage != nil {
				fmt.Printf("> %s\n", failMessage[1])
			}
		}
		fmt.Printf("\nFor help with troubleshooting, visit: https://docs.pivotal.io/redis/smoke-tests.html\n\n")
	}
}

func (report *SmokeTestReport) getTitleFromComponents(summary *types.SpecSummary) (title string) {
	if len(summary.ComponentTexts) > 0 {
		title = summary.ComponentTexts[len(summary.ComponentTexts)-1]
	}
	return
}

func (report *SmokeTestReport) printMessageTitle(message string) {
	border := strings.Repeat("-", len(message)+2)
	fmt.Printf("\n\n|%s|\n", border)
	fmt.Printf("| %s |\n", message)
	fmt.Printf("|%s|\n\n", border)
}
