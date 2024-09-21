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
	} else {
		handler = applyProductionCorsMiddleware(handler)
	}

	logger.Info("Starting server listening on port %s", config.BackendExecutablePort)
	err := http.ListenAndServe(":"+config.BackendExecutablePort, handler)
	if err != nil {
		logger.Fatal("Failed to start server: " + err.Error())
	}
}

func applyProductionCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// TODO The JWT token can have the exact same value as the cookie for simplification.
		// TODO JWT token must be get rid of, before proxying to another app.
		/* TODO requirements
		check if origin is present and allowed ("." + HOST), or maybe a list of allowed urls (ocelot + available apps)
		check if target URL is allowed ("." + HOST)

		if target URL = ocelot-cloud.localhost
			set CORS headers
			if method is OPTIONS
				set CORS headers and return
		isCrossRequest = ...(e.g. nocodb.localhost is accessing ocelot-cloud.localhost)
		if isCrossRequest:
			check if target URL is "ocelot-cloud.localhost", if not, return with error
			if r.path == "/api/auth-check"
				handle request # handler should respond with JWT token
			else
				return with error, "path is not allowed for cross origin requests"
		else:
			check cookie (or JWT token), if valid handle request
		*/

		isOriginAllowed := strings.HasSuffix(origin, "."+config.RootDomain)
		if origin != "" && isOriginAllowed && isAllowedPath(r.URL.Path) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		} else {
			// TODO return and respond with "not allowed"?
		}

		// TODO Should the browser first start with a request using the "OPTIONS" method, or can I directly use POST requests? Not sure if that is allowed, when CORS headers were not set previously.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isAllowedPath(path string) bool {
	// Allow CORS only for the /abc path
	if path == "/abc" {
		return true
	}
	return false
}

func initializeDockerNetwork() {
	// TODO I remember that this is somewhere else used. So duplication? Maybe in ci-runner?
	_ = shared.ExecuteShellCommand("docker network ls | grep -q ocelot-net || docker network create ocelot-net")
}

// TODO When implementing users and groups, here should be a check whether the user is authorized or not to access the app.

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	ocelotDomain := "ocelot-cloud." + config.RootDomain
	// TODO Surprising, why would I need a localDomain? Remove or add an explanation
	localDomain := config.RootDomain + ":" + config.DockerContainerPort
	if r.Host == ocelotDomain || r.Host == localDomain {
		router.ServeHTTP(w, r)
	} else {
		logger.Info("Proxying request to container: %s%s", r.Host, r.URL)
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
