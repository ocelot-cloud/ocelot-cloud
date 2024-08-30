package setup

import (
	"github.com/gorilla/mux" // TODO To be wrapped?
	"github.com/ocelot-cloud/shared"
	"net/http"
	"net/http/httputil"
	"net/url"
	"ocelot/backend/apps"
	"ocelot/backend/security"
	"ocelot/backend/tools"
	"strings"
)

var (
	securityModule     *security.SecurityModule
	router             *mux.Router
	stackService       apps.StackService
	config             *tools.GlobalConfig
	stackConfigService apps.StackConfigService
)

func InitializeApplication(routerArg *mux.Router, configArg *tools.GlobalConfig, securityModuleArg *security.SecurityModule) {
	router = routerArg
	config = configArg
	securityModule = securityModuleArg

	apps.StackFileDir = getStackFileDir()
	stackConfigService = apps.ProvideStackConfigService(apps.StackFileDir)
	stackService = getStackService(stackConfigService)
	initializeDockerNetwork()
	initializeHandlers()
}

func getStackFileDir() string {
	if config.UseDummyStacks {
		return "stacks/dummy"
	} else {
		return "stacks/local"
	}
}

func getStackService(stackConfigService apps.StackConfigService) apps.StackService {
	if config.AreMocksEnabled {
		apps.Logger.Debug("Using mock DockerService")
		return apps.ProvideStackServiceMocked(stackConfigService)
	} else {
		apps.Logger.Debug("Using real DockerService")
		return apps.ProvideStackServiceReal(stackConfigService)
	}
}

func initializeDockerNetwork() {
	// TODO I remember that this is somewhere else used. So duplication? Maybe in ci-runner?
	_ = shared.ExecuteShellCommand("docker network ls | grep -q ocelot-net || docker network create ocelot-net")
}

func initializeHandlers() {
	initializeFunctionalEndpoints()
	proxyHandler := buildProxyHandler()
	apps.Logger.Info("Starting server listening on port ", config.BackendExecutablePort)
	err := http.ListenAndServe(":"+config.BackendExecutablePort, http.HandlerFunc(proxyHandler))
	if err != nil {
		apps.Logger.Fatal("Failed to start server: " + err.Error())
	}
}

func buildProxyHandler() func(w http.ResponseWriter, r *http.Request) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		ocelotDomain := "ocelot-cloud." + config.RootDomain
		// TODO Surprising, why would I need a localDomain? Remove or add an explanation
		localDomain := config.RootDomain + ":" + config.DockerContainerPort
		if r.Host == ocelotDomain || r.Host == localDomain {
			router.ServeHTTP(w, r)
		} else {
			proxyRequestToTheDockerContainer(w, r)
		}
	}
	return handler
}

// TODO Make sure to remove the ocelot cookie before proxying a request to the service behind, so that it can't read/steal it.
func proxyRequestToTheDockerContainer(w http.ResponseWriter, r *http.Request) {
	apps.Logger.Trace("Proxying request with target host %s", r.Host)
	targetContainer := strings.TrimSuffix(r.Host, "."+config.RootDomain)
	targetPort := stackConfigService.GetStackConfig(targetContainer).Port
	targetURL, err := url.Parse("http://" + targetContainer + ":" + targetPort)
	if err != nil {
		apps.Logger.Error("error when parsing URL, %s", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// the path of original request is preserved
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	r.URL.Host = targetContainer
	r.URL.Scheme = "http"
	r.Header.Set("X-Forwarded-Host", r.Host)
	proxy.ServeHTTP(w, r)
}

func initializeFunctionalEndpoints() {
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/check-session", checkSessionHandler).Methods("GET")
	api.HandleFunc("/hello", helloHandler)

	registerSecuredEndpoint("/stacks/read", createReadHandler(stackService))
	registerSecuredEndpoint("/stacks/deploy", createDeployHandler(stackService))
	registerSecuredEndpoint("/stacks/stop", createStopHandler(stackService))

	if config.IsGuiEnabled {
		initializeFrontendResourceDelivery()
	}
}

func initializeFrontendResourceDelivery() {
	router.PathPrefix("/").Handler(securityModule.ApplyAuthMiddlewares(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Attempt to open the requested file within the ./dist directory.
		_, err := http.Dir("./dist").Open(r.URL.Path)

		// If the requested file does not exist (err is not nil) and the path does not seem to
		// refer to a static file (no extension), then serve index.html. This caters to SPA routing needs,
		// allowing frontend routes to be handled by index.html.
		// This means that users can directly access pages with paths such as "example.com/some/path".
		if err != nil && !strings.Contains(r.URL.Path, ".") {
			apps.Logger.Debug("Serving index.html for SPA route ''", r.URL.Path)
			http.ServeFile(w, r, "./dist/index.html")
			return
		}

		// If the request is for a static file or if the file exists, serve it directly.
		// This handles requests for JS, CSS, images, etc.
		apps.Logger.Debug("Serving static content at '%s'", r.URL.Path)
		http.FileServer(http.Dir("./dist")).ServeHTTP(w, r)
	})))
}

func registerSecuredEndpoint(path string, handlerFunc http.HandlerFunc) {
	router.Handle("/api"+path, securityModule.ApplyAuthMiddlewares(handlerFunc))
}
