package src

import (
	"fmt"
	"github.com/ocelot-cloud/task-runner"
	"os/exec"
)

const (
	INITIAL_ADMIN_NAME_ENV     = "INITIAL_ADMIN_NAME=admin"
	INITIAL_ADMIN_PASSWORD_ENV = "INITIAL_ADMIN_PASSWORD=password"
)

func TestBackendCore() {
	tr.PrintTaskDescription("Executing backend unit tests")
	defer tr.Cleanup()
	tr.ExecuteInDir(backendAppsDir, "go test -v -count=1 .")
	tr.ExecuteInDir(backendAppsDir+"/download", "go test -v -count=1 ./...")
	tr.ExecuteInDir(backendAppsDir+"/yaml", "go test -v -count=1 ./...")
	tr.ExecuteInDir(backendRepoDir, "go test -v -count=1 ./...")
	tr.ExecuteInDir(backendToolsDir, "go test -v -count=1 ./...")
}

func TestBackendComponentMocked() {
	tr.PrintTaskDescription("Testing mocked backend component")
	defer tr.Cleanup()
	tr.ExecuteInDir(backendDir, "rm -rf data")
	Build(Backend)
	// TODO Aggregate the envs
	// TODO Dummy stacks should not be necessary when there are mocks used.
	tr.StartDaemon(backendDir, "./backend", getTestProfileEnv(), getEnableDummyStacksEnv(true))
	tr.WaitUntilPortIsReady("8080")
	tr.ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 -tags functional ./...", getTestProfileEnv())
}

// TODO There are quite a lot of envs. Maybe I should refactor that into sth like "envs := getEnvs(...)".
func TestCloudAcceptance() {
	tr.PrintTaskDescription("Testing acceptance")
	defer tr.Cleanup()
	deployContainer(getEnableDummyStacksEnv(true))
	tr.ExecuteInDir(acceptanceTestsDir, cypressCommand)
}

func DeployContainer() {
	tr.PrintTaskDescription("Running a production server")
	deployContainer()
}

func DeployContainerWithDummies() {
	tr.PrintTaskDescription("Running a server using dummy stacks")
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
	tr.StartDaemon(ocelotStackDir, dockerCmd, envs...)
	tr.WaitForWebPageToBeReady(ocelotUrl)
}

func TestCi() {
	tr.PrintTaskDescription("Running CI tests")
	TestBackend()
	TestFrontend()
	TestCloudAcceptance()
	TestIntegration()
	TestHubAll()
}

func TestIntegration() {
	tr.PrintTaskDescription("Testing integration between cloud and hub")
	defer tr.Cleanup()
	tr.StartDaemon(hubDir, "bash run-development-setup.sh")
	tr.WaitUntilPortIsReady("8082")
	tr.ExecuteInDir(cloudHubClientDir, "go test -v -count=1 ./...", getTestProfileEnv())
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
	tr.PrintTaskDescription("Testing PROD backend API with real docker service")
	defer tr.Cleanup()
	deployContainer(getEnableDummyStacksEnv(true), "ENABLE_DATA_WIPE_ENDPOINT=true")
	tr.ExecuteInDir(backendComponentTestsDir, "go test -v -count=1 -tags security ./...", getProdProfileEnv())
}

func testBackendImageDownload() {
	tr.ExecuteInDir(backendAppsDir, "go test -v -count=1 -run TestDownloadProcessProviderReal", "IS_IMAGE_DOWNLOAD_TEST=true")
}

func TestFrontend() {
	tr.PrintTaskDescription("Testing Components In DevelopmentMode")
	defer tr.Cleanup()
	tr.ExecuteInDir(backendDir, "rm -rf data")
	Build(Backend)
	tr.StartDaemon(backendDir, "./backend", getTestProfileEnv(), getEnableDummyStacksEnv(true))
	tr.WaitUntilPortIsReady("8080")

	Build(Frontend)
	tr.StartDaemon(frontendDir, "npm run serve", "VITE_APP_PROFILE="+TestProfile)
	tr.WaitForWebPageToBeReady(frontendServerUrl)
	tr.ExecuteInDir(acceptanceTestsDir, cypressCommand, "CYPRESS_PROFILE="+TestProfile)
}

func testRunScript() {
	tr.PrintTaskDescription("Testing run script")
	defer tr.Cleanup()
	Build(DockerImage)
	tr.ExecuteInDir(scriptsDir, "bash run-dummy.sh")
	tr.WaitForWebPageToBeReady(ocelotUrl)
	tr.ExecuteInDir(acceptanceTestsDir, cypressCommand)
}
