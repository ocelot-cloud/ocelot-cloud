package apps

import (
	"github.com/gorilla/mux"
	"ocelot/backend/tools"
)

// TODO add router to global config
var (
	logger = tools.Logger // TODO should be private
	config *tools.GlobalConfig
	router *mux.Router
	// TODO to small to be its own file
	// TODO The definition of the stack file dir  depending on the global config should be here I guess.
	appFileDir         string
	stackService       appServiceType
	stackConfigService configServiceType // TODO why is this needed? Should be rather a pointer?
)

func InitializeAppService(routerArg *mux.Router, configArg *tools.GlobalConfig) {
	router = routerArg
	config = configArg

	appFileDir = getStackFileDir(config)
	stackConfigService = provideAppConfigService(appFileDir)
	stackService = getStackService(config, stackConfigService)

	registerSecuredEndpoint("/stacks/read", createReadHandler)
	registerSecuredEndpoint("/stacks/deploy", createDeployHandler)
	registerSecuredEndpoint("/stacks/stop", createStopHandler)
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
