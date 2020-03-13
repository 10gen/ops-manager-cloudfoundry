//steps for 
package steps

import (
	"fmt"

	"smoke-tests/cf"
	c "smoke-tests/service/configuration"
	"smoke-tests/service/data"
	"smoke-tests/service/reporter"

	"github.com/pivotal-cf-experimental/cf-test-helpers/services"
)

//InstallServiceSteps is set of steps for installing service with config, pass nil if you don't need to use -c
func InstallServiceSteps() {

}

//DeleteServiceSteps is set of steps for deleting particular service instance
func DeleteServiceSteps(cf cf.CF, i c.InstanceNames, planName string) []*reporter.Step{
	specSteps := []*reporter.Step{
		reporter.NewStep(
			fmt.Sprintf("Unbind the %q plan instance", planName),
			cf.UnbindService(i.AppName, i.ServiceTestName),
		),
		reporter.NewStep(
			fmt.Sprintf("Delete the service key %s for the %q plan instance", i.ServiceKeyName, planName),
			cf.DeleteServiceKey(i.ServiceTestName, i.ServiceKeyName),
		),
		reporter.NewStep(
			fmt.Sprintf("Delete the %q plan instance", planName),
			cf.DeleteService(i.ServiceTestName),
		),
		reporter.NewStep(
			fmt.Sprintf("Ensure service instance for plan %q has been deleted", planName),
			cf.EnsureServiceInstanceGone(i.ServiceTestName),
		),
		reporter.NewStep(
			"Delete the app",
			cf.Delete(i.AppName),
		),
		reporter.NewStep(
			fmt.Sprintf("Delete security group '%s'", i.SecurityGroupName),
			cf.DeleteSecurityGroup("mongodb-smoke-tests-sg"),
		),
	}
	for _, task := range specSteps {
		task.Perform()
	}
	return specSteps
}

//PushApplication is set of steps for pushing our test application data.PushArgs
func PushApplication(cf cf.CF, i c.InstanceNames, cfTestConfig services.Config) []*reporter.Step{ //TODO refactor
	specSteps := []*reporter.Step{
		reporter.NewStep(
			"Connect to CloudFoundry",
			cf.API(cfTestConfig.ApiEndpoint, cfTestConfig.SkipSSLValidation),
		),
		reporter.NewStep(
			"Log in as admin",
			cf.Auth(cfTestConfig.AdminUser, cfTestConfig.AdminPassword),
		),
		reporter.NewStep(
			fmt.Sprintf("Target '%s' org and '%s' space", i.Context.Org, i.Context.Space),
			cf.TargetOrgAndSpace(i.Context.Org, i.Context.Space),
		),
		reporter.NewStep(
			"Push the MongoDB sample app to Cloud Foundry",
			cf.Push(i.AppName, data.PushArgs...),
		),
	}

	// smokeTestReporter.ClearSpecSteps()
	// smokeTestReporter.RegisterSpecSteps(specSteps)
	for _, task := range specSteps {
		task.Perform()
	}
	return specSteps
}

//CreateOrganization - create organization steps.
func CreateOrganization(cf cf.CF, i c.InstanceNames, cfTestConfig services.Config, serviceName string) []*reporter.Step {
	beforeSuiteSteps := []*reporter.Step{
		reporter.NewStep(
			"Connect to CloudFoundry",
			cf.API(cfTestConfig.ApiEndpoint, cfTestConfig.SkipSSLValidation),
		),
		reporter.NewStep(
			"Log in as admin",
			cf.Auth(cfTestConfig.AdminUser, cfTestConfig.AdminPassword),
		),
		reporter.NewStep(
			"Create 'mongodb-smoke-tests' quota",
			cf.CreateQuota("mongodb-smoke-test-quota", data.CreateQuotaArgs...),
		),
		reporter.NewStep(
			fmt.Sprintf("Create '%s' org", i.Context.Org),
			cf.CreateOrg(i.Context.Org, "mongodb-smoke-test-quota"),
		),
		reporter.NewStep(
			fmt.Sprintf("Enable service access for '%s' org", i.Context.Org),
			cf.EnableServiceAccess(i.Context.Org, serviceName), //TODO refactor 
		),
		reporter.NewStep(
			fmt.Sprintf("Target '%s' org", i.Context.Org),
			cf.TargetOrg(i.Context.Org),
		),
		reporter.NewStep(
			fmt.Sprintf("Create '%s' space", i.Context.Space),
			cf.CreateSpace(i.Context.Space),
		),
		reporter.NewStep(
			fmt.Sprintf("Target '%s' org and '%s' space", i.Context.Org, i.Context.Space),
			cf.TargetOrgAndSpace(i.Context.Org, i.Context.Space),
		),
		reporter.NewStep(
			"Log out",
			cf.Logout(),
		),
	}
	for _, task := range beforeSuiteSteps {
		task.Perform()
	}
	return beforeSuiteSteps
}

//DeleteOrganization - delete organization steps
func DeleteOrganization(cf cf.CF, i c.InstanceNames) []*reporter.Step {
	afterSuiteSteps := []*reporter.Step{
		reporter.NewStep(
			"Ensure no service-instances left",
			cf.EnsureAllServiceInstancesGone(),
		),
		reporter.NewStep(
			fmt.Sprintf("Delete org '%s'", i.Context.Org),
			cf.DeleteOrg(i.Context.Org),
		),
		reporter.NewStep(
			"Log out",
			cf.Logout(),
		),
	}
	for _, step := range afterSuiteSteps {
		step.Perform()
	}
	return afterSuiteSteps
}
