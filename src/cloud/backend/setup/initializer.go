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

// TODO Write test for all that validation logic
func applyProductionCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO The JWT token can have the exact same value as the cookie for simplification.
		// TODO When proxying to an app, the JWT token must be removed.

		/* TODO Implement
		origin := r.Header.Get("Origin")
		host := r.Host

		if !strings.HasSuffix(origin, ".localhost") { // TODO use "." + config.RootDomain or so
			w.Write([]byte("invalid origin suffix")) // TODO more detail
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !strings.HasSuffix(host, ".localhost") { // TODO abstract the suffix
			w.Write([]byte("invalid host suffix")) // TODO more detail
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		*/

		// The "Host" is the destination domain of the request. The "Origin" is the URL displayed in the browser.
		// For common requests, both are equal. For cross requests, the two are different.
		// For example, if a browser visits "nocodb.localhost", that is the origin. If the same page makes
		// a cross-request to "ocelot-cloud.localhost", then that is the host.

		/* TODO Implement
		if host == "ocelot-cloud.localhost" { // TODO abstract
			w.Header().Set("Access-Control-Allow-Origin", "*") // TODO check that subsequent logic handles not allowed requests
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		isCrossRequest := host != origin
		if isCrossRequest {
			if host != "ocelot-cloud.localhost" { // TODO abstract
				w.Write([]byte("cross requests are only allowed to target domain 'ocelot-cloud.localhost', but you tried to access: " + r.URL.Host)) // TODO more detail
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else if r.URL.Path == "/api/auth-check" { // TODO abstract
				next.ServeHTTP(w, r) // TODO handler should respond with JWT token
			} else {
				w.Write([]byte("path is not allowed for cross origin requests")) // TODO more detail
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		} else {
			// TODO check cookie (or JWT token), if valid handle request
			next.ServeHTTP(w, r)
		}

		*/
		next.ServeHTTP(w, r)
	})
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
