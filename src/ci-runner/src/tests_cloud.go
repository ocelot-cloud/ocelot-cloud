package src

import (
	"fmt"
	"ocelot/ci-runner/cli"
	"os/exec"
)

const (
	INITIAL_ADMIN_NAME_ENV     = "INITIAL_ADMIN_NAME=admin"
	INITIAL_ADMIN_PASSWORD_ENV = "INITIAL_ADMIN_PASSWORD=password"
)

func TestBackendCore() {
	printTaskDescription("Executing backend unit tests")
	defer Cleanup()
	cli.ExecuteInDir(backendAppsDir, "go test -v -count=1 .")
	cli.ExecuteInDir(backendAppsDir+"/download", "go test -v -count=1 ./...")
	cli.ExecuteInDir(backendAppsDir+"/yaml", "go test -v -count=1 ./...")
	cli.ExecuteInDir(backendSecurityDir, "go test -v -count=1 ./...")
	cli.ExecuteInDir(backendToolsDir, "go test -v -count=1 ./...")
}

func TestBackendComponentMocked() {
	printTaskDescription("Testing mocked backend component")
	defer Cleanup()
	cli.ExecuteInDir(backendDir, "rm -rf data")
	Build(Backend)
	// TODO Aggregate the envs
	StartDaemon(backendDir, "./backend", getTestProfileEnv(), getEnableDummyStacksEnv(true))
	cli.WaitUntilPortIsReady("8080")
	cli.ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 -tags functional ./...", getTestProfileEnv())
}

// TODO There are quite a lot of envs. Maybe I should refactor that into sth like "envs := getEnvs(...)".
func TestCloudAcceptance() {
	printTaskDescription("Testing acceptance")
	defer Cleanup()
	deployContainer(getEnableDummyStacksEnv(true))
	cli.ExecuteInDir(acceptanceTestsDir, cypressCommand)
}

func DeployContainer() {
	printTaskDescription("Running a production server")
	deployContainer()
}

func DeployContainerWithDummies() {
	printTaskDescription("Running a server using dummy stacks")
	deployContainer(getEnableDummyStacksEnv(true))
}

func deployContainer(additionalEnvs ...string) {
	exec.Command("/bin/sh", "-c", "docker network ls | grep -q ocelot-net || docker network create ocelot-net").Run()
	Build(DockerImage)
	envs := []string{
		"HOST=http://localhost",
		INITIAL_ADMIN_NAME_ENV,
		INITIAL_ADMIN_PASSWORD_ENV,
		"LOG_LEVEL=DEBUG",
	}
	envs = append(envs, additionalEnvs...)
	dockerCmd := fmt.Sprintf("bash -c '%s && docker logs -f ocelot-cloud'", ocelotContainerRunCommandDetached)
	StartDaemon(ocelotStackDir, dockerCmd, envs...)
	cli.WaitForIndexPageToBeReady(ocelotUrl)
}

func TestCi() {
	printTaskDescription("Running CI tests")
	TestBackend()
	TestFrontend()
	TestSecurity()
	TestCloudAcceptance()
	TestHubAll()
}

func TestSecurity() {
	printTaskDescription("Running security tests")
	defer Cleanup()
	deployContainer(getEnableDummyStacksEnv(true))
	cli.ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 -tags security ./...")
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
	cli.ExecuteInDir(backendDir, "rm -rf data")
	Build(Backend)
	StartDaemon(backendDir, "./backend", "USE_DUMMY_STACKS=true", "HOST=http://localhost:8080", INITIAL_ADMIN_NAME_ENV, INITIAL_ADMIN_PASSWORD_ENV, "ENABLE_DATA_WIPE_ENDPOINT=true")
	cli.WaitUntilPortIsReady("8080")
	cli.ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 -tags functional ./...", getProdProfileEnv())
}

func testBackendImageDownload() {
	cli.ExecuteInDir(backendAppsDir, "go test -v -count=1 -run TestDownloadProcessProviderReal", "IS_IMAGE_DOWNLOAD_TEST=true")
}

func printTaskDescription(text string) {
	cli.ColoredPrintln("\n=== %s ===\n", text)
}

func TestFrontend() {
	printTaskDescription("Testing Components In DevelopmentMode")
	defer Cleanup()
	cli.ExecuteInDir(backendDir, "rm -rf data")
	Build(Backend)
	StartDaemon(backendDir, "./backend", getTestProfileEnv(), getEnableDummyStacksEnv(true))
	cli.WaitUntilPortIsReady("8080")

	Build(Frontend)
	StartDaemon(frontendDir, "npm run serve", "VITE_APP_PROFILE="+TestProfile)
	cli.WaitForIndexPageToBeReady(frontendServerUrl)
	cli.ExecuteInDir(acceptanceTestsDir, cypressCommand, "CYPRESS_PROFILE="+TestProfile)
}

func testRunScript() {
	printTaskDescription("Testing run script")
	defer Cleanup()
	Build(DockerImage)
	cli.ExecuteInDir(scriptsDir, "bash run-dummy.sh")
	cli.WaitForIndexPageToBeReady(ocelotUrl)
	cli.ExecuteInDir(acceptanceTestsDir, cypressCommand)
}
