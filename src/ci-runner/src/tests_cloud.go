package src

import (
	"os/exec"
)

const (
	INITIAL_ADMIN_NAME_ENV     = "INITIAL_ADMIN_NAME=admin"
	INITIAL_ADMIN_PASSWORD_ENV = "INITIAL_ADMIN_PASSWORD=password"
)

func TestBackendCore() {
	printTaskDescription("Executing backend unit tests")
	defer Cleanup()
	ExecuteInDir(backendAppsDir, "go test -v -count=1 ./...")
	ExecuteInDir(backendSecurityDir, "go test -v -count=1 ./...")
	ExecuteInDir(backendToolsDir, "go test -v -count=1 ./...")
}

func TestBackendComponentMocked() {
	printTaskDescription("Testing mocked backend component")
	defer Cleanup()
	ExecuteInDir(backendDir, "rm -rf data")
	Build(Backend)
	// TODO Aggregate the envs
	StartDaemon(backendDir, "./backend", getTestProfileEnv(), getEnableDummyStacksEnv(true))
	WaitUntilPortIsReady("localhost:8080")
	ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 component_test.go", getTestProfileEnv())
}

// TODO There are quite a lot of envs. Maybe I should refactor that into sth like "envs := getEnvs(...)".
func TestCloudAcceptance() {
	printTaskDescription("Testing acceptance")
	defer Cleanup()
	exec.Command("/bin/sh", "-c", "docker network ls | grep -q ocelot-net || docker network create ocelot-net").Run()
	Build(DockerImage)
	StartDaemon(ocelotStackDir, ocelotContainerRunCommand, "USE_DUMMY_STACKS=true", "HOST=http://localhost", INITIAL_ADMIN_NAME_ENV, INITIAL_ADMIN_PASSWORD_ENV)
	WaitForIndexPageToBeReady(ocelotUrl)
	ExecuteInDir(acceptanceTestsDir, cypressCommand)
}

func DeployLocally() {
	printTaskDescription("Running a production server")
	exec.Command("/bin/sh", "-c", "docker network ls | grep -q ocelot-net || docker network create ocelot-net").Run()
	Build(DockerImage)
	StartDaemon(ocelotStackDir, ocelotContainerRunCommandDetached, "HOST=http://localhost", INITIAL_ADMIN_NAME_ENV, INITIAL_ADMIN_PASSWORD_ENV)
	WaitForIndexPageToBeReady(ocelotUrl)
}

func TestCi() {
	printTaskDescription("Running CI tests")
	TestBackend()
	TestFrontend()
	TestCloudAcceptance()
	TestHubAll()
}

func TestBackend() {
	TestBackendCore()
	TestBackendComponentMocked()
	TestProdBackendApi()
}

func RunScheduledTests() {
	testRunScript()
	testBackendImageDownload()
}

func TestProdBackendApi() {
	printTaskDescription("Testing PROD backend API with real docker service")
	defer Cleanup()
	ExecuteInDir(backendDir, "rm -rf data")
	Build(Backend)
	StartDaemon(backendDir, "./backend", "USE_DUMMY_STACKS=true", "HOST=http://localhost:8080", INITIAL_ADMIN_NAME_ENV, INITIAL_ADMIN_PASSWORD_ENV)
	WaitUntilPortIsReady("localhost:8080")
	ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 ./...", getProdProfileEnv())
}

func testBackendImageDownload() {
	ExecuteInDir(backendAppsDir, "go test -v -count=1 -run TestDownloadProcessProviderReal", "IS_IMAGE_DOWNLOAD_TEST=true")
}

func printTaskDescription(text string) {
	ColoredPrintln("\n=== %s ===\n", text)
}

func TestFrontend() {
	printTaskDescription("Testing Components In DevelopmentMode")
	defer Cleanup()
	ExecuteInDir(backendDir, "rm -rf data")
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
