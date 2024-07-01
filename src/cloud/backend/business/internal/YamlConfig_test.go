package internal

import (
	"github.com/ocelot-cloud/shared/assert"
	"ocelot/backend/config"
	"testing"
)

var DefaultStackFileDir = "../../stacks/dummy"

func init() {
	StackFileDir = DefaultStackFileDir
}

func TestWhetherExistingUrlPathIsCorrectlyRead(t *testing.T) {
	yamlConfig := ProvideStackConfigService(StackFileDir)
	limesurveyUrlPath := yamlConfig.GetStackConfig(tools.NginxCustomPath).UrlPath
	assert.Equal(t, "/custom-path", limesurveyUrlPath)
}

func TestMissingYamlFileLeadsToReturnOfIndexPath(t *testing.T) {
	assertEmptyUrlPathForStack(t, tools.NginxDefault2)
}

func assertEmptyUrlPathForStack(t *testing.T, stackName string) {
	yamlConfig := ProvideStackConfigService(StackFileDir)
	missingYamlFileUrlPathDefaultValue := yamlConfig.GetStackConfig(stackName).UrlPath
	assert.Equal(t, "/", missingYamlFileUrlPathDefaultValue)
}

func TestMissingUrlPathVariableInYamlFileLeadsToReturnOfIndexPath(t *testing.T) {
	assertEmptyUrlPathForStack(t, tools.NginxDefault)
}

func TestNonExistentStackShouldReturnDefaultConfig(t *testing.T) {
	yamlConfig := ProvideStackConfigService(StackFileDir)
	resultConfig := yamlConfig.GetStackConfig("non-existent-stack")
	assert.Equal(t, "/", resultConfig.UrlPath)
	assert.Equal(t, "80", resultConfig.Port)
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
			stackConfigService := ProvideStackConfigService(DefaultStackFileDir)
			config := stackConfigService.GetStackConfig(tc.StackName)
			assert.Equal(t, tc.ExpectedPort, config.Port)
			assert.Equal(t, tc.ExpectedPath, config.UrlPath)
		})
	}
}
