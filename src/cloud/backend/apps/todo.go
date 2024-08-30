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
	Logger             = tools.Logger // TODO should be private
	config             *tools.GlobalConfig
	router             *mux.Router
	stackService       StackService
	stackConfigService StackConfigService // TODO why is this needed? Should be rather a pointer?
)

func InitializeAppService(routerArg *mux.Router, configArg *tools.GlobalConfig) {
	config = configArg
	// TODO Are local variables needed?
	stackFileDir = getStackFileDir(config)
	stackConfigService = provideStackConfigService(stackFileDir)
	stackService = getStackService(config, stackConfigService)

	router = routerArg
	registerSecuredEndpoint("/stacks/read", createReadHandler(stackService))
	registerSecuredEndpoint("/stacks/deploy", createDeployHandler(stackService))
	registerSecuredEndpoint("/stacks/stop", createStopHandler(stackService))
}

func getStackFileDir(config *tools.GlobalConfig) string {
	if config.UseDummyStacks {
		return "stacks/dummy"
	} else {
		return "stacks/local"
	}
}

func getStackService(config *tools.GlobalConfig, stackConfigService StackConfigService) StackService {
	if config.AreMocksEnabled {
		Logger.Debug("Using mock DockerService")
		return ProvideStackServiceMocked(stackConfigService)
	} else {
		Logger.Debug("Using real DockerService")
		return ProvideStackServiceReal(stackConfigService)
	}
}

func registerSecuredEndpoint(path string, handlerFunc http.HandlerFunc) {
	router.Handle("/api"+path, security.ApplyAuthMiddlewares(handlerFunc))
}

func createReadHandler(stackService StackService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Only GET method is supported.", http.StatusMethodNotAllowed)
			return
		}

		stackStateInfo := stackService.GetStackStateInfo()
		response := make([]tools.ResponsePayloadDto, 0)
		for stackName, stackDetails := range stackStateInfo {
			response = append(response, tools.ResponsePayloadDto{stackName, stackDetails.State.String(), stackDetails.Path})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func createDeployHandler(stackService StackService) http.HandlerFunc {
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

		if err = stackService.DeployStack(stackName); err != nil {
			if err != nil {
				Logger.Error("Deploying stack failed: " + stackName + "\n" + err.Error() + "\n")
				http.Error(w, "Deploying stack failed: "+stackName, http.StatusInternalServerError)
			}
			return
		}
	}
}

func createStopHandler(stackService StackService) http.HandlerFunc {
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

		if err = stackService.StopStack(stackName); err != nil {
			if err != nil {
				Logger.Warn("error when trying to stop stack, %s", err.Error())
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
