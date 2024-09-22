package apps

import (
	"github.com/gorilla/mux"
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
	router             *mux.Router
	config             *tools.GlobalConfig
	appService         appServiceType
	stackConfigService yaml.ConfigServiceType // TODO why is this needed? Should be rather a pointer?
)

func InitializeAppService(routerArg *mux.Router, configArg *tools.GlobalConfig) {
	config = configArg
	router = routerArg

	vars.AppFileDir = getStackFileDir(config)
	stackConfigService = yaml.ProvideAppConfigService()
	appService = getStackService(config, stackConfigService)

	routes := []Route{
		{"/stacks/read", appReadHandler},
		{"/stacks/deploy", appDeployHandler},
		{"/stacks/stop", appStopHandler},
	}
	if config.OpenDataWipeEndpoint {
		router.HandleFunc("/api/stacks/wipe-data", wipeDataHandler)
	}
	RegisterRoutes(routes)
}

// TODO Not sure, but this should maybe be put to "tools" for later reuse? Maybe also put the router and global config there for simplification.
type Route struct {
	Path        string
	HandlerFunc http.HandlerFunc
}

func RegisterRoutes(routes []Route) {
	for _, r := range routes {
		router.Handle("/api"+r.Path, r.HandlerFunc)
	}
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
