package src

import (
	"os"
	"os/exec"
	"path/filepath"
)

const ocelotContainerRunCommand = "docker-compose -p ocelot-cloud up"
const ocelotContainerRunCommandDetached = "docker-compose -p ocelot-cloud up -d"

var projectDir = GetProjectDir()
var scriptsDir = projectDir + "/scripts"
var ComponentDir = projectDir + "/components"
var backendDir = ComponentDir + "/backend"
var backendComponentTestsDir = backendDir + "/modules/component-tests"
var frontendDir = ComponentDir + "/frontend"
var acceptanceTestsDir = ComponentDir + "/acceptance-tests"
var ocelotStackDir = backendDir + "/stacks/core/ocelot-cloud"
var backendBusinessInternalDir = backendDir + "/modules/business/internal"

const BackendModeProduction = "production"
const BackendModeDependenciesMocked = "dependencies-mocked"
const BackendModeDevelopmentSetup = "development-setup"

const FrontendModeDevelopmentSetup = "development-setup"
const FrontendModeBackendMock = "backend-mock"

func GetProjectDir() string {
	devopsRunnerDir, _ := os.Getwd()
	componentDir := filepath.Dir(devopsRunnerDir)
	projectDir := filepath.Dir(componentDir)
	return projectDir
}

func BuildBackendAndFrontend() {
	printTestDescription("Building backend and frontend")
	ExecuteInDir(backendDir, "go build")
	ExecuteInDir(frontendDir, "npm install")
	ExecuteInDir(frontendDir, "npm run build")
}

func TestBackendTestingApi() {
	printTestDescription("Testing backend testing API")
	defer Cleanup()
	ExecuteInDir(backendBusinessInternalDir, "go test -v -count=1 ./...")
	ExecuteInDir(backendDir+"/modules/tools", "go test -v -count=1 ./...")
	ExecuteInDir(backendDir+"/modules/security/internal", "go test -v -count=1 ./...")
}

func TestBackendComponent() {
	testWithDefaultConfig()
	testCorsDisabling()
}

func testWithDefaultConfig() {
	printTestDescription("Testing backend component")
	defer Cleanup()
	ExecuteInDir(backendDir, "go build")
	StartBackendDaemonInMode(BackendModeProduction)
	ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 component_test.go", addBackendModeEnvPrefix(BackendModeProduction))
}

func addBackendModeEnvPrefix(profile string) string {
	return "BACKEND_COMPONENT_TEST_PROFILE=" + profile
}

func testCorsDisabling() {
	printTestDescription("Testing whether backend sets CORS headers to disable CORS policy")
	defer Cleanup()
	ExecuteInDir(backendDir, "go build")
	StartBackendDaemonInMode(BackendModeDevelopmentSetup)
	ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 -run=TestWhetherCorsPolicyDisablingHeadersAreInResponse ./...", addBackendModeEnvPrefix(BackendModeDevelopmentSetup))
}

func TestBackendComponentMocked() {
	printTestDescription("Testing mocked backend component")
	defer Cleanup()
	ExecuteInDir(backendDir, "go build")
	StartBackendDaemonInMode(BackendModeDependenciesMocked)
	ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 component_test.go", addBackendModeEnvPrefix(BackendModeDependenciesMocked))
}

func TestAcceptance() {
	printTestDescription("Testing acceptance")
	defer Cleanup()
	exec.Command("/bin/sh", "-c", "docker network ls | grep -q ocelot-net || docker network create ocelot-net").Run()
	ExecuteInDir(scriptsDir, "./build.sh")
	StartDaemon(ocelotStackDir, ocelotContainerRunCommand, "USE_DUMMY_STACKS=true")
	WaitForIndexPageToBeReady(ocelotUrl)
	ExecuteInDir(acceptanceTestsDir, "npm install")
	ExecuteInDir(acceptanceTestsDir, "npx cypress run --headless")
}

func RunProduction() {
	printTestDescription("Running a production server")
	exec.Command("/bin/sh", "-c", "docker network ls | grep -q ocelot-net || docker network create ocelot-net").Run()
	ExecuteInDir(scriptsDir, "./build.sh")
	StartDaemon(ocelotStackDir, ocelotContainerRunCommandDetached)
	WaitForIndexPageToBeReady(ocelotUrl)
}

func TestBackendFast() {
	printTestDescription("Running all fast backend tests")
	TestBackendTestingApi()
	TestBackendComponentMocked()
}

func TestCi() {
	printTestDescription("Running CI tests")
	// Starting with fastest tests, ending with slowest.
	TestBackendTestingApi()
	TestBackendComponentMocked()
	TestBackendComponent()
	TestFrontendFast()
	TestAcceptance()
}

func RunScheduledTests() {
	TestComponentsInDevelopmentSetupMode()
	TestBackendImageDownload()
	TestRunScript()
}

func TestBackendImageDownload() {
	ExecuteInDir(backendBusinessInternalDir, "go test -v -count=1 -run TestDownloadProcessProviderReal", "IS_IMAGE_DOWNLOAD_TEST=true")
}

func printTestDescription(text string) {
	ColoredPrint("\n=== %s ===\n", text)
}

func TestFrontendFast() {
	printTestDescription("Testing Frontend Fast")
	defer Cleanup()
	ExecuteInDir(frontendDir, "npm install")
	StartDaemon(frontendDir, "npm run serve", "VUE_APP_PROFILE="+FrontendModeBackendMock)
	WaitForIndexPageToBeReady(frontendServerUrl)
	ExecuteInDir(acceptanceTestsDir, "npm install")
	ExecuteInDir(acceptanceTestsDir, "npx cypress run --headless", "CYPRESS_PROFILE="+FrontendModeBackendMock)
}

func TestComponentsInDevelopmentSetupMode() {
	printTestDescription("Testing Components In DevelopmentMode")
	defer Cleanup()
	ExecuteInDir(backendDir, "go build")
	StartBackendDaemonInMode(BackendModeDevelopmentSetup)
	StartDaemon(frontendDir, "npm run serve", "VUE_APP_PROFILE="+FrontendModeDevelopmentSetup)
	WaitForIndexPageToBeReady(frontendServerUrl)
	ExecuteInDir(acceptanceTestsDir, "npx cypress run --headless", "CYPRESS_PROFILE="+FrontendModeDevelopmentSetup)
}

func TestRunScript() {
	printTestDescription("Testing Components In DevelopmentMode")
	defer Cleanup()
	ExecuteInDir(scriptsDir, "bash build.sh")
	ExecuteInDir(scriptsDir, "bash run-dummy.sh")
	WaitForIndexPageToBeReady(ocelotUrl)
	ExecuteInDir(acceptanceTestsDir, "npm install")
	ExecuteInDir(acceptanceTestsDir, "npx cypress run --headless")
}
