package src

import (
	"ocelot/ci-runner/cli"
	"os/exec"
)

const (
	INITIAL_ADMIN_NAME_ENV     = "INITIAL_ADMIN_NAME=admin"
	INITIAL_ADMIN_PASSWORD_ENV = "INITIAL_ADMIN_PASSWORD=password"
)

func TestBackendCore() {
	printTaskDescription("Executing backend unit tests")
	defer cli.Cleanup()
	cli.ExecuteInDir(backendAppsDir, "go test -v -count=1 .")
	cli.ExecuteInDir(backendAppsDir+"/download", "go test -v -count=1 ./...")
	cli.ExecuteInDir(backendAppsDir+"/yaml", "go test -v -count=1 ./...")
	cli.ExecuteInDir(backendSecurityDir, "go test -v -count=1 ./...")
	cli.ExecuteInDir(backendToolsDir, "go test -v -count=1 ./...")
}

func TestBackendComponentMocked() {
	printTaskDescription("Testing mocked backend component")
	defer cli.Cleanup()
	cli.ExecuteInDir(backendDir, "rm -rf data")
	Build(Backend)
	// TODO Aggregate the envs
	cli.StartDaemon(backendDir, "./backend", getTestProfileEnv(), getEnableDummyStacksEnv(true))
	cli.WaitUntilPortIsReady("8080")
	cli.ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 ./...", getTestProfileEnv())
}

// TODO There are quite a lot of envs. Maybe I should refactor that into sth like "envs := getEnvs(...)".
func TestCloudAcceptance() {
	printTaskDescription("Testing acceptance")
	defer cli.Cleanup()
	exec.Command("/bin/sh", "-c", "docker network ls | grep -q ocelot-net || docker network create ocelot-net").Run()
	Build(DockerImage)
	cli.StartDaemon(ocelotStackDir, ocelotContainerRunCommand, "USE_DUMMY_STACKS=true", "HOST=http://localhost", INITIAL_ADMIN_NAME_ENV, INITIAL_ADMIN_PASSWORD_ENV)
	cli.WaitForIndexPageToBeReady(ocelotUrl)
	cli.ExecuteInDir(acceptanceTestsDir, cypressCommand)
}

func DeployLocally() {
	printTaskDescription("Running a production server")
	exec.Command("/bin/sh", "-c", "docker network ls | grep -q ocelot-net || docker network create ocelot-net").Run()
	Build(DockerImage)
	cli.StartDaemon(ocelotStackDir, ocelotContainerRunCommandDetached, "HOST=http://localhost", INITIAL_ADMIN_NAME_ENV, INITIAL_ADMIN_PASSWORD_ENV)
	cli.WaitForIndexPageToBeReady(ocelotUrl)
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
	defer cli.Cleanup()
	cli.ExecuteInDir(backendDir, "rm -rf data")
	Build(Backend)
	cli.StartDaemon(backendDir, "./backend", "USE_DUMMY_STACKS=true", "HOST=http://localhost:8080", INITIAL_ADMIN_NAME_ENV, INITIAL_ADMIN_PASSWORD_ENV, "ENABLE_DATA_WIPE_ENDPOINT=true")
	cli.WaitUntilPortIsReady("8080")
	cli.ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 ./...", getProdProfileEnv())
}

func testBackendImageDownload() {
	cli.ExecuteInDir(backendAppsDir, "go test -v -count=1 -run TestDownloadProcessProviderReal", "IS_IMAGE_DOWNLOAD_TEST=true")
}

func printTaskDescription(text string) {
	cli.ColoredPrintln("\n=== %s ===\n", text)
}

func TestFrontend() {
	printTaskDescription("Testing Components In DevelopmentMode")
	defer cli.Cleanup()
	cli.ExecuteInDir(backendDir, "rm -rf data")
	Build(Backend)
	cli.StartDaemon(backendDir, "./backend", getTestProfileEnv(), getEnableDummyStacksEnv(true))
	cli.WaitUntilPortIsReady("8080")

	Build(Frontend)
	cli.StartDaemon(frontendDir, "npm run serve", "VITE_APP_PROFILE="+TestProfile)
	cli.WaitForIndexPageToBeReady(frontendServerUrl)
	cli.ExecuteInDir(acceptanceTestsDir, cypressCommand, "CYPRESS_PROFILE="+TestProfile)
}

func testRunScript() {
	printTaskDescription("Testing run script")
	defer cli.Cleanup()
	Build(DockerImage)
	cli.ExecuteInDir(scriptsDir, "bash run-dummy.sh")
	cli.WaitForIndexPageToBeReady(ocelotUrl)
	cli.ExecuteInDir(acceptanceTestsDir, cypressCommand)
}
