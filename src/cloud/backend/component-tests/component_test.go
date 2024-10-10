package component_tests

import (
	"github.com/ocelot-cloud/shared"
	"github.com/ocelot-cloud/shared/assert"
	"github.com/ocelot-cloud/shared/utils"
	"ocelot/backend/tools"
	"os"
	"testing"
	"time"
)

// TODO Add a "before each" function to initialize a global "var cloud" for these tests?
// TODO Add a "wipe" endpoint that stops all stacks and it also deletes all users except "admin"
// TODO replace existing component-tests request logic with the CloudClient
// TODO test /api/check-auth, get user name and isAdmin == true
// TODO user registration, authorization and authentication etc
func TestLogin(t *testing.T) {
	cloud := getCloud()
	assert.Nil(t, cloud.parent.Cookie)
	assert.Nil(t, cloud.Login())
	cookie := cloud.parent.Cookie
	assert.NotNil(t, cookie)
	assert.Equal(t, 64, len(cookie.Value))
	assert.True(t, cookie.Expires.After(time.Now().AddDate(0, 0, 29)))
	assert.True(t, cookie.Expires.Before(time.Now().AddDate(0, 0, 31)))
}

func TestHappyPathDeployAndStop(t *testing.T) {
	cloud := getClientAndLogin(t)
	_, err := cloud.readApps()
	assert.Nil(t, err)
	assert.Nil(t, cloud.assertState("Uninitialized"))
	assert.Nil(t, cloud.startApp())
	assert.Nil(t, cloud.assertState("Available"))
	assert.Nil(t, cloud.stopApp())
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
	// TODO Adapt mock to make this test also usable for TEST profile
	cloud := getClientAndLogin(t) // TODO Maybe add a wait until all services are uninitialized? But only for prod
	cloud.appToOperateOn = tools.NginxSlowStart
	assert.Nil(t, cloud.startApp())
	// TODO Add assertion for Downloading and Stopping
	// assert.Nil(t, cloud.assertState("Downloading"))
	assert.Nil(t, cloud.assertState("Starting"))
	assert.Nil(t, cloud.assertState("Available"))
	// assert.Nil(t, cloud.assertState("Stopping"))
}

func onlyExecuteTestForProfile(t *testing.T, profileEnablingTheTest string) {
	setEnvProfile, _ := os.LookupEnv("PROFILE")
	if setEnvProfile != profileEnablingTheTest {
		t.Skip()
	}
}
