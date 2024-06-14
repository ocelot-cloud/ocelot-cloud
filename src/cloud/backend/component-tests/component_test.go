package component_tests

import (
	"bytes"
	"encoding/json"
	"github.com/ocelot-cloud/shared"
	"net/http"
	"ocelot/backend/tools"
	"os"
	"testing"
	"time"
)

var logger = shared.ProvideLogger()

const endpoint = "http://localhost:8080/api/stacks/"
const stackOneName = tools.NginxDefault
const stackTwoName = tools.NginxDefault2

func TestHappyPathDeployAndStop(t *testing.T) {
	postJsonWithoutAssertions(endpoint+"stop", tools.StackInfo{stackOneName})
	postJsonWithoutAssertions(endpoint+"stop", tools.StackInfo{stackTwoName})

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

func postJsonWithoutAssertions(endpoint string, data tools.StackInfo) {
	jsonData, _ := json.Marshal(data)
	http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
}

func getAndRead(t *testing.T, endpoint string) []tools.ResponsePayloadDto {
	resp, err := http.Get(endpoint)
	shared.AssertNil(t, err)
	defer resp.Body.Close()

	var stackStates []tools.ResponsePayloadDto
	err = json.NewDecoder(resp.Body).Decode(&stackStates)
	shared.AssertNil(t, err)

	return stackStates
}

func assertState(t *testing.T, info []tools.ResponsePayloadDto, name string, state string) {
	for _, singleInfo := range info {
		if singleInfo.Name == name {
			shared.AssertEqual(t, state, singleInfo.State, "Stack '"+name+"' was present but had wrong state.")
			return
		}
	}
	shared.AssertFail(t, "Stack was not present at all.")
}

func postJSON(t *testing.T, endpoint string, stackName string) *http.Response {
	stackNameJson := tools.StackInfo{Name: stackName}
	jsonData, marshalErr := json.Marshal(stackNameJson)
	shared.AssertNil(t, marshalErr)
	resp, postErr := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	shared.AssertNil(t, postErr)
	shared.AssertEqual(t, 200, resp.StatusCode)
	return resp
}

func TestDeployStackNotExisting(t *testing.T) {
	postStackAndCheckResponse(t, "deploy", http.StatusInternalServerError)
}

func TestStopStackNotExisting(t *testing.T) {
	postStackAndCheckResponse(t, "stop", http.StatusInternalServerError)
}

func postStackAndCheckResponse(t *testing.T, action string, expectedHttpStatus int) {
	data := tools.StackInfo{"not-existing-stack"}
	jsonData, err := json.Marshal(data)
	shared.AssertNil(t, err)
	resp, err := http.Post(endpoint+action, "application/json", bytes.NewBuffer(jsonData))
	shared.AssertNil(t, err)
	shared.AssertEqual(t, expectedHttpStatus, resp.StatusCode)
}

func TestAbsenceOfCorsPolicyDisablingHeadersInResponse(t *testing.T) {
	AssertCorsHeaders(t, "", "", "")
}

func AssertCorsHeaders(t *testing.T, expectedAllowOrigin, expectedAllowMethods, expectedAllowHeaders string) {
	resp, err := http.Get("http://localhost:8080/api/stacks/read")
	shared.AssertNil(t, err)
	defer resp.Body.Close()

	allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
	shared.AssertEqual(t, expectedAllowOrigin, allowOrigin)

	allowMethods := resp.Header.Get("Access-Control-Allow-Methods")
	shared.AssertEqual(t, expectedAllowMethods, allowMethods)

	allowHeaders := resp.Header.Get("Access-Control-Allow-Headers")
	shared.AssertEqual(t, expectedAllowHeaders, allowHeaders)
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
	shared.AssertTrue(t, isCustomPathNginxPathOk)
	shared.AssertTrue(t, isDefaultNginxPathOk)
}

func TestNetworkCreationOnStackDeployment(t *testing.T) {
	dontExecuteTestForProfile(t, tools.BackendModeDependenciesMocked)

	_ = shared.ExecuteShellCommand("docker network ls | grep -q nginx-default-net || docker network rm nginx-default-net")
	postJSON(t, endpoint+"deploy", tools.NginxDefault)
	err := shared.ExecuteShellCommand("docker network ls | grep -q nginx-default-net")
	shared.AssertNil(t, err)
}

func TestWhetherCorsPolicyDisablingHeadersAreInResponse(t *testing.T) {
	onlyExecuteTestForProfile(t, tools.BackendModeDevelopmentSetup)
	AssertCorsHeaders(t, "*", "POST, GET, OPTIONS, PUT, DELETE", "Accept, Content-Type, Content-Length, Authorization")
}

func TestHealthStateOfSlowStartingStack(t *testing.T) {
	onlyExecuteTestForProfile(t, tools.BackendModeProdWithGui)

	postJsonWithoutAssertions(endpoint+"stop", tools.StackInfo{tools.NginxSlowStart})
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

func isStackInState(stackName string, expectedState string, responsePayload []tools.ResponsePayloadDto) bool {
	for _, singleInfo := range responsePayload {
		if singleInfo.Name == stackName && singleInfo.State == expectedState {
			return true
		}
	}
	return false
}

func onlyExecuteTestForProfile(t *testing.T, profileEnablingTheTest string) {
	setEnvProfile, _ := os.LookupEnv("BACKEND_COMPONENT_TEST_PROFILE")
	if setEnvProfile != profileEnablingTheTest {
		t.Skip()
	}
}

func dontExecuteTestForProfile(t *testing.T, profileDisablingTheTest string) {
	setEnvProfile, _ := os.LookupEnv("BACKEND_COMPONENT_TEST_PROFILE")
	if setEnvProfile == profileDisablingTheTest {
		t.Skip()
	}
}
