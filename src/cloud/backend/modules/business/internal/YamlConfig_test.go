package internal

import (
	"ocelot/tools"
	"testing"
)

func init() {
	StackFileDir = "../../../stacks/dummy"
}

func TestWhetherExistingUrlPathIsCorrectlyRead(t *testing.T) {
	yamlConfig := ProvideStackConfigService(StackFileDir)
	limesurveyUrlPath := yamlConfig.GetStackConfig(tools.NginxCustomPath).UrlPath
	tools.AssertEqual(t, "/custom-path", limesurveyUrlPath)
}

func TestMissingYamlFileLeadsToReturnOfIndexPath(t *testing.T) {
	assertEmptyUrlPathForStack(t, tools.NginxDefault2)
}

func assertEmptyUrlPathForStack(t *testing.T, stackName string) {
	yamlConfig := ProvideStackConfigService(StackFileDir)
	missingYamlFileUrlPathDefaultValue := yamlConfig.GetStackConfig(stackName).UrlPath
	tools.AssertEqual(t, "/", missingYamlFileUrlPathDefaultValue)
}

func TestMissingUrlPathVariableInYamlFileLeadsToReturnOfIndexPath(t *testing.T) {
	assertEmptyUrlPathForStack(t, tools.NginxDefault)
}

func TestNonExistentStackShouldReturnDefaultConfig(t *testing.T) {
	yamlConfig := ProvideStackConfigService(StackFileDir)
	resultConfig := yamlConfig.GetStackConfig("non-existent-stack")
	tools.AssertEqual(t, "/", resultConfig.UrlPath)
	tools.AssertEqual(t, "80", resultConfig.Port)
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
			tools.AssertEqual(t, tc.ExpectedPort, config.Port)
			tools.AssertEqual(t, tc.ExpectedPath, config.UrlPath)
		})
	}
}
