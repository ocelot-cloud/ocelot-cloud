package yaml_config

import (
	"github.com/ocelot-cloud/shared/assert"
	"ocelot/backend/apps/global_config"
	"ocelot/backend/tools"
	"testing"
)

func init() {
	global_config.AppFileDir = "../" + global_config.DummyAppAssetsDirForTests
}

func TestWhetherExistingUrlPathIsCorrectlyRead(t *testing.T) {
	yamlConfig := ProvideAppConfigService()
	path := yamlConfig.GetAppConfig(tools.NginxCustomPath).UrlPath
	assert.Equal(t, "/custom-path", path)
}

func TestMissingYamlFileLeadsToReturnOfIndexPath(t *testing.T) {
	assertEmptyUrlPathForStack(t, tools.NginxDefault2)
}

func assertEmptyUrlPathForStack(t *testing.T, stackName string) {
	yamlConfig := ProvideAppConfigService()
	missingYamlFileUrlPathDefaultValue := yamlConfig.GetAppConfig(stackName).UrlPath
	assert.Equal(t, "/", missingYamlFileUrlPathDefaultValue)
}

func TestMissingUrlPathVariableInYamlFileLeadsToReturnOfIndexPath(t *testing.T) {
	assertEmptyUrlPathForStack(t, tools.NginxDefault)
}

func TestNonExistentStackShouldReturnDefaultConfig(t *testing.T) {
	yamlConfig := ProvideAppConfigService()
	resultConfig := yamlConfig.GetAppConfig("non-existent-stack")
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
			stackConfigService := ProvideAppConfigService()
			config := stackConfigService.GetAppConfig(tc.StackName)
			assert.Equal(t, tc.ExpectedPort, config.Port)
			assert.Equal(t, tc.ExpectedPath, config.UrlPath)
		})
	}
}
