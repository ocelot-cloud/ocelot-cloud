package component_tests

import (
	"github.com/ocelot-cloud/shared"
	"github.com/ocelot-cloud/shared/assert"
	"github.com/ocelot-cloud/shared/utils"
	"ocelot/backend/tools"
	"os"
	"testing"
)

var logger = tools.Logger

const (
	backendBaseUrl = "http://localhost:8080"
	TestProfile    = "TEST"
	ProdProfile    = "PROD"
)

func TestHappyPathDeployAndStop(t *testing.T) {
	cloud := getClientAndLogin(t)
	_, err := cloud.readApps()
	assert.Nil(t, err)
	assert.Nil(t, cloud.assertState("Uninitialized"))

	assert.Nil(t, cloud.startApp())
	_, err = cloud.readApps()
	assert.Nil(t, err)
	assert.Nil(t, cloud.assertState("Available"))

	assert.Nil(t, cloud.stopApp())
	_, err = cloud.readApps()
	assert.Nil(t, err)
	assert.Nil(t, cloud.assertState("Uninitialized"))
}

// TODO I think it should be a 400 error, since its the users fault
func TestDeployStackNotExisting(t *testing.T) {
	cloud := getClientAndLogin(t)
	cloud.appToOperateOn = "not-existing-stack"
	err := cloud.startApp()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(500, "Deploying stack failed: not-existing-stack"), err.Error())
}

func TestStopStackNotExisting(t *testing.T) {
	cloud := getClientAndLogin(t)
	cloud.appToOperateOn = "not-existing-stack"
	err := cloud.stopApp()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(500, "Stopping stack failed: not-existing-stack"), err.Error())
}

func TestAbsenceOfCorsPolicyDisablingHeadersInResponse(t *testing.T) {
	onlyExecuteTestForProfile(t, ProdProfile)
	AssertCorsHeaders(t, "", "", "", "")
}

func AssertCorsHeaders(t *testing.T, expectedAllowOrigin, expectedAllowMethods, expectedAllowHeaders, expectedAllowCredentials string) {
	cloud := getClientAndLogin(t)
	resp, err := cloud.parent.DoRequestWithFullResponse("/api/stacks/read", nil, "")
	assert.Nil(t, err)

	allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
	assert.Equal(t, expectedAllowOrigin, allowOrigin)

	allowMethods := resp.Header.Get("Access-Control-Allow-Methods")
	assert.Equal(t, expectedAllowMethods, allowMethods)

	allowHeaders := resp.Header.Get("Access-Control-Allow-Headers")
	assert.Equal(t, expectedAllowHeaders, allowHeaders)

	allowCredentials := resp.Header.Get("Access-Control-Allow-Credentials")
	assert.Equal(t, expectedAllowCredentials, allowCredentials)
}

func TestUrlPaths(t *testing.T) {
	cloud := getClientAndLogin(t)
	apps, err := cloud.readApps()
	assert.Nil(t, err)
	isCustomPathNginxPathOk := false
	isDefaultNginxPathOk := false
	for _, app := range *apps {
		if app.Name == tools.NginxCustomPath && app.UrlPath == "/custom-path" {
			isCustomPathNginxPathOk = true
		} else if app.Name == tools.NginxDefault && app.UrlPath == "/" {
			isDefaultNginxPathOk = true
		}
	}
	assert.True(t, isCustomPathNginxPathOk)
	assert.True(t, isDefaultNginxPathOk)
}

func TestNetworkCreationOnStackDeployment(t *testing.T) {
	onlyExecuteTestForProfile(t, ProdProfile)
	cloud := getClientAndLogin(t)

	_ = shared.ExecuteShellCommand("docker network ls | grep -q nginx-default-net || docker network rm nginx-default-net")
	assert.Nil(t, cloud.startApp())
	err := shared.ExecuteShellCommand("docker network ls | grep -q nginx-default-net")
	assert.Nil(t, err)
}

func TestWhetherCorsPolicyDisablingHeadersAreInResponse(t *testing.T) {
	onlyExecuteTestForProfile(t, TestProfile)
	AssertCorsHeaders(t, backendBaseUrl, "POST, GET, OPTIONS, PUT, DELETE", "Accept, Content-Type, Content-Length, Authorization", "true")
}

func TestHealthStateOfSlowStartingStack(t *testing.T) {
	onlyExecuteTestForProfile(t, ProdProfile)
	cloud := getClientAndLogin(t)
	cloud.appToOperateOn = tools.NginxSlowStart
	assert.Nil(t, cloud.startApp())
	assert.Nil(t, cloud.assertState("Starting"))
	assert.Nil(t, cloud.assertState("Available"))
}

func onlyExecuteTestForProfile(t *testing.T, profileEnablingTheTest string) {
	setEnvProfile, _ := os.LookupEnv("PROFILE")
	if setEnvProfile != profileEnablingTheTest {
		t.Skip()
	}
}
