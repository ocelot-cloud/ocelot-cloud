package src

import (
	"os/exec"
)

func testBackendCore() {
	printTaskDescription("Executing backend unit tests")
	defer Cleanup()
	ExecuteInDir(backendBusinessInternalDir, "go test -v -count=1 ./...")
	ExecuteInDir(backendSecurityInternalDir, "go test -v -count=1 ./...")
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
	StartDaemon(ocelotStackDir, ocelotContainerRunCommand, getEnableMocksEnv(true))
	WaitForIndexPageToBeReady(ocelotUrl)
	Build(Acceptance)
	ExecuteInDir(acceptanceTestsDir, cypressCommand)
}

func DeployLocally() {
	printTaskDescription("Running a production server")
	exec.Command("/bin/sh", "-c", "docker network ls | grep -q ocelot-net || docker network create ocelot-net").Run()
	Build(DockerImage)
	StartDaemon(ocelotStackDir, ocelotContainerRunCommandDetached)
	WaitForIndexPageToBeReady(ocelotUrl)
}

func TestCi() {
	printTaskDescription("Running CI tests")
	// Starting with the fastest tests, ending with slowest.

	// TODO backend units + mocked
	testBackendCore()
	TestBackendComponentMocked()

	// TODO development setup: test backend mocked + GUI
	checkComponentsWithTestProfile()

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
	// TODO Get rid of "DISABLE_SECURITY", it should always be enabled by default
	Build(Backend)
	StartDaemon(backendDir, "./backend", "DISABLE_SECURITY=true", "USE_DUMMY_STACKS=true")
	WaitUntilPortIsReady("localhost:8080")
	ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 ./...")
}

func testBackendImageDownload() {
	ExecuteInDir(backendBusinessInternalDir, "go test -v -count=1 -run TestDownloadProcessProviderReal", "IS_IMAGE_DOWNLOAD_TEST=true")
}

func printTaskDescription(text string) {
	ColoredPrintln("\n=== %s ===\n", text)
}

func checkComponentsWithTestProfile() {
	printTaskDescription("Testing Components In DevelopmentMode")
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
	printTaskDescription("Testing run script")
	defer Cleanup()
	Build(DockerImage)
	ExecuteInDir(scriptsDir, "bash run-dummy.sh")
	WaitForIndexPageToBeReady(ocelotUrl)
	Build(Acceptance)
	ExecuteInDir(acceptanceTestsDir, cypressCommand)
}
