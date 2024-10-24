package apps

import (
	"net/http"
	"ocelot/backend/apps/docker"
	"ocelot/backend/apps/vars"
	_ "ocelot/backend/apps/vars"
	"ocelot/backend/apps/yaml"
	"ocelot/backend/tools"
)

// TODO add router to global config
var (
	logger             = tools.Logger
	appService         appServiceType
	stackConfigService yaml.ConfigServiceType // TODO why is this needed? Should be rather a pointer?
)

func InitializeAppService() {
	vars.AppFileDir = getStackFileDir(tools.Config)
	stackConfigService = yaml.ProvideAppConfigService()
	appService = getStackService(tools.Config, stackConfigService)

	routes := []tools.Route{
		{"/stacks/read", appReadHandler},
		{"/stacks/deploy", appDeployHandler},
		{"/stacks/stop", appStopHandler},
	}
	if tools.Config.OpenDataWipeEndpoint {
		tools.Router.HandleFunc("/api/stacks/wipe-data", wipeDataHandler)
	}
	tools.RegisterRoutes(routes)
}

func getStackFileDir(config *tools.GlobalConfig) string {
	if config.UseDummyStacks {
		return vars.DummyAppAssetsDir
	} else {
		return vars.RealAppAssetsDir
	}
}

func getStackService(config *tools.GlobalConfig, stackConfigService yaml.ConfigServiceType) appServiceType {
	if config.AreMocksEnabled {
		logger.Debug("Using mock DockerService")
		return provideAppServiceMocked(stackConfigService)
	} else {
		logger.Debug("Using real DockerService")
		return provideAppServiceReal(stackConfigService)
	}
}

func wipeDataHandler(w http.ResponseWriter, r *http.Request) {
	apps := appService.getAppStateInfo()
	for appName, appDetails := range apps {
		if appDetails.State != docker.Uninitialized {
			err := appService.stopApp(appName)
			if err != nil {
				logger.Error("Couldn't stop app: %s, error: %v", appName, err)
			}
		}
	}
}
