package internal

import (
	"errors"
	"fmt"
	"os"
)

type StackServiceImpl struct {
	DockerService        DockerService
	StackConfigService   StackConfigService
	StackDownloadManager StackDownloadManager
	lastActionOnStack    map[string]StackAction
}

func ProvideStackServiceMocked(stackConfigService StackConfigService) StackService {
	return &StackServiceImpl{ProvideServiceMock(), stackConfigService, ProvideDownloadManagerMock(), make(map[string]StackAction)}
}

func ProvideStackServiceReal(stackConfigService StackConfigService) StackService {
	return &StackServiceImpl{&DockerServiceReal{}, stackConfigService, ProvideStackDownloadManagerReal(), make(map[string]StackAction)}
}

type StackAction int

const (
	Deploy StackAction = iota
	Stop
)

type StackService interface {
	DeployStack(stackName string) error
	StopStack(stackName string) error
	GetStackStateInfo() map[string]StackDetails
}

type StackDetails struct {
	State StackState
	Path  string
}

type DockerService interface {
	DeployStack(stackName string) error
	StopStack(stackName string) error
	GetRunningStackStateInfo() (map[string]StackDetails, error)
}

type StackConfigService interface {
	GetStackConfig(stackName string) StackConfig
}

type StackDownloadManager interface {
	GetStackDownloadStates() map[string]DownloadState
	DownloadStack(stackName string)
}

func (sm *StackServiceImpl) DeployStack(stackName string) error {
	sm.lastActionOnStack[stackName] = Deploy
	sm.StackDownloadManager.DownloadStack(stackName)
	return sm.DockerService.DeployStack(stackName)
}

func (sm *StackServiceImpl) GetStackStateInfo() map[string]StackDetails {
	Logger.Trace("Stack state info was requested.")
	resultInfos, err := sm.DockerService.GetRunningStackStateInfo()

	stacksInDir, err := sm.stackNamesInDirectory()
	if err != nil {
		Logger.Error("error when reading stack names from directory: %s", err.Error())
		return nil
	}

	resultInfos = sm.addUninitializedStacks(resultInfos, stacksInDir)
	delete(resultInfos, "ocelot-cloud")

	for stackName, stackDetail := range resultInfos {
		newPath := sm.StackConfigService.GetStackConfig(stackName).UrlPath
		resultInfos[stackName] = StackDetails{stackDetail.State, newPath}
	}

	downloadStates := sm.StackDownloadManager.GetStackDownloadStates()
	for stackName, stackDetails := range resultInfos {
		if _, ok := downloadStates[stackName]; ok {
			if downloadStates[stackName] == Ongoing {
				resultInfos[stackName] = StackDetails{Downloading, stackDetails.Path}
			} else if stackDetails.State == Uninitialized && sm.lastActionOnStack[stackName] == Deploy {
				resultInfos[stackName] = StackDetails{Starting, stackDetails.Path}
			} else if stackDetails.State != Uninitialized && sm.lastActionOnStack[stackName] == Stop {
				resultInfos[stackName] = StackDetails{Stopping, stackDetails.Path}
			}
		}
	}

	logStackStateInfo(resultInfos)
	return resultInfos
}

func logStackStateInfo(info map[string]StackDetails) {
	var logString = ""
	currentIndex := 0
	for stackName, stackDetails := range info {
		if currentIndex == 0 {
			logString += fmt.Sprintf("\n  {%s: %s}", stackName, stackDetails.State.String())
		} else {
			logString += fmt.Sprintf("\n,  {%s: %s}", stackName, stackDetails.State.String())
		}
		currentIndex++
	}
	Logger.Trace("Stack state info is returned: [%s\n]", logString)
}

func (sm *StackServiceImpl) stackNamesInDirectory() ([]string, error) {
	files, err := os.ReadDir(StackFileDir)
	if err != nil {
		Logger.Warn("Could not read stack from directory '" + StackFileDir + "': " + err.Error())
		return nil, err
	}

	var stackNames []string
	for _, f := range files {
		if f.IsDir() {
			stackNames = append(stackNames, f.Name())
		}
	}
	return stackNames, nil
}

func (sm *StackServiceImpl) addUninitializedStacks(resultInfos map[string]StackDetails, stacksInDir []string) map[string]StackDetails {
	for _, stackName := range stacksInDir {
		if _, ok := resultInfos[stackName]; !ok {
			resultInfos[stackName] = StackDetails{Uninitialized, "/"}
		}
	}
	return resultInfos
}

func (sm *StackServiceImpl) StopStack(stackToStopName string) error {
	sm.lastActionOnStack[stackToStopName] = Stop
	Logger.Info("Stopping stack: %s", stackToStopName)
	stackStateInfo := sm.GetStackStateInfo()
	var doesStackExist = false
	var existingStack StackDetails
	for stackName, stackDetails := range stackStateInfo {
		if stackName == stackToStopName {
			doesStackExist = true
			existingStack = stackDetails
			break
		}
	}
	if doesStackExist == false {
		return logAndCreateStackNotFoundError(stackToStopName)
	} else if !(existingStack.State == Starting || existingStack.State == Available || existingStack.State == Stopping) {
		Logger.Warn("only 'Starting' and 'Available' stacks can be stopped. State is: %s", existingStack.State.String())
		return errors.New("error - stopping stack failed")
	} else {
		Logger.Debug("Stack does exist and is now stopped: %s", stackToStopName)
		return sm.DockerService.StopStack(stackToStopName)
	}
}

func (sm *StackServiceImpl) StopAllStacks() error {
	stackStateInfo := sm.GetStackStateInfo()

	for stackName, stackDetails := range stackStateInfo {
		if stackDetails.State == Starting || stackDetails.State == Available {
			stackName := stackName
			err := sm.StopStack(stackName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
