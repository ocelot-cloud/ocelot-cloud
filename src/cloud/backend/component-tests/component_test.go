package component_tests

import (
	"bytes"
	"encoding/json"
	"github.com/ocelot-cloud/shared"
	"github.com/ocelot-cloud/shared/assert"
	"github.com/ocelot-cloud/shared/utils"
	"net/http"
	"ocelot/backend/tools"
	"os"
	"testing"
	"time"
)

var logger = tools.Logger

const (
	backendBaseUrl = "http://localhost:8080"
	endpoint       = backendBaseUrl + "/api/stacks/"
	stackOneName   = tools.NginxDefault
	stackTwoName   = tools.NginxDefault2
	TestProfile    = "TEST"
	ProdProfile    = "PROD"
)

func TestHappyPathDeployAndStop(t *testing.T) {
	postJsonWithoutAssertions(endpoint+"stop", utils.SingleString{stackOneName})
	postJsonWithoutAssertions(endpoint+"stop", utils.SingleString{stackTwoName})

	responsePayloadsBeforeDeploy := getAndRead(t, endpoint+"read")
	assertState(t, responsePayloadsBeforeDeploy, stackOneName, "Uninitialized")
	assertState(t, responsePayloadsBeforeDeploy, stackTwoName, "Uninitialized")

	postJSON(t, endpoint+"deploy", stackOneName)
	time.Sleep(3 * time.Second)
	responsePayloadsAfterDeploy := getAndRead(t, endpoint+"read")
	assertState(t, responsePayloadsAfterDeploy, stackOneName, "Available")
	assertState(t, responsePayloadsAfterDeploy, stackTwoName, "Uninitialized")

	postJSON(t, endpoint+"stop", stackOneName)
	responsePayloadsAfterStop := getAndRead(t, endpoint+"read")
	assertState(t, responsePayloadsAfterStop, stackOneName, "Uninitialized")
	assertState(t, responsePayloadsAfterStop, stackTwoName, "Uninitialized")
}

func postJsonWithoutAssertions(endpoint string, data utils.SingleString) {
	jsonData, _ := json.Marshal(data)
	http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
}

func getAndRead(t *testing.T, endpoint string) []tools.AppInfo {
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	assert.Nil(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "auth",
		Value: "valid",
	})

	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	var stackStates []tools.AppInfo
	err = json.NewDecoder(resp.Body).Decode(&stackStates)
	assert.Nil(t, err)

	return stackStates
}

func assertState(t *testing.T, info []tools.AppInfo, name string, state string) {
	for _, singleInfo := range info {
		if singleInfo.Name == name {
			assert.Equal(t, state, singleInfo.State, "Stack '"+name+"' was present but had wrong state.")
			return
		}
	}
	assert.Fail(t, "Stack was not present at all.")
}

func postJSON(t *testing.T, endpoint string, stackName string) *http.Response {
	stackNameJson := utils.SingleString{stackName}
	jsonData, marshalErr := json.Marshal(stackNameJson)
	assert.Nil(t, marshalErr)

	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "auth",
		Value: "valid",
	})

	resp, postErr := client.Do(req)
	assert.Nil(t, postErr)
	assert.Equal(t, 200, resp.StatusCode)
	return resp
}

func TestDeployStackNotExisting(t *testing.T) {
	postStackAndCheckResponse(t, "deploy", http.StatusInternalServerError)
}

func TestStopStackNotExisting(t *testing.T) {
	postStackAndCheckResponse(t, "stop", http.StatusInternalServerError)
}

func postStackAndCheckResponse(t *testing.T, action string, expectedHttpStatus int) {
	data := utils.SingleString{"not-existing-stack"}
	jsonData, err := json.Marshal(data)
	assert.Nil(t, err)

	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint+action, bytes.NewBuffer(jsonData))
	assert.Nil(t, err)

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Add the "auth=valid" cookie to the request
	req.AddCookie(&http.Cookie{
		Name:  "auth",
		Value: "valid",
	})

	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, expectedHttpStatus, resp.StatusCode)
}

