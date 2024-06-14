package internal

import (
	"github.com/ocelot-cloud/shared"
	"ocelot/tools"
	"testing"
)

func init() {
	StackFileDir = "../../../stacks/dummy"
}

func TestWhetherExistingUrlPathIsCorrectlyRead(t *testing.T) {
	yamlConfig := ProvideStackConfigService(StackFileDir)
	limesurveyUrlPath := yamlConfig.GetStackConfig(tools.NginxCustomPath).UrlPath
	shared.AssertEqual(t, "/custom-path", limesurveyUrlPath)
}

func TestMissingYamlFileLeadsToReturnOfIndexPath(t *testing.T) {
	assertEmptyUrlPathForStack(t, tools.NginxDefault2)
}

func assertEmptyUrlPathForStack(t *testing.T, stackName string) {
	yamlConfig := ProvideStackConfigService(StackFileDir)
	missingYamlFileUrlPathDefaultValue := yamlConfig.GetStackConfig(stackName).UrlPath
	shared.AssertEqual(t, "/", missingYamlFileUrlPathDefaultValue)
}

func TestMissingUrlPathVariableInYamlFileLeadsToReturnOfIndexPath(t *testing.T) {
	assertEmptyUrlPathForStack(t, tools.NginxDefault)
}

func TestNonExistentStackShouldReturnDefaultConfig(t *testing.T) {
	yamlConfig := ProvideStackConfigService(StackFileDir)
	resultConfig := yamlConfig.GetStackConfig("non-existent-stack")
	shared.AssertEqual(t, "/", resultConfig.UrlPath)
	shared.AssertEqual(t, "80", resultConfig.Port)
}

func TestStackConfig(t *testing.T) {
	type testCase struct {
		Name         string
		StackName    string
		ExpectedPort string
		ExpectedPath string
	}
	testCases := []testCase{
		{"nginx-default2", "nginx-default2", "80", "/"},
		{"nginx-custom-path", "nginx-custom-path", "80", "/custom-path"},
		{"nginx-custom-port", "nginx-custom-port", "3000", "/"},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			stackConfigService := ProvideStackConfigService("../../../stacks/dummy")
			config := stackConfigService.GetStackConfig(tc.StackName)
			shared.AssertEqual(t, tc.ExpectedPort, config.Port)
			shared.AssertEqual(t, tc.ExpectedPath, config.UrlPath)
		})
	}
}
