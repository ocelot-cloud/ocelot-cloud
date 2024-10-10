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
	cli.PrintTaskDescription("Executing backend unit tests")
	defer cli.Cleanup()
	cli.ExecuteInDir(backendAppsDir, "go test -v -count=1 .")
	cli.ExecuteInDir(backendAppsDir+"/download", "go test -v -count=1 ./...")
	cli.ExecuteInDir(backendAppsDir+"/yaml", "go test -v -count=1 ./...")
	cli.ExecuteInDir(backendRepoDir, "go test -v -count=1 ./...")
	cli.ExecuteInDir(backendToolsDir, "go test -v -count=1 ./...")
}

func TestBackendComponentMocked() {
	cli.PrintTaskDescription("Testing mocked backend component")
	defer cli.Cleanup()
	cli.ExecuteInDir(backendDir, "rm -rf data")
	Build(Backend)
	// TODO Aggregate the envs
	// TODO Dummy stacks should not be necessary when there are mocks used.
	cli.StartDaemon(backendDir, "./backend", getTestProfileEnv(), getEnableDummyStacksEnv(true))
	cli.WaitUntilPortIsReady("8080")
	cli.ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 -tags functional ./...", getTestProfileEnv())
}

// TODO There are quite a lot of envs. Maybe I should refactor that into sth like "envs := getEnvs(...)".
func TestCloudAcceptance() {
	cli.PrintTaskDescription("Testing acceptance")
	defer cli.Cleanup()
	deployContainer(getEnableDummyStacksEnv(true))
	cli.ExecuteInDir(acceptanceTestsDir, cypressCommand)
}

func DeployContainer() {
	cli.PrintTaskDescription("Running a production server")
	deployContainer()
}

func DeployContainerWithDummies() {
	cli.PrintTaskDescription("Running a server using dummy stacks")
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
	cli.StartDaemon(ocelotStackDir, dockerCmd, envs...)
	cli.WaitForWebPageToBeReady(ocelotUrl)
}

func TestCi() {
	cli.PrintTaskDescription("Running CI tests")
	TestBackend()
	TestFrontend()
	TestCloudAcceptance()
	TestIntegration()
	TestHubAll()
}

func TestIntegration() {
	cli.PrintTaskDescription("Testing integration between cloud and hub")
	defer cli.Cleanup()
	cli.StartDaemon(hubDir, "bash run-development-setup.sh")
	cli.WaitUntilPortIsReady("8082")
	cli.ExecuteInDir(cloudHubClientDir, "go test -v -count=1 ./...", getTestProfileEnv())
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

// TODO Maybe dont seaprate between build tags "functional" and "security"? I want them to run together.
func TestProdBackendApi() {
	cli.PrintTaskDescription("Testing PROD backend API with real docker service")
	defer cli.Cleanup()
	deployContainer(getEnableDummyStacksEnv(true), "ENABLE_DATA_WIPE_ENDPOINT=true")
	cli.ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 -tags security ./...", getProdProfileEnv())
}

func testBackendImageDownload() {
	cli.ExecuteInDir(backendAppsDir, "go test -v -count=1 -run TestDownloadProcessProviderReal", "IS_IMAGE_DOWNLOAD_TEST=true")
}

func TestFrontend() {
	cli.PrintTaskDescription("Testing Components In DevelopmentMode")
	defer cli.Cleanup()
	cli.ExecuteInDir(backendDir, "rm -rf data")
	Build(Backend)
	cli.StartDaemon(backendDir, "./backend", getTestProfileEnv(), getEnableDummyStacksEnv(true))
	cli.WaitUntilPortIsReady("8080")

	Build(Frontend)
	cli.StartDaemon(frontendDir, "npm run serve", "VITE_APP_PROFILE="+TestProfile)
	cli.WaitForWebPageToBeReady(frontendServerUrl)
	cli.ExecuteInDir(acceptanceTestsDir, cypressCommand, "CYPRESS_PROFILE="+TestProfile)
}

func testRunScript() {
	cli.PrintTaskDescription("Testing run script")
	defer cli.Cleanup()
	Build(DockerImage)
	cli.ExecuteInDir(scriptsDir, "bash run-dummy.sh")
	cli.WaitForWebPageToBeReady(ocelotUrl)
	cli.ExecuteInDir(acceptanceTestsDir, cypressCommand)
}