func TestAbsenceOfCorsPolicyDisablingHeadersInResponse(t *testing.T) {
	onlyExecuteTestForProfile(t, ProdProfile)
	AssertCorsHeaders(t, "", "", "", "")
}

func AssertCorsHeaders(t *testing.T, expectedAllowOrigin, expectedAllowMethods, expectedAllowHeaders, expectedAllowCredentials string) {
	req, err := http.NewRequest("GET", backendBaseUrl+"/api/stacks/read", nil)
	if err != nil {
		logger.Error("error: %v", err)
		t.Fail()
	}
	req.Header.Set("Origin", backendBaseUrl)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("error: %v", err)
		t.Fail()
	}
	defer resp.Body.Close()

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
	responsePayloads := getAndRead(t, endpoint+"read")
	isCustomPathNginxPathOk := false
	isDefaultNginxPathOk := false
	for _, responsePayload := range responsePayloads {
		if responsePayload.Name == tools.NginxCustomPath && responsePayload.UrlPath == "/custom-path" {
			isCustomPathNginxPathOk = true
		} else if responsePayload.Name == tools.NginxDefault && responsePayload.UrlPath == "/" {
			isDefaultNginxPathOk = true
		}
	}
	assert.True(t, isCustomPathNginxPathOk)
	assert.True(t, isDefaultNginxPathOk)
}

func TestNetworkCreationOnStackDeployment(t *testing.T) {
	onlyExecuteTestForProfile(t, ProdProfile)

	_ = shared.ExecuteShellCommand("docker network ls | grep -q nginx-default-net || docker network rm nginx-default-net")
	postJSON(t, endpoint+"deploy", tools.NginxDefault)
	err := shared.ExecuteShellCommand("docker network ls | grep -q nginx-default-net")
	assert.Nil(t, err)
}

func TestWhetherCorsPolicyDisablingHeadersAreInResponse(t *testing.T) {
	onlyExecuteTestForProfile(t, TestProfile)
	AssertCorsHeaders(t, backendBaseUrl, "POST, GET, OPTIONS, PUT, DELETE", "Accept, Content-Type, Content-Length, Authorization", "true")
}

func TestHealthStateOfSlowStartingStack(t *testing.T) {
	onlyExecuteTestForProfile(t, ProdProfile)

	postJsonWithoutAssertions(endpoint+"stop", utils.SingleString{tools.NginxSlowStart})
	logger.Info("Deploying stack '%s'", tools.NginxSlowStart)
	postJSON(t, endpoint+"deploy", tools.NginxSlowStart)

	assertWithinLongerTimeRangeThatStackStateBecomesExpectedState(t, tools.NginxSlowStart, "Starting")
	assertWithinLongerTimeRangeThatStackStateBecomesExpectedState(t, tools.NginxSlowStart, "Available")
}

func assertWithinLongerTimeRangeThatStackStateBecomesExpectedState(t *testing.T, stackName string, expectedState string) {
	const maxAttempts = 30
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		responsePayload := getAndRead(t, endpoint+"read")
		if isStackInState(stackName, expectedState, responsePayload) {
			return
		}
		logger.Info("Attempt %v: Stack '%s' is not in state '%s' yet. Re-try in one second...", attempt, stackName, expectedState)
		time.Sleep(1 * time.Second)
	}
	t.Fail()
}

func isStackInState(stackName string, expectedState string, responsePayload []tools.AppInfo) bool {
	for _, singleInfo := range responsePayload {
		if singleInfo.Name == stackName && singleInfo.State == expectedState {
			return true
		}
	}
	return false
}

func onlyExecuteTestForProfile(t *testing.T, profileEnablingTheTest string) {
	setEnvProfile, _ := os.LookupEnv("PROFILE")
	if setEnvProfile != profileEnablingTheTest {
		t.Skip()
	}
}
