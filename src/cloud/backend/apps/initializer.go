package apps

import (
	"github.com/gorilla/mux"
	"net/http"
	"ocelot/backend/apps/global_config"
	_ "ocelot/backend/apps/global_config"
	"ocelot/backend/apps/yaml_config"
	"ocelot/backend/security"
	"ocelot/backend/tools"
)

// TODO add router to global config
var (
	logger             = tools.Logger
	router             *mux.Router
	config             *tools.GlobalConfig
	appService         appServiceType
	stackConfigService yaml_config.ConfigServiceType // TODO why is this needed? Should be rather a pointer?
)

func InitializeAppService(routerArg *mux.Router, configArg *tools.GlobalConfig) {
	config = configArg
	router = routerArg

	global_config.AppFileDir = getStackFileDir(config)
	stackConfigService = yaml_config.ProvideAppConfigService()
	appService = getStackService(config, stackConfigService)

	routes := []security.Route{
		{"/stacks/read", appReadHandler},
		{"/stacks/deploy", appDeployHandler},
		{"/stacks/stop", appStopHandler},
	}
	if config.OpenDataWipeEndpoint {
		router.HandleFunc("/api/stacks/wipe-data", wipeDataHandler)
	}
	security.RegisterRoutes(routes)
}

func getStackFileDir(config *tools.GlobalConfig) string {
	if config.UseDummyStacks {
		return global_config.DummyAppAssetsDir
	} else {
		return global_config.RealAppAssetsDir
	}
}

func getStackService(config *tools.GlobalConfig, stackConfigService yaml_config.ConfigServiceType) appServiceType {
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
		if appDetails.State != Uninitialized {
			err := appService.stopApp(appName)
			if err != nil {
				logger.Error("Couldn't stop app: %s, error: %v", appName, err)
			}
		}
	}
}
