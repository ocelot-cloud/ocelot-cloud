package apps

import (
	"github.com/gorilla/mux"
	"net/http"
	"ocelot/backend/security"
	"ocelot/backend/tools"
)

// TODO add router to global config
var (
	logger = tools.Logger
	router *mux.Router
	config *tools.GlobalConfig
	// TODO The definition of the stack file dir  depending on the global config should be here I guess.
	appFileDir         string
	appService         appServiceType
	stackConfigService configServiceType // TODO why is this needed? Should be rather a pointer?
)

func InitializeAppService(routerArg *mux.Router, configArg *tools.GlobalConfig) {
	config = configArg
	router = routerArg

	appFileDir = getStackFileDir(config)
	stackConfigService = provideAppConfigService(appFileDir)
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
		return dummyAppAssetsDir
	} else {
		return realAppAssetsDir
	}
}

func getStackService(config *tools.GlobalConfig, stackConfigService configServiceType) appServiceType {
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
