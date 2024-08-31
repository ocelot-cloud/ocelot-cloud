package apps

import (
	"github.com/ocelot-cloud/shared/assert"
	"ocelot/backend/tools"
	"testing"
)

var stackToDeploy = tools.NginxDefault
var stack2ToDeploy = tools.NginxDefault2

func createStackService() *appServiceImpl {
	appFileDir = DefaultStackFileDir
	return &appServiceImpl{provideServiceMock(), provideStackConfigService(appFileDir), provideDownloaderMock(), make(map[string]appAction)}
}

func TestHappyPathDeployAndStop(t *testing.T) {
	stackService := createStackService()

	err := stackService.deployApp(stackToDeploy)
	assert.Nil(t, err)

	currentStackInfo := stackService.getAppStateInfo()
	assertState(t, currentStackInfo, stackToDeploy, Available)
	assertState(t, currentStackInfo, stack2ToDeploy, Uninitialized)

	err = stackService.stopApp(stackToDeploy)
	assert.Nil(t, err)

	infoAfterStop := stackService.getAppStateInfo()
	assertState(t, infoAfterStop, stackToDeploy, Uninitialized)
	assertState(t, infoAfterStop, stack2ToDeploy, Uninitialized)
}

func assertState(t *testing.T, stackInfo map[string]appDetailsType, name string, state stackState) {
	if _, ok := stackInfo[name]; ok {
		assert.Equal(t, state, stackInfo[name].State, "Stack was present but had wrong state.")
	} else {
		assert.Fail(t, "Stack was not present at all.")
	}
}

func TestAllStacksStop(t *testing.T) {
	stackService := createStackService()
	assert.Nil(t, stackService.deployApp(stackToDeploy))
	assert.Nil(t, stackService.deployApp(stack2ToDeploy))

	infoAfterDeploy := stackService.getAppStateInfo()
	assertState(t, infoAfterDeploy, stackToDeploy, Available)
	assertState(t, infoAfterDeploy, stack2ToDeploy, Available)

	assert.Nil(t, stackService.StopAllStacks())

	infoAfterStopAll := stackService.getAppStateInfo()
	assertState(t, infoAfterStopAll, stackToDeploy, Uninitialized)
	assertState(t, infoAfterStopAll, stack2ToDeploy, Uninitialized)
}

func TestToDeploySameStackTwice(t *testing.T) {
	stackService := createStackService()
	assert.Nil(t, stackService.deployApp(stackToDeploy))
	assert.Nil(t, stackService.deployApp(stackToDeploy))
}

func TestToNotRunningStack(t *testing.T) {
	stackService := createStackService()
	err := stackService.stopApp(stackToDeploy)
	assert.NotNil(t, err)
	assert.Equal(t, "error - stopping stack failed", err.Error())
}

func TestIgnoreStackInStackInfo(t *testing.T) {
	stackService := createStackService()
	stackName := "ocelot-cloud"
	assert.Nil(t, stackService.deployApp(stackName))

	stackStateInfo := stackService.getAppStateInfo()
	if _, ok := stackStateInfo[stackName]; ok {
		logger.Error("%s was not ignored as expected.", stackName)
		t.Fail()
	}
}

func TestNginxCustomUrlPath(t *testing.T) {
	stackService := createStackService()
	assert.Nil(t, stackService.deployApp(tools.NginxCustomPath))
	actualUrlPath := getUrlPathForStack(t, stackService, tools.NginxCustomPath)
	assert.Equal(t, "/custom-path", actualUrlPath)
}

func TestNginxDefaultUrlPath(t *testing.T) {
	stackService := createStackService()
	assert.Nil(t, stackService.deployApp(tools.NginxDefault))
	actualUrlPath := getUrlPathForStack(t, stackService, tools.NginxDefault)
	assert.Equal(t, "/", actualUrlPath)
}

func getUrlPathForStack(t *testing.T, stackService appServiceType, stackName string) string {
	stackStateInfo := stackService.getAppStateInfo()
	if _, ok := stackStateInfo[stackName]; ok {
		return stackStateInfo[stackName].Path
	} else {
		t.Fatalf("stack '%s' was not found", stackName)
		return ""
	}
}

type StackServiceTestApi struct {
	t            *testing.T
	stackService appServiceType
	stackName    string
}

func (s *StackServiceTestApi) deploy() *StackServiceTestApi {
	assert.Nil(s.t, s.stackService.deployApp(s.stackName))
	return s
}

func (s *StackServiceTestApi) stop() *StackServiceTestApi {
	assert.Nil(s.t, s.stackService.stopApp(s.stackName))
	return s
}

func (s *StackServiceTestApi) assertState(expectedState stackState) *StackServiceTestApi {
	stackStateInfo := s.stackService.getAppStateInfo()
	if _, ok := stackStateInfo[s.stackName]; ok {
		assert.Equal(s.t, expectedState, stackStateInfo[s.stackName].State)
	} else {
		s.t.Fatalf("Stack '%s' not found", s.stackName)
	}
	return s
}

func TestHealthStateHandling(t *testing.T) {
	api := StackServiceTestApi{t, createStackService(), tools.NginxSlowStart}
	api.deploy().assertState(Starting).assertState(Available)
	api.stop().assertState(Uninitialized)
	api.deploy().assertState(Starting).stop().assertState(Uninitialized)
}

func TestDownloadStateHandling(t *testing.T) {
	api := StackServiceTestApi{t, createStackService(), tools.NginxDownloading}
	api.deploy().assertState(Downloading).assertState(Starting).assertState(Available)
	api.stop().assertState(Uninitialized)
}
