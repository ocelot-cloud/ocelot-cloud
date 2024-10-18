package apps

import (
	"github.com/ocelot-cloud/shared/assert"
	"ocelot/backend/apps/docker"
	"ocelot/backend/apps/download"
	"ocelot/backend/apps/vars"
	"ocelot/backend/apps/yaml"
	"ocelot/backend/tools"
	"testing"
)

var appToDeploy = tools.NginxDefault
var app2ToDeploy = tools.NginxDefault2

func createAppService() *appServiceImpl {
	vars.AppFileDir = vars.DummyAppAssetsDirForTests
	return &appServiceImpl{docker.ProvideServiceMock(), yaml.ProvideAppConfigService(), download.ProvideDownloaderMock(), make(map[string]appAction)}
}

func TestHappyPathDeployAndStop(t *testing.T) {
	appService := createAppService()

	err := appService.deployApp(appToDeploy)
	assert.Nil(t, err)

	currentStackInfo := appService.getAppStateInfo()
	assertState(t, currentStackInfo, appToDeploy, docker.Available)
	assertState(t, currentStackInfo, app2ToDeploy, docker.Uninitialized)

	err = appService.stopApp(appToDeploy)
	assert.Nil(t, err)

	infoAfterStop := appService.getAppStateInfo()
	assertState(t, infoAfterStop, appToDeploy, docker.Uninitialized)
	assertState(t, infoAfterStop, app2ToDeploy, docker.Uninitialized)
}

func assertState(t *testing.T, stackInfo map[string]docker.AppDetailsType, name string, state docker.AppState) {
	if _, ok := stackInfo[name]; ok {
		assert.Equal(t, state, stackInfo[name].State, "Stack was present but had wrong state.")
	} else {
		assert.Fail(t, "Stack was not present at all.")
	}
}

func TestAllStacksStop(t *testing.T) {
	appService := createAppService()
	assert.Nil(t, appService.deployApp(appToDeploy))
	assert.Nil(t, appService.deployApp(app2ToDeploy))

	infoAfterDeploy := appService.getAppStateInfo()
	assertState(t, infoAfterDeploy, appToDeploy, docker.Available)
	assertState(t, infoAfterDeploy, app2ToDeploy, docker.Available)

	assert.Nil(t, appService.stopAllApps())

	infoAfterStopAll := appService.getAppStateInfo()
	assertState(t, infoAfterStopAll, appToDeploy, docker.Uninitialized)
	assertState(t, infoAfterStopAll, app2ToDeploy, docker.Uninitialized)
}

func TestToDeploySameStackTwice(t *testing.T) {
	appService := createAppService()
	assert.Nil(t, appService.deployApp(appToDeploy))
	assert.Nil(t, appService.deployApp(appToDeploy))
}

func TestToNotRunningStack(t *testing.T) {
	appService := createAppService()
	err := appService.stopApp(appToDeploy)
	assert.NotNil(t, err)
	assert.Equal(t, "error - stopping app failed", err.Error())
}

func TestIgnoreStackInStackInfo(t *testing.T) {
	appService := createAppService()
	stackName := "ocelot-cloud"
	assert.Nil(t, appService.deployApp(stackName))

	stackStateInfo := appService.getAppStateInfo()
	if _, ok := stackStateInfo[stackName]; ok {
		logger.Error("%s was not ignored as expected.", stackName)
		t.Fail()
	}
}

func TestNginxCustomUrlPath(t *testing.T) {
	appService := createAppService()
	assert.Nil(t, appService.deployApp(tools.NginxCustomPath))
	actualUrlPath := getUrlPathForStack(t, appService, tools.NginxCustomPath)
	assert.Equal(t, "/custom-path", actualUrlPath)
}

func TestNginxDefaultUrlPath(t *testing.T) {
	appService := createAppService()
	assert.Nil(t, appService.deployApp(tools.NginxDefault))
	actualUrlPath := getUrlPathForStack(t, appService, tools.NginxDefault)
	assert.Equal(t, "/", actualUrlPath)
}

func getUrlPathForStack(t *testing.T, appService appServiceType, stackName string) string {
	stackStateInfo := appService.getAppStateInfo()
	if _, ok := stackStateInfo[stackName]; ok {
		return stackStateInfo[stackName].Path
	} else {
		t.Fatalf("stack '%s' was not found", stackName)
		return ""
	}
}

type appServiceTestApi struct {
	t          *testing.T
	appService appServiceType
	stackName  string
}

func (s *appServiceTestApi) deploy() *appServiceTestApi {
	assert.Nil(s.t, s.appService.deployApp(s.stackName))
	return s
}

func (s *appServiceTestApi) stop() *appServiceTestApi {
	assert.Nil(s.t, s.appService.stopApp(s.stackName))
	return s
}

func (s *appServiceTestApi) assertState(expectedState docker.AppState) *appServiceTestApi {
	stackStateInfo := s.appService.getAppStateInfo()
	if _, ok := stackStateInfo[s.stackName]; ok {
		assert.Equal(s.t, expectedState, stackStateInfo[s.stackName].State)
	} else {
		s.t.Fatalf("Stack '%s' not found", s.stackName)
	}
	return s
}

func TestHealthStateHandling(t *testing.T) {
	api := appServiceTestApi{t, createAppService(), tools.NginxSlowStart}
	api.deploy().assertState(docker.Starting).assertState(docker.Available)
	api.stop().assertState(docker.Uninitialized)
	api.deploy().assertState(docker.Starting).stop().assertState(docker.Uninitialized)
}

func TestDownloadStateHandling(t *testing.T) {
	api := appServiceTestApi{t, createAppService(), tools.NginxDownloading}
	api.deploy().assertState(docker.Downloading).assertState(docker.Starting).assertState(docker.Available)
	api.stop().assertState(docker.Uninitialized)
}
