package setup

import (
	"github.com/gorilla/mux" // TODO To be wrapped?
	"github.com/ocelot-cloud/shared"
	"github.com/ocelot-cloud/shared/utils"
	"net/http"
	"ocelot/backend/apps"
	"ocelot/backend/security"
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
	if config.IsGuiEnabled {
		initializeFrontendResourceDelivery()
	}

	proxy := http.HandlerFunc(proxyHandler)
	handler := security.ApplyAuthMiddleware(proxy)
	if config.AreCrossOriginRequestsAllowed {
		handler = utils.GetCorsDisablingHandler(handler)
	}

	logger.Info("Starting server listening on port %s", config.BackendExecutablePort)
	err := http.ListenAndServe(":"+config.BackendExecutablePort, handler)
	if err != nil {
		logger.Fatal("Failed to start server: " + err.Error())
	}
}

func initializeDockerNetwork() {
	// TODO I remember that this is somewhere else used. So duplication? Maybe in ci-runner?
	_ = shared.ExecuteShellCommand("docker network ls | grep -q ocelot-net || docker network create ocelot-net")
}

// TODO When implementing users and groups, here should be a check whether the user is authorized or not to access the app.

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	ocelotDomain := "ocelot-cloud." + config.RootDomain // TODO Should be abstracted.
	// TODO Surprising, why would I need a localDomain? Remove or add an explanation
	localDomain := config.RootDomain + ":" + config.DockerContainerPort
	if r.Host == ocelotDomain || r.Host == localDomain {
		// TODO Logic is unclear to me, when is this case triggered, and where does it go?
		router.ServeHTTP(w, r)
	} else {
		logger.Info("proxying request to container: %s%s", r.Host, r.URL)
		apps.ProxyRequestToTheDockerContainer(w, r)
	}
}

func initializeFrontendResourceDelivery() {
	router.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Attempt to open the requested file within the ./dist directory.
		_, err := http.Dir("./dist").Open(r.URL.Path)

		// If the requested file does not exist (err is not nil) and the path does not seem to refer to
		// a static file (i.e. no dot extension like ".css"), then serve index.html. This caters to SPA routing needs,
		// allowing frontend routes to be handled by index.html.
		// This means that users can directly access pages with paths such as "example.com/some/path".
		if err != nil && !strings.Contains(r.URL.Path, ".") {
			logger.Debug("Serving index.html for SPA route: %s", r.URL.Path)
			http.ServeFile(w, r, "./dist/index.html")
			return
		}

		// If the request is for a static file or if the file exists, serve it directly.
		// This handles requests for JS, CSS, images, etc.
		logger.Debug("Serving static content at '%s'", r.URL.Path)
		http.FileServer(http.Dir("./dist")).ServeHTTP(w, r)
	}))
}
