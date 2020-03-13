package service_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"mongodb-service-adapter/adapter"

	smokeTestCF "smoke-tests/cf"
	"smoke-tests/mongodb"
	prepare "smoke-tests/service/configuration"
	"smoke-tests/service/reporter"
	"smoke-tests/service/steps"


	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/pborman/uuid"
	"github.com/pivotal-cf-experimental/cf-test-helpers/services"
)

var _ = Describe("MongoDB Service", func() {
	var (
		testCF = smokeTestCF.CF{
			ShortTimeout: time.Minute * 3,
			LongTimeout:  time.Minute * 15,
			RetryBackoff: mongodbConfig.Retry.Backoff(),
			MaxRetries:   mongodbConfig.Retry.MaxRetries(),
		}
		retryInterval = time.Second

		planName string

		caseInstance prepare.InstanceNames
	)

	SynchronizedBeforeSuite(func() []byte {
		caseInstance.GenerateInstanceNames()
		regularContext := services.NewContext(cfTestConfig, "mongodb-test").RegularUserContext()
		caseInstance.Context.Org = regularContext.Org
		caseInstance.Context.Space = regularContext.Space

		smokeTestReporter.RegisterBeforeSuiteSteps(steps.CreateOrganization(testCF, caseInstance, cfTestConfig, mongodbConfig.ServiceName))

		rawTestContext, err := json.Marshal(caseInstance.Context)
		Expect(err).NotTo(HaveOccurred())

		return rawTestContext
	}, func(data []byte) {
		err := json.Unmarshal(data, &caseInstance.Context)
		Expect(err).NotTo(HaveOccurred())

		// Set $CF_HOME so that cf cli state is not shared between nodes
		cfHomeDir, err := ioutil.TempDir("", "cf-mongodb-smoke-tests")
		Expect(err).NotTo(HaveOccurred())
		os.Setenv("CF_HOME", cfHomeDir)
	})

	BeforeEach(func() {
		smokeTestReporter.ClearSpecSteps()
		smokeTestReporter.RegisterSpecSteps(steps.PushApplication(testCF, caseInstance, cfTestConfig))
	})

	AfterEach(func() {
		smokeTestReporter.RegisterSpecSteps(steps.DeleteServiceSteps(testCF, caseInstance, planName))
	})

	SynchronizedAfterSuite(func() {}, func() {
		smokeTestReporter.RegisterAfterSuiteSteps(steps.DeleteOrganization(testCF, caseInstance))
	})

	AssertLifeCycleBehavior := func(sp prepare.ServiceParameters) {
		It(sp.PrintParameters()+": create, bind to, write to, read from, unbind, and destroy a service instance", func() {
			var skip bool
			uri := fmt.Sprintf("https://%s.%s", caseInstance.AppName, cfTestConfig.AppsDomain)
			app := mongodb.NewApp(uri, testCF.ShortTimeout, retryInterval)
			testValue := uuid.NewRandom().String()
			serviceConfig := fmt.Sprintf(`{"enable_backup":"%s", "version":"%s"}`, sp.BackupEnable, sp.MongoDBVersion)
			fmt.Println("serviceName : ", sp.MarketPlaceServiceName, " planName: ", sp.PlanName, " serviceInstanceName: ", caseInstance.ServiceTestName,
				"configuration : -c ", serviceConfig)

			serviceCreateStep := reporter.NewStep(
				fmt.Sprintf("Create a '%s' plan instance of MongoDB", sp.PlanName),
				testCF.CreateService(sp.MarketPlaceServiceName, sp.PlanName, caseInstance.ServiceTestName, serviceConfig, &skip),
			)

			smokeTestReporter.RegisterSpecSteps([]*reporter.Step{serviceCreateStep})

			specSteps := []*reporter.Step{
				reporter.NewStep(
					fmt.Sprintf("Bind the mongodb sample app '%s' to the '%s' plan instance '%s' of MongoDB", caseInstance.AppName, sp.PlanName, caseInstance.ServiceTestName),
					testCF.BindService(caseInstance.AppName, caseInstance.ServiceTestName),
				),
				reporter.NewStep(
					fmt.Sprintf("Create service key for the '%s' plan instance '%s' of MongoDB", sp.PlanName, caseInstance.ServiceTestName),
					testCF.CreateServiceKey(caseInstance.ServiceTestName, caseInstance.ServiceKeyName),
				),
				reporter.NewStep(
					fmt.Sprintf("Create and bind security group '%s' for running smoke tests", caseInstance.SecurityGroupName),
					testCF.CreateAndBindSecurityGroup(caseInstance.SecurityGroupName, caseInstance.Context.Org, caseInstance.Context.Space),
				),
				reporter.NewStep(
					"Start the app",
					testCF.Start(caseInstance.AppName),
				),
				reporter.NewStep(
					"Verify that the app is responding",
					app.IsRunning(),
				),
				reporter.NewStep(
					"Write to MongoDB",
					app.Write("testkey", testValue),
				),
				reporter.NewStep(
					"Read from MongoDB",
					app.ReadAssert("testkey", testValue),
				),
				reporter.NewStep(
					"Backup configuration",
					func() {
						client := adapter.NewPCGCClient(mongodbConfig.OpsMan.URL, mongodbConfig.OpsMan.UserName, mongodbConfig.OpsMan.UserAPIKey)
						groupID := testCF.GetGroupID(caseInstance.ServiceTestName)
						fmt.Printf("Got groupID/ProjectID: %s", groupID)
						backupConfigs, err := client.GetBackupConfigs(groupID)
						Expect(err).NotTo(HaveOccurred())
						for _, config := range backupConfigs.BackupResults {
							Expect(config).To(MatchFields(IgnoreExtras, Fields{
								"StatusName": Equal(convertedBackupStatus(sp.BackupEnable)),
							}))
						}
					},
				),
				reporter.NewStep(
					"Check requared version",
					func() {
						client := adapter.NewPCGCClient(mongodbConfig.OpsMan.URL, mongodbConfig.OpsMan.UserName, mongodbConfig.OpsMan.UserAPIKey)
						groupID := testCF.GetGroupID(caseInstance.ServiceTestName)
						autoConfig, err := client.GetAutomationConfig(groupID)
						Expect(err).NotTo(HaveOccurred())
						for _, process := range autoConfig.Processes {
							Expect(process).To(PointTo(MatchFields(IgnoreExtras, Fields{
								"Version": Equal(sp.MongoDBVersion),
							})))
						}
					},
				),
			}

			smokeTestReporter.RegisterSpecSteps(specSteps)

			serviceCreateStep.Perform()
			serviceCreateStep.Description = fmt.Sprintf("Create a '%s' plan instance of MongoDB", sp.PlanName)

			if skip {
				serviceCreateStep.Result = "SKIPPED"
			} else {
				for _, task := range specSteps {
					task.Perform()
				}
			}
		})
	}

	Context("for each plan", func() {
		mongodbConfig.SetDefaultForNonDefinedParameters()
		cases := prepare.GenerateTestServiceParameters(mongodbConfig)
		prepare.PrintGeneratedServiceParameters(cases)
		for _, oneCase := range cases {
			planName = oneCase.PlanName
			// oneCase.Instance = caseInstance
			AssertLifeCycleBehavior(oneCase)
		}
	})
})

//TODO move somewhere
func convertedBackupStatus(backupEnable string) string {
	if backupEnable == "true" {
		return "STARTED"
	}
	return "INACTIVE"
}
