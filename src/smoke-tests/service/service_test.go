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
	"smoke-tests/service/reporter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/pborman/uuid"
	"github.com/pivotal-cf-experimental/cf-test-helpers/services"
)

type CFTestContext struct {
	Org, Space string
}

var _ = Describe("MongoDB Service", func() {
	var (
		// client = adapter.NewPCGCClient(mongodbConfig.URL, mongodbConfig.UserName, mongodbConfig.UserAPIKey)
		testCF = smokeTestCF.CF{
			ShortTimeout: time.Minute * 3,
			LongTimeout:  time.Minute * 15,
			RetryBackoff: mongodbConfig.Retry.Backoff(),
			MaxRetries:   mongodbConfig.Retry.MaxRetries(),
		}

		retryInterval = time.Second

		appPath             = "../assets/cf-mongo-example-app"
		serviceInstanceName string
		appName             string
		planName            string
		securityGroupName   string
		serviceKeyName      string

		cfTestContext CFTestContext

		context services.Context
	)

	SynchronizedBeforeSuite(func() []byte {
		context = services.NewContext(cfTestConfig, "mongodb-test")

		createQuotaArgs := []string{
			"-m", "10G",
			"-r", "1000",
			"-s", "100",
			"--allow-paid-service-plans",
		}

		regularContext := context.RegularUserContext()

		beforeSuiteSteps := []*reporter.Step{
			reporter.NewStep(
				"Connect to CloudFoundry",
				testCF.API(cfTestConfig.ApiEndpoint, cfTestConfig.SkipSSLValidation),
			),
			reporter.NewStep(
				"Log in as admin",
				testCF.Auth(cfTestConfig.AdminUser, cfTestConfig.AdminPassword),
			),
			reporter.NewStep(
				"Create 'mongodb-smoke-tests' quota",
				testCF.CreateQuota("mongodb-smoke-test-quota", createQuotaArgs...),
			),
			reporter.NewStep(
				fmt.Sprintf("Create '%s' org", regularContext.Org),
				testCF.CreateOrg(regularContext.Org, "mongodb-smoke-test-quota"),
			),
			reporter.NewStep(
				fmt.Sprintf("Enable service access for '%s' org", regularContext.Org),
				testCF.EnableServiceAccess(regularContext.Org, mongodbConfig.ServiceName),
			),
			reporter.NewStep(
				fmt.Sprintf("Target '%s' org", regularContext.Org),
				testCF.TargetOrg(regularContext.Org),
			),
			reporter.NewStep(
				fmt.Sprintf("Create '%s' space", regularContext.Space),
				testCF.CreateSpace(regularContext.Space),
			),
			reporter.NewStep(
				fmt.Sprintf("Target '%s' org and '%s' space", regularContext.Org, regularContext.Space),
				testCF.TargetOrgAndSpace(regularContext.Org, regularContext.Space),
			),
			reporter.NewStep(
				"Log out",
				testCF.Logout(),
			),
		}

		smokeTestReporter.RegisterBeforeSuiteSteps(beforeSuiteSteps)

		for _, task := range beforeSuiteSteps {
			task.Perform()
		}

		cfTestContext = CFTestContext{
			Org:   regularContext.Org,
			Space: regularContext.Space,
		}

		rawTestContext, err := json.Marshal(cfTestContext)
		Expect(err).NotTo(HaveOccurred())

		return rawTestContext
	}, func(data []byte) {
		err := json.Unmarshal(data, &cfTestContext)
		Expect(err).NotTo(HaveOccurred())

		// Set $CF_HOME so that cf cli state is not shared between nodes
		cfHomeDir, err := ioutil.TempDir("", "cf-mongodb-smoke-tests")
		Expect(err).NotTo(HaveOccurred())
		os.Setenv("CF_HOME", cfHomeDir)
	})

	BeforeEach(func() {
		appName = randomName()
		serviceInstanceName = randomName()
		securityGroupName = randomName()
		serviceKeyName = randomName()

		pushArgs := []string{
			"-m", "256M",
			"-p", appPath,
			"-s", "cflinuxfs3",
			"--no-start",
		}

		specSteps := []*reporter.Step{
			reporter.NewStep(
				"Connect to CloudFoundry",
				testCF.API(cfTestConfig.ApiEndpoint, cfTestConfig.SkipSSLValidation),
			),
			reporter.NewStep(
				"Log in as admin",
				testCF.Auth(cfTestConfig.AdminUser, cfTestConfig.AdminPassword),
			),
			reporter.NewStep(
				fmt.Sprintf("Target '%s' org and '%s' space", cfTestContext.Org, cfTestContext.Space),
				testCF.TargetOrgAndSpace(cfTestContext.Org, cfTestContext.Space),
			),
			reporter.NewStep(
				"Push the MongoDB sample app to Cloud Foundry",
				testCF.Push(appName, pushArgs...),
			),
		}

		smokeTestReporter.ClearSpecSteps()
		smokeTestReporter.RegisterSpecSteps(specSteps)

		for _, task := range specSteps {
			task.Perform()
		}
	})

	AfterEach(func() {
		specSteps := []*reporter.Step{
			reporter.NewStep(
				fmt.Sprintf("Unbind the %q plan instance", planName),
				testCF.UnbindService(appName, serviceInstanceName),
			),
			reporter.NewStep(
				fmt.Sprintf("Delete the service key %s for the %q plan instance", serviceKeyName, planName),
				testCF.DeleteServiceKey(serviceInstanceName, serviceKeyName),
			),
			reporter.NewStep(
				fmt.Sprintf("Delete the %q plan instance", planName),
				testCF.DeleteService(serviceInstanceName),
			),
			reporter.NewStep(
				fmt.Sprintf("Ensure service instance for plan %q has been deleted", planName),
				testCF.EnsureServiceInstanceGone(serviceInstanceName),
			),
			reporter.NewStep(
				"Delete the app",
				testCF.Delete(appName),
			),
			reporter.NewStep(
				fmt.Sprintf("Delete security group '%s'", securityGroupName),
				testCF.DeleteSecurityGroup("mongodb-smoke-tests-sg"),
			),
		}

		smokeTestReporter.RegisterSpecSteps(specSteps)

		for _, task := range specSteps {
			task.Perform()
		}
	})

	SynchronizedAfterSuite(func() {}, func() {
		afterSuiteSteps := []*reporter.Step{
			reporter.NewStep(
				"Ensure no service-instances left",
				testCF.EnsureAllServiceInstancesGone(),
			),
			reporter.NewStep(
				fmt.Sprintf("Delete org '%s'", cfTestContext.Org),
				testCF.DeleteOrg(cfTestContext.Org),
			),
			reporter.NewStep(
				"Log out",
				testCF.Logout(),
			),
		}

		smokeTestReporter.RegisterAfterSuiteSteps(afterSuiteSteps)

		afterSuiteSteps[0].Perform()
	})

	AssertLifeCycleBehavior := func(sp ServiceParameters) {
		It(sp.PrintParameters()+": create, bind to, write to, read from, unbind, and destroy a service instance", func() {
			var skip bool

			uri := fmt.Sprintf("https://%s.%s", appName, cfTestConfig.AppsDomain)
			app := mongodb.NewApp(uri, testCF.ShortTimeout, retryInterval)
			testValue := randomName()
			serviceConfig := fmt.Sprintf(`{"enable_backup":"%s", "version":"%s"}`, sp.BackupEnable, sp.MongoDBVersion)
			fmt.Println("serviceName : ", sp.ServiceName, " planName: ", sp.PlanName, " serviceInstanceName: ", serviceInstanceName,
				"configuration : -c ", serviceConfig)

			serviceCreateStep := reporter.NewStep(
				fmt.Sprintf("Create a '%s' plan instance of MongoDB", sp.PlanName),
				testCF.CreateService(sp.ServiceName, sp.PlanName, serviceInstanceName, serviceConfig, &skip),
			)

			smokeTestReporter.RegisterSpecSteps([]*reporter.Step{serviceCreateStep})

			specSteps := []*reporter.Step{
				reporter.NewStep(
					fmt.Sprintf("Bind the mongodb sample app '%s' to the '%s' plan instance '%s' of MongoDB", appName, planName, serviceInstanceName),
					testCF.BindService(appName, serviceInstanceName),
				),
				reporter.NewStep(
					fmt.Sprintf("Create service key for the '%s' plan instance '%s' of MongoDB", planName, serviceInstanceName),
					testCF.CreateServiceKey(serviceInstanceName, serviceKeyName),
				),
				reporter.NewStep(
					fmt.Sprintf("Create and bind security group '%s' for running smoke tests", securityGroupName),
					testCF.CreateAndBindSecurityGroup(securityGroupName, cfTestContext.Org, cfTestContext.Space),
				),
				reporter.NewStep(
					"Start the app",
					testCF.Start(appName),
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
						groupID := testCF.GetGroupID(serviceInstanceName)
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
			}

			smokeTestReporter.RegisterSpecSteps(specSteps)

			serviceCreateStep.Perform()
			serviceCreateStep.Description = fmt.Sprintf("Create a '%s' plan instance of MongoDB", planName)

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
		cases := generateTestServiceParameters(mongodbConfig)
		printGeneratedServiceParameters(cases)
		for _, oneCase := range cases {
			planName = oneCase.PlanName
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

func randomName() string {
	return uuid.NewRandom().String()
}
