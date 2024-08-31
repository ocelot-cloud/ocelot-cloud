package apps

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"ocelot/backend/security"
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

	registerSecuredEndpoint("/stacks/read", createReadHandler(stackService))
	registerSecuredEndpoint("/stacks/deploy", createDeployHandler(stackService))
	registerSecuredEndpoint("/stacks/stop", createStopHandler(stackService))
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

func registerSecuredEndpoint(path string, handlerFunc http.HandlerFunc) {
	router.Handle("/api"+path, security.ApplyAuthMiddleware(handlerFunc))
}

func createReadHandler(stackService appServiceType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Only GET method is supported.", http.StatusMethodNotAllowed)
			return
		}

		stackStateInfo := stackService.getAppStateInfo()
		response := make([]tools.ResponsePayloadDto, 0)
		for stackName, stackDetails := range stackStateInfo {
			response = append(response, tools.ResponsePayloadDto{stackName, stackDetails.State.toString(), stackDetails.Path})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func createDeployHandler(stackService appServiceType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Only POST method is supported.", http.StatusMethodNotAllowed)
			return
		}

		stackName, err := decodeStackInfo(r)
		if err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		if err = stackService.deployApp(stackName); err != nil {
			if err != nil {
				logger.Error("Deploying stack failed: " + stackName + "\n" + err.Error() + "\n")
				http.Error(w, "Deploying stack failed: "+stackName, http.StatusInternalServerError)
			}
			return
		}
	}
}

func createStopHandler(stackService appServiceType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Only POST method is supported.", http.StatusMethodNotAllowed)
			return
		}

		stackName, err := decodeStackInfo(r)
		if err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		if err = stackService.stopApp(stackName); err != nil {
			if err != nil {
				logger.Warn("error when trying to stop stack, %s", err.Error())
				http.Error(w, "Stopping stack failed: "+stackName, http.StatusInternalServerError)
			}
			return
		}
	}
}

func decodeStackInfo(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	var stackInfo tools.StackInfo
	err = json.Unmarshal(body, &stackInfo)
	if err != nil {
		return "", err
	}
	return stackInfo.Name, nil
}
