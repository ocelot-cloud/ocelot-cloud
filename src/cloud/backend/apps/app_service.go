package apps

import (
	"errors"
	"fmt"
	"os"
)

type appServiceImpl struct {
	dockerService     dockerService
	appConfigService  configServiceType
	downloadManager   downloadManager
	lastActionOnStack map[string]appAction
}

func provideAppServiceMocked(appConfigService configServiceType) appServiceType {
	return &appServiceImpl{provideServiceMock(), appConfigService, provideDownloaderMock(), make(map[string]appAction)}
}

func provideAppServiceReal(stackConfigService configServiceType) appServiceType {
	return &appServiceImpl{&dockerServiceReal{}, stackConfigService, provideDownloaderReal(), make(map[string]appAction)}
}

type appAction int

const (
	Deploy appAction = iota
	Stop
)

type appServiceType interface {
	deployApp(appName string) error
	stopApp(appName string) error
	getAppStateInfo() map[string]appDetailsType
}

type appDetailsType struct {
	State stackState
	Path  string
}

type dockerService interface {
	deployStack(appName string) error
	stopStack(appName string) error
	getRunningStackStateInfo() (map[string]appDetailsType, error)
}

type configServiceType interface {
	getAppConfig(appName string) appConfig
}

type downloadManager interface {
	getDownloadStates() map[string]downloadState
	download(appName string)
}

func (sm *appServiceImpl) deployApp(appName string) error {
	sm.lastActionOnStack[appName] = Deploy
	sm.downloadManager.download(appName)
	return sm.dockerService.deployStack(appName)
}

func (sm *appServiceImpl) getAppStateInfo() map[string]appDetailsType {
	logger.Trace("App state info was requested.")
	resultInfos, err := sm.dockerService.getRunningStackStateInfo()

	appsInDir, err := sm.appNamesInDirectory()
	if err != nil {
		logger.Error("error when reading stack names from directory: %s", err.Error())
		return nil
	}

	resultInfos = sm.addUninitializedStacks(resultInfos, appsInDir)
	delete(resultInfos, "ocelot-cloud")

	for appName, stackDetail := range resultInfos {
		newPath := sm.appConfigService.getAppConfig(appName).UrlPath
		resultInfos[appName] = appDetailsType{stackDetail.State, newPath}
	}

	downloadStates := sm.downloadManager.getDownloadStates()
	for stackName, stackDetails := range resultInfos {
		if _, ok := downloadStates[stackName]; ok {
			if downloadStates[stackName] == ongoing {
				resultInfos[stackName] = appDetailsType{Downloading, stackDetails.Path}
			} else if stackDetails.State == Uninitialized && sm.lastActionOnStack[stackName] == Deploy {
				resultInfos[stackName] = appDetailsType{Starting, stackDetails.Path}
			} else if stackDetails.State != Uninitialized && sm.lastActionOnStack[stackName] == Stop {
				resultInfos[stackName] = appDetailsType{Stopping, stackDetails.Path}
			}
		}
	}

	logAppStateInfo(resultInfos)
	return resultInfos
}

func logAppStateInfo(info map[string]appDetailsType) {
	var logString = ""
	currentIndex := 0
	for stackName, stackDetails := range info {
		if currentIndex == 0 {
			logString += fmt.Sprintf("\n  {%s: %s}", stackName, stackDetails.State.toString())
		} else {
			logString += fmt.Sprintf("\n,  {%s: %s}", stackName, stackDetails.State.toString())
		}
		currentIndex++
	}
	logger.Trace("App state info is returned: [%s\n]", logString)
}

func (sm *appServiceImpl) appNamesInDirectory() ([]string, error) {
	files, err := os.ReadDir(appFileDir)
	if err != nil {
		logger.Warn("Could not read stack from directory '" + appFileDir + "': " + err.Error())
		return nil, err
	}

	var appNames []string
	for _, f := range files {
		if f.IsDir() {
			appNames = append(appNames, f.Name())
		}
	}
	return appNames, nil
}

func (sm *appServiceImpl) addUninitializedStacks(resultInfos map[string]appDetailsType, stacksInDir []string) map[string]appDetailsType {
	for _, stackName := range stacksInDir {
		if _, ok := resultInfos[stackName]; !ok {
			resultInfos[stackName] = appDetailsType{Uninitialized, "/"}
		}
	}
	return resultInfos
}

func (sm *appServiceImpl) stopApp(stackToStopName string) error {
	sm.lastActionOnStack[stackToStopName] = Stop
	logger.Info("Stopping stack: %s", stackToStopName)
	stackStateInfo := sm.getAppStateInfo()
	var doesStackExist = false
	var existingStack appDetailsType
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
		logger.Warn("only 'Starting' and 'Available' stacks can be stopped. State is: %s", existingStack.State.toString())
		return errors.New("error - stopping stack failed")
	} else {
		logger.Debug("Stack does exist and is now stopped: %s", stackToStopName)
		return sm.dockerService.stopStack(stackToStopName)
	}
}

func (sm *appServiceImpl) StopAllStacks() error {
	stackStateInfo := sm.getAppStateInfo()

	for stackName, stackDetails := range stackStateInfo {
		if stackDetails.State == Starting || stackDetails.State == Available {
			stackName := stackName
			err := sm.stopApp(stackName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
