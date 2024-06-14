package internal

import (
	"github.com/ocelot-cloud/shared"
	"ocelot/tools"
	"testing"
)

var stackToDeploy = tools.NginxDefault
var stack2ToDeploy = tools.NginxDefault2

func createStackService() *StackServiceImpl {
	StackFileDir = "../../../stacks/dummy"
	return &StackServiceImpl{ProvideServiceMock(), ProvideStackConfigService(StackFileDir), ProvideDownloadManagerMock(), make(map[string]StackAction)}
}

func TestHappyPathDeployAndStop(t *testing.T) {
	stackService := createStackService()

	err := stackService.DeployStack(stackToDeploy)
	shared.AssertNil(t, err)

	currentStackInfo := stackService.GetStackStateInfo()
	assertState(t, currentStackInfo, stackToDeploy, Available)
	assertState(t, currentStackInfo, stack2ToDeploy, Uninitialized)

	err = stackService.StopStack(stackToDeploy)
	shared.AssertNil(t, err)

	infoAfterStop := stackService.GetStackStateInfo()
	assertState(t, infoAfterStop, stackToDeploy, Uninitialized)
	assertState(t, infoAfterStop, stack2ToDeploy, Uninitialized)
}

func assertState(t *testing.T, stackInfo map[string]StackDetails, name string, state StackState) {
	if _, ok := stackInfo[name]; ok {
		shared.AssertEqual(t, state, stackInfo[name].State, "Stack was present but had wrong state.")
	} else {
		shared.AssertFail(t, "Stack was not present at all.")
	}
}

func TestAllStacksStop(t *testing.T) {
	stackService := createStackService()
	shared.AssertNil(t, stackService.DeployStack(stackToDeploy))
	shared.AssertNil(t, stackService.DeployStack(stack2ToDeploy))

	infoAfterDeploy := stackService.GetStackStateInfo()
	assertState(t, infoAfterDeploy, stackToDeploy, Available)
	assertState(t, infoAfterDeploy, stack2ToDeploy, Available)

	shared.AssertNil(t, stackService.StopAllStacks())

	infoAfterStopAll := stackService.GetStackStateInfo()
	assertState(t, infoAfterStopAll, stackToDeploy, Uninitialized)
	assertState(t, infoAfterStopAll, stack2ToDeploy, Uninitialized)
}

func TestToDeploySameStackTwice(t *testing.T) {
	stackService := createStackService()
	shared.AssertNil(t, stackService.DeployStack(stackToDeploy))
	shared.AssertNil(t, stackService.DeployStack(stackToDeploy))
}

func TestToNotRunningStack(t *testing.T) {
	stackService := createStackService()
	err := stackService.StopStack(stackToDeploy)
	shared.AssertNotNil(t, err)
	shared.AssertEqual(t, "error - stopping stack failed", err.Error())
}

func TestIgnoreStackInStackInfo(t *testing.T) {
	stackService := createStackService()
	stackName := "ocelot-cloud"
	shared.AssertNil(t, stackService.DeployStack(stackName))

	stackStateInfo := stackService.GetStackStateInfo()
	if _, ok := stackStateInfo[stackName]; ok {
		Logger.Error("%s was not ignored as expected.", stackName)
		t.Fail()
	}
}

func TestNginxCustomUrlPath(t *testing.T) {
	stackService := createStackService()
	shared.AssertNil(t, stackService.DeployStack(tools.NginxCustomPath))
	actualUrlPath := getUrlPathForStack(t, stackService, tools.NginxCustomPath)
	shared.AssertEqual(t, "/custom-path", actualUrlPath)
}

func TestNginxDefaultUrlPath(t *testing.T) {
	stackService := createStackService()
	shared.AssertNil(t, stackService.DeployStack(tools.NginxDefault))
	actualUrlPath := getUrlPathForStack(t, stackService, tools.NginxDefault)
	shared.AssertEqual(t, "/", actualUrlPath)
}

func getUrlPathForStack(t *testing.T, stackService StackService, stackName string) string {
	stackStateInfo := stackService.GetStackStateInfo()
	if _, ok := stackStateInfo[stackName]; ok {
		return stackStateInfo[stackName].Path
	} else {
		t.Fatalf("stack '%s' was not found", stackName)
		return ""
	}
}

type StackServiceTestApi struct {
	t            *testing.T
	stackService StackService
	stackName    string
}

func (s *StackServiceTestApi) deploy() *StackServiceTestApi {
	shared.AssertNil(s.t, s.stackService.DeployStack(s.stackName))
	return s
}

func (s *StackServiceTestApi) stop() *StackServiceTestApi {
	shared.AssertNil(s.t, s.stackService.StopStack(s.stackName))
	return s
}

func (s *StackServiceTestApi) assertState(expectedState StackState) *StackServiceTestApi {
	stackStateInfo := s.stackService.GetStackStateInfo()
	if _, ok := stackStateInfo[s.stackName]; ok {
		shared.AssertEqual(s.t, expectedState, stackStateInfo[s.stackName].State)
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
