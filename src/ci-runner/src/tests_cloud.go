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

const BackendModeProduction = "production"
const BackendModeDependenciesMocked = "dependencies-mocked"
const BackendModeDevelopmentSetup = "development-setup"

const FrontendModeDevelopmentSetup = "development-setup"
const FrontendModeBackendMock = "backend-mock"

func GetProjectDir() string {
	devopsRunnerDir, _ := os.Getwd()
	src := filepath.Dir(devopsRunnerDir)
	return filepath.Dir(src)
}

func BuildBackendAndFrontend() {
	printTestDescription("Building backend and frontend")
	Build(Backend)
	Build(Frontend)
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
	StartBackendDaemon(BackendModeProduction, "PROD")
	ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 component_test.go", addBackendProfileEnvPrefix(BackendModeProduction))
}

func addBackendProfileEnvPrefix(profile string) string {
	return "BACKEND_COMPONENT_TEST_PROFILE=" + profile
}

func testCorsDisabling() {
	printTestDescription("Testing whether backend sets CORS headers to disable CORS policy")
	defer Cleanup()
	Build(Backend)
	StartBackendDaemon(BackendModeDevelopmentSetup, "TEST")
	ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 -run=TestWhetherCorsPolicyDisablingHeadersAreInResponse ./...", addBackendProfileEnvPrefix(BackendModeDevelopmentSetup))
}

func TestBackendComponentMocked() {
	printTestDescription("Testing mocked backend component")
	defer Cleanup()
	Build(Backend)
	StartBackendDaemon(BackendModeDependenciesMocked, "TEST")
	ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 component_test.go", addBackendProfileEnvPrefix(BackendModeDependenciesMocked))
}

func TestCloudAcceptance() {
	printTestDescription("Testing acceptance")
	defer Cleanup()
	exec.Command("/bin/sh", "-c", "docker network ls | grep -q ocelot-net || docker network create ocelot-net").Run()
	Build(DockerImage)
	StartDaemon(ocelotStackDir, ocelotContainerRunCommand, "PROFILE=PROD")
	WaitForIndexPageToBeReady(ocelotUrl)
	Build(Acceptance)
	ExecuteInDir(acceptanceTestsDir, cypressCommand)
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
	testComponentsInDevelopmentSetupMode()

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
	// TODO abstract profiles into variable
	// TODO Adapt profiles in frontend and acceptance
	Build(Backend)
	StartDaemon(backendDir, "./backend -enable-dummy-stacks -disable-security", "PROFILE=PROD", "ENABLE_MOCKS=false")
	ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 ./...")
}

func testBackendImageDownload() {
	ExecuteInDir(backendBusinessInternalDir, "go test -v -count=1 -run TestDownloadProcessProviderReal", "IS_IMAGE_DOWNLOAD_TEST=true")
}

func printTestDescription(text string) {
	ColoredPrintln("\n=== %s ===\n", text)
}

func testComponentsInDevelopmentSetupMode() {
	printTestDescription("Testing Components In DevelopmentMode")
	defer Cleanup()
	Build(Backend)
	StartBackendDaemon(BackendModeDevelopmentSetup, "TEST")
	Build(Frontend)
	StartDaemon(frontendDir, "npm run serve", "VUE_APP_PROFILE="+FrontendModeDevelopmentSetup)
	WaitForIndexPageToBeReady(frontendServerUrl)
	Build(Acceptance)
	ExecuteInDir(acceptanceTestsDir, cypressCommand, "CYPRESS_PROFILE="+FrontendModeDevelopmentSetup)
}

func testRunScript() {
	printTestDescription("Testing Components In DevelopmentMode")
	defer Cleanup()
	Build(DockerImage)
	ExecuteInDir(scriptsDir, "bash run-dummy.sh")
	WaitForIndexPageToBeReady(ocelotUrl)
	Build(Acceptance)
	ExecuteInDir(acceptanceTestsDir, cypressCommand)
}
