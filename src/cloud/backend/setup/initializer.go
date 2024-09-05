package setup

import (
	"fmt"
	"github.com/gorilla/mux" // TODO To be wrapped?
	"github.com/ocelot-cloud/shared"
	"github.com/ocelot-cloud/shared/utils"
	"net/http"
	"ocelot/backend/apps"
	"ocelot/backend/tools"
	"strings"
)

var (
	logger = tools.Logger
	router *mux.Router
	config *tools.GlobalConfig
)

func InitializeApplication(routerArg *mux.Router, configArg *tools.GlobalConfig) {
	router = routerArg
	config = configArg

	apps.InitializeAppService(router, config)
	initializeDockerNetwork()
	initializeFunctionalEndpoints()

	proxyHandler := buildProxyHandler()
	logger.Info("Starting server listening on port %s", config.BackendExecutablePort)
	// TODO utils.GetCorsDisablingHandler should only be enabled in TEST profile
	err := http.ListenAndServe(":"+config.BackendExecutablePort, utils.GetCorsDisablingHandler(http.HandlerFunc(proxyHandler)))
	if err != nil {
		logger.Fatal("Failed to start server: " + err.Error())
	}
}

func initializeDockerNetwork() {
	// TODO I remember that this is somewhere else used. So duplication? Maybe in ci-runner?
	_ = shared.ExecuteShellCommand("docker network ls | grep -q ocelot-net || docker network create ocelot-net")
}

func buildProxyHandler() func(w http.ResponseWriter, r *http.Request) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		ocelotDomain := "ocelot-cloud." + config.RootDomain
		// TODO Surprising, why would I need a localDomain? Remove or add an explanation
		localDomain := config.RootDomain + ":" + config.DockerContainerPort
		if r.Host == ocelotDomain || r.Host == localDomain {
			router.ServeHTTP(w, r)
		} else {
			apps.ProxyRequestToTheDockerContainer(w, r)
		}
	}
	return handler
}

func initializeFunctionalEndpoints() {
	// TODO Is that still necessary?
	router.HandleFunc("/api/hello", helloHandler)

	if config.IsGuiEnabled {
		initializeFrontendResourceDelivery()
	}
}

func initializeFrontendResourceDelivery() {
	// TODO utils.GetCorsDisablingHandler should be used only once.
	router.PathPrefix("/").Handler(utils.GetCorsDisablingHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Attempt to open the requested file within the ./dist directory.
		_, err := http.Dir("./dist").Open(r.URL.Path)

		// If the requested file does not exist (err is not nil) and the path does not seem to
		// refer to a static file (no extension), then serve index.html. This caters to SPA routing needs,
		// allowing frontend routes to be handled by index.html.
		// This means that users can directly access pages with paths such as "example.com/some/path".
		if err != nil && !strings.Contains(r.URL.Path, ".") {
			logger.Debug("Serving index.html for SPA route ''", r.URL.Path)
			http.ServeFile(w, r, "./dist/index.html")
			return
		}

		// If the request is for a static file or if the file exists, serve it directly.
		// This handles requests for JS, CSS, images, etc.
		logger.Debug("Serving static content at '%s'", r.URL.Path)
		http.FileServer(http.Dir("./dist")).ServeHTTP(w, r)
	})))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<html><body>Hello</body></html>")
}
