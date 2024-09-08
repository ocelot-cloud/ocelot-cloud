package apps

import (
	"errors"
	"fmt"
	"ocelot/backend/apps/docker"
	"ocelot/backend/apps/download"
	"ocelot/backend/apps/vars"
	"ocelot/backend/apps/yaml"
	"os"
)

type appServiceImpl struct {
	dockerService    dockerService
	appConfigService yaml.ConfigServiceType
	downloadManager  download.DownloadManager
	lastActionOnApp  map[string]appAction
}

func provideAppServiceMocked(appConfigService yaml.ConfigServiceType) appServiceType {
	return &appServiceImpl{docker.ProvideServiceMock(), appConfigService, download.ProvideDownloaderMock(), make(map[string]appAction)}
}

func provideAppServiceReal(appConfigService yaml.ConfigServiceType) appServiceType {
	return &appServiceImpl{&docker.DockerServiceReal{}, appConfigService, download.ProvideDownloaderReal(), make(map[string]appAction)}
}

type appAction int

const (
	Deploy appAction = iota
	Stop
)

type appServiceType interface {
	deployApp(appName string) error
	stopApp(appName string) error
	getAppStateInfo() map[string]docker.AppDetailsType
}

type dockerService interface {
	DeployApp(appName string) error
	StopApp(appName string) error
	GetRunningAppStateInfo() (map[string]docker.AppDetailsType, error)
}

func (sm *appServiceImpl) deployApp(appName string) error {
	sm.lastActionOnApp[appName] = Deploy
	sm.downloadManager.Download(appName)
	return sm.dockerService.DeployApp(appName)
}

func (sm *appServiceImpl) getAppStateInfo() map[string]docker.AppDetailsType {
	logger.Trace("App state info was requested.")
	resultInfos, err := sm.dockerService.GetRunningAppStateInfo()

	appsInDir, err := sm.appNamesInDirectory()
	if err != nil {
		logger.Error("error when reading app names from directory: %s", err.Error())
		return nil
	}

	resultInfos = sm.addUninitializedApps(resultInfos, appsInDir)
	delete(resultInfos, "ocelot-cloud")

	for appName, appDetail := range resultInfos {
		newPath := sm.appConfigService.GetAppConfig(appName).UrlPath
		resultInfos[appName] = docker.AppDetailsType{appDetail.State, newPath}
	}

	downloadStates := sm.downloadManager.GetDownloadStates()
	for appName, appDetails := range resultInfos {
		if _, ok := downloadStates[appName]; ok {
			if downloadStates[appName] == download.Ongoing {
				resultInfos[appName] = docker.AppDetailsType{docker.Downloading, appDetails.Path}
			} else if appDetails.State == docker.Uninitialized && sm.lastActionOnApp[appName] == Deploy {
				resultInfos[appName] = docker.AppDetailsType{docker.Starting, appDetails.Path}
			} else if appDetails.State != docker.Uninitialized && sm.lastActionOnApp[appName] == Stop {
				resultInfos[appName] = docker.AppDetailsType{docker.Stopping, appDetails.Path}
			}
		}
	}

	logAppStateInfo(resultInfos)
	return resultInfos
}

func logAppStateInfo(info map[string]docker.AppDetailsType) {
	var logString = ""
	currentIndex := 0
	for appName, appDetails := range info {
		if currentIndex == 0 {
			logString += fmt.Sprintf("\n  {%s: %s}", appName, appDetails.State.ToString())
		} else {
			logString += fmt.Sprintf("\n,  {%s: %s}", appName, appDetails.State.ToString())
		}
		currentIndex++
	}
	logger.Trace("App state info is returned: [%s\n]", logString)
}

func (sm *appServiceImpl) appNamesInDirectory() ([]string, error) {
	files, err := os.ReadDir(vars.AppFileDir)
	if err != nil {
		logger.Warn("Could not read app from directory '" + vars.AppFileDir + "': " + err.Error())
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

func (sm *appServiceImpl) addUninitializedApps(resultInfos map[string]docker.AppDetailsType, appsInDir []string) map[string]docker.AppDetailsType {
	for _, appName := range appsInDir {
		if _, ok := resultInfos[appName]; !ok {
			resultInfos[appName] = docker.AppDetailsType{docker.Uninitialized, "/"}
		}
	}
	return resultInfos
}

func (sm *appServiceImpl) stopApp(appToStop string) error {
	sm.lastActionOnApp[appToStop] = Stop
	logger.Info("Stopping app: %s", appToStop)
	appStateInfo := sm.getAppStateInfo()
	var doesAppExist = false
	var existingApp docker.AppDetailsType
	for appName, appDetails := range appStateInfo {
		if appName == appToStop {
			doesAppExist = true
			existingApp = appDetails
			break
		}
	}
	if doesAppExist == false {
		return docker.LogAndCreateAppNotFoundError(appToStop)
	} else if !(existingApp.State == docker.Starting || existingApp.State == docker.Available || existingApp.State == docker.Stopping) {
		logger.Warn("only 'Starting' and 'Available' apps can be stopped. State is: %s", existingApp.State.ToString())
		return errors.New("error - stopping app failed")
	} else {
		logger.Debug("App does exist and is now stopped: %s", appToStop)
		return sm.dockerService.StopApp(appToStop)
	}
}

func (sm *appServiceImpl) stopAllApps() error {
	appStateInfo := sm.getAppStateInfo()

	for appName, appDetails := range appStateInfo {
		if appDetails.State == docker.Starting || appDetails.State == docker.Available {
			err := sm.stopApp(appName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
