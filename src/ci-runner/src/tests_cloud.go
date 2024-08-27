package src

import (
	"os"
	"os/exec"
	"path/filepath"
)

const ocelotContainerRunCommand = "docker compose -p ocelot-cloud up"
const ocelotContainerRunCommandDetached = "docker compose -p ocelot-cloud up -d"
const cypressCommand = "npx cypress run --spec cypress/e2e/cloud.cy.ts --headless"

var projectDir = GetProjectDir()
var scriptsDir = projectDir + "/scripts"
var srcDir = projectDir + "/src"
var cloudDir = srcDir + "/cloud"
var backendDir = cloudDir + "/backend"
var backendComponentTestsDir = backendDir + "/component-tests"
var frontendDir = cloudDir + "/frontend"
var acceptanceTestsDir = cloudDir + "/acceptance-tests"
var ocelotStackDir = backendDir + "/stacks/core/ocelot-cloud"
var backendBusinessInternalDir = backendDir + "/business/internal"
var backendSecurityInternalDir = backendDir + "/security/internal"
var hubDir = srcDir + "/hub"

// TestProfile There is also the "PROD" profile, but it should be used automatically if no profile is given.
const TestProfile = "TEST"

func GetProjectDir() string {
	devopsRunnerDir, _ := os.Getwd()
	src := filepath.Dir(devopsRunnerDir)
	return filepath.Dir(src)
}

func testBackendCore() {
	printTestDescription("Testing backend testing API")
	defer Cleanup()
	ExecuteInDir(backendBusinessInternalDir, "go test -v -count=1 ./...")
	ExecuteInDir(backendSecurityInternalDir, "go test -v -count=1 ./...")
}

func testWithDefaultConfig() {
	printTestDescription("Testing backend component")
	defer Cleanup()
	Build(Backend)
	StartDaemon(backendDir, "./backend", "DISABLE_SECURITY=true")
	WaitUntilPortIsReady("localhost:8080")
	ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 component_test.go")
}

func testCorsDisabling() {
	printTestDescription("Testing whether backend sets CORS headers to disable CORS policy")
	defer Cleanup()
	Build(Backend)
	StartDaemon(backendDir, "./backend", getTestProfileEnv())
	WaitUntilPortIsReady("localhost:8080")

	ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 -run=TestWhetherCorsPolicyDisablingHeadersAreInResponse ./...")
}

func TestBackendComponentMocked() {
	printTestDescription("Testing mocked backend component")
	defer Cleanup()
	Build(Backend)
	StartDaemon(backendDir, "./backend", getTestProfileEnv(), getEnableDummyStacksEnv(true))
	WaitUntilPortIsReady("localhost:8080")
	ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 component_test.go")
}

func getTestProfileEnv() string {
	return "PROFILE=" + TestProfile
}

func getEnableDummyStacksEnv(enabled bool) string {
	prefix := "USE_DUMMY_STACKS="
	if enabled {
		return prefix + "true"
	} else {
		return prefix + "false"
	}
}

func TestCloudAcceptance() {
	printTestDescription("Testing acceptance")
	defer Cleanup()
	exec.Command("/bin/sh", "-c", "docker network ls | grep -q ocelot-net || docker network create ocelot-net").Run()
	Build(DockerImage)
	StartDaemon(ocelotStackDir, ocelotContainerRunCommand, getEnableMocksEnv(true))
	WaitForIndexPageToBeReady(ocelotUrl)
	Build(Acceptance)
	ExecuteInDir(acceptanceTestsDir, cypressCommand)
}

func getEnableMocksEnv(enabled bool) string {
	prefix := "ENABLE_MOCKS="
	if enabled {
		return prefix + "true"
	} else {
		return prefix + "false"
	}
}

func DeployLocally() {
	printTestDescription("Running a production server")
	exec.Command("/bin/sh", "-c", "docker network ls | grep -q ocelot-net || docker network create ocelot-net").Run()
	Build(DockerImage)
	StartDaemon(ocelotStackDir, ocelotContainerRunCommandDetached)
	WaitForIndexPageToBeReady(ocelotUrl)
}

func TestCi() {
	printTestDescription("Running CI tests")
	// Starting with the fastest tests, ending with slowest.

	// TODO backend units + mocked
	testBackendCore()
	TestBackendComponentMocked()

	// TODO development setup: test backend mocked + GUI
	checkComponentsWithTestProfile()

	// TODO test backend no mocks, just API
	testProdBackendApi()
	testWithDefaultConfig()
	testCorsDisabling()

	// TODO acceptance, backend mocked
	TestCloudAcceptance()

	TestHubAll()
}

func RunScheduledTests() {
	testRunScript()
	testBackendImageDownload()
}

func testProdBackendApi() {
	printTestDescription("Testing PROD backend")
	defer Cleanup()
	// TODO Get rid of "disable security" and the other CLI args
	Build(Backend)
	StartDaemon(backendDir, "./backend", getEnableMocksEnv(false), "DISABLE_SECURITY=true")
	ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 ./...")
}

func testBackendImageDownload() {
	ExecuteInDir(backendBusinessInternalDir, "go test -v -count=1 -run TestDownloadProcessProviderReal", "IS_IMAGE_DOWNLOAD_TEST=true")
}

func printTestDescription(text string) {
	ColoredPrintln("\n=== %s ===\n", text)
}

func checkComponentsWithTestProfile() {
	printTestDescription("Testing Components In DevelopmentMode")
	defer Cleanup()
	Build(Backend)
	StartDaemon(backendDir, "./backend", getTestProfileEnv(), getEnableDummyStacksEnv(true))
	WaitUntilPortIsReady("localhost:8080")

	Build(Frontend)
	StartDaemon(frontendDir, "npm run serve", "VITE_APP_PROFILE="+TestProfile)
	WaitForIndexPageToBeReady(frontendServerUrl)
	Build(Acceptance)
	ExecuteInDir(acceptanceTestsDir, cypressCommand, "CYPRESS_PROFILE="+TestProfile)
}

func testRunScript() {
	printTestDescription("Testing run script")
	defer Cleanup()
	Build(DockerImage)
	ExecuteInDir(scriptsDir, "bash run-dummy.sh")
	WaitForIndexPageToBeReady(ocelotUrl)
	Build(Acceptance)
	ExecuteInDir(acceptanceTestsDir, cypressCommand)
}
