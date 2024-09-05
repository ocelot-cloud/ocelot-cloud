package src

import (
	"os/exec"
)

func testBackendCore() {
	printTaskDescription("Executing backend unit tests")
	defer Cleanup()
	ExecuteInDir(backendAppsDir, "go test -v -count=1 ./...")
	ExecuteInDir(backendSecurityDir, "go test -v -count=1 ./...")
	ExecuteInDir(backendToolsDir, "go test -v -count=1 ./...")
}

func testCorsDisabling() {
	printTaskDescription("Testing whether backend sets CORS headers to disable CORS policy")
	defer Cleanup()
	Build(Backend)
	StartDaemon(backendDir, "./backend", getTestProfileEnv())
	WaitUntilPortIsReady("localhost:8080")
	ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 -run=TestWhetherCorsPolicyDisablingHeadersAreInResponse ./...")
}

func TestBackendComponentMocked() {
	printTaskDescription("Testing mocked backend component")
	defer Cleanup()
	Build(Backend)
	StartDaemon(backendDir, "./backend", getTestProfileEnv(), getEnableDummyStacksEnv(true))
	WaitUntilPortIsReady("localhost:8080")
	ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 component_test.go")
}

func TestCloudAcceptance() {
	printTaskDescription("Testing acceptance")
	defer Cleanup()
	exec.Command("/bin/sh", "-c", "docker network ls | grep -q ocelot-net || docker network create ocelot-net").Run()
	Build(DockerImage)
	// TODO USE_DUMMY_STACKS should not be necessary here, since we use the mocks
	StartDaemon(ocelotStackDir, ocelotContainerRunCommand, "USE_DUMMY_STACKS=true", "HOST=http://localhost")
	WaitForIndexPageToBeReady(ocelotUrl)
	ExecuteInDir(acceptanceTestsDir, cypressCommand)
}

func DeployLocally() {
	printTaskDescription("Running a production server")
	exec.Command("/bin/sh", "-c", "docker network ls | grep -q ocelot-net || docker network create ocelot-net").Run()
	Build(DockerImage)
	StartDaemon(ocelotStackDir, ocelotContainerRunCommandDetached, "HOST=http://localhost")
	WaitForIndexPageToBeReady(ocelotUrl)
}

func TestCi() {
	printTaskDescription("Running CI tests")
	// Starting with the fastest tests, ending with slowest.

	// TODO backend units + mocked
	testBackendCore()
	TestBackendComponentMocked()

	// TODO development setup: test backend mocked + GUI
	TestCloudComponentsWithTestProfile()

	// TODO test backend no mocks, just API
	testProdBackendApi()
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
	printTaskDescription("Testing PROD backend API with real docker service")
	defer Cleanup()
	Build(Backend)
	StartDaemon(backendDir, "./backend", "USE_DUMMY_STACKS=true", "HOST=http://localhost:8080")
	WaitUntilPortIsReady("localhost:8080")
	ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 ./...")
}

func testBackendImageDownload() {
	ExecuteInDir(backendAppsDir, "go test -v -count=1 -run TestDownloadProcessProviderReal", "IS_IMAGE_DOWNLOAD_TEST=true")
}

func printTaskDescription(text string) {
	ColoredPrintln("\n=== %s ===\n", text)
}

func TestCloudComponentsWithTestProfile() {
	printTaskDescription("Testing Components In DevelopmentMode")
	defer Cleanup()
	Build(Backend)
	StartDaemon(backendDir, "./backend", getTestProfileEnv(), getEnableDummyStacksEnv(true))
	WaitUntilPortIsReady("localhost:8080")

	Build(Frontend)
	StartDaemon(frontendDir, "npm run serve", "VITE_APP_PROFILE="+TestProfile)
	WaitForIndexPageToBeReady(frontendServerUrl)
	ExecuteInDir(acceptanceTestsDir, cypressCommand, "CYPRESS_PROFILE="+TestProfile)
}

func testRunScript() {
	printTaskDescription("Testing run script")
	defer Cleanup()
	Build(DockerImage)
	ExecuteInDir(scriptsDir, "bash run-dummy.sh")
	WaitForIndexPageToBeReady(ocelotUrl)
	ExecuteInDir(acceptanceTestsDir, cypressCommand)
}
