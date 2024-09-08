package apps

import (
	"errors"
	"fmt"
	"ocelot/backend/apps/image_download"
	"os"
)

type appServiceImpl struct {
	dockerService    dockerService
	appConfigService configServiceType
	downloadManager  image_download.DownloadManager
	lastActionOnApp  map[string]appAction
}

func provideAppServiceMocked(appConfigService configServiceType) appServiceType {
	return &appServiceImpl{provideServiceMock(), appConfigService, image_download.ProvideDownloaderMock(), make(map[string]appAction)}
}

func provideAppServiceReal(appConfigService configServiceType) appServiceType {
	return &appServiceImpl{&dockerServiceReal{}, appConfigService, image_download.ProvideDownloaderReal(), make(map[string]appAction)}
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
	State appState
	Path  string
}

type dockerService interface {
	deployApp(appName string) error
	stopApp(appName string) error
	getRunningAppStateInfo() (map[string]appDetailsType, error)
}

type configServiceType interface {
	getAppConfig(appName string) appConfig
}

func (sm *appServiceImpl) deployApp(appName string) error {
	sm.lastActionOnApp[appName] = Deploy
	sm.downloadManager.Download(appName)
	return sm.dockerService.deployApp(appName)
}

func (sm *appServiceImpl) getAppStateInfo() map[string]appDetailsType {
	logger.Trace("App state info was requested.")
	resultInfos, err := sm.dockerService.getRunningAppStateInfo()

	appsInDir, err := sm.appNamesInDirectory()
	if err != nil {
		logger.Error("error when reading app names from directory: %s", err.Error())
		return nil
	}

	resultInfos = sm.addUninitializedApps(resultInfos, appsInDir)
	delete(resultInfos, "ocelot-cloud")

	for appName, appDetail := range resultInfos {
		newPath := sm.appConfigService.getAppConfig(appName).UrlPath
		resultInfos[appName] = appDetailsType{appDetail.State, newPath}
	}

	downloadStates := sm.downloadManager.GetDownloadStates()
	for appName, appDetails := range resultInfos {
		if _, ok := downloadStates[appName]; ok {
			if downloadStates[appName] == image_download.Ongoing {
				resultInfos[appName] = appDetailsType{Downloading, appDetails.Path}
			} else if appDetails.State == Uninitialized && sm.lastActionOnApp[appName] == Deploy {
				resultInfos[appName] = appDetailsType{Starting, appDetails.Path}
			} else if appDetails.State != Uninitialized && sm.lastActionOnApp[appName] == Stop {
				resultInfos[appName] = appDetailsType{Stopping, appDetails.Path}
			}
		}
	}

	logAppStateInfo(resultInfos)
	return resultInfos
}

func logAppStateInfo(info map[string]appDetailsType) {
	var logString = ""
	currentIndex := 0
	for appName, appDetails := range info {
		if currentIndex == 0 {
			logString += fmt.Sprintf("\n  {%s: %s}", appName, appDetails.State.toString())
		} else {
			logString += fmt.Sprintf("\n,  {%s: %s}", appName, appDetails.State.toString())
		}
		currentIndex++
	}
	logger.Trace("App state info is returned: [%s\n]", logString)
}

func (sm *appServiceImpl) appNamesInDirectory() ([]string, error) {
	files, err := os.ReadDir(appFileDir)
	if err != nil {
		logger.Warn("Could not read app from directory '" + appFileDir + "': " + err.Error())
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

func (sm *appServiceImpl) addUninitializedApps(resultInfos map[string]appDetailsType, appsInDir []string) map[string]appDetailsType {
	for _, appName := range appsInDir {
		if _, ok := resultInfos[appName]; !ok {
			resultInfos[appName] = appDetailsType{Uninitialized, "/"}
		}
	}
	return resultInfos
}

func (sm *appServiceImpl) stopApp(appToStop string) error {
	sm.lastActionOnApp[appToStop] = Stop
	logger.Info("Stopping app: %s", appToStop)
	appStateInfo := sm.getAppStateInfo()
	var doesAppExist = false
	var existingApp appDetailsType
	for appName, appDetails := range appStateInfo {
		if appName == appToStop {
			doesAppExist = true
			existingApp = appDetails
			break
		}
	}
	if doesAppExist == false {
		return logAndCreateAppNotFoundError(appToStop)
	} else if !(existingApp.State == Starting || existingApp.State == Available || existingApp.State == Stopping) {
		logger.Warn("only 'Starting' and 'Available' apps can be stopped. State is: %s", existingApp.State.toString())
		return errors.New("error - stopping app failed")
	} else {
		logger.Debug("App does exist and is now stopped: %s", appToStop)
		return sm.dockerService.stopApp(appToStop)
	}
}

func (sm *appServiceImpl) stopAllApps() error {
	appStateInfo := sm.getAppStateInfo()

	for appName, appDetails := range appStateInfo {
		if appDetails.State == Starting || appDetails.State == Available {
			err := sm.stopApp(appName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
