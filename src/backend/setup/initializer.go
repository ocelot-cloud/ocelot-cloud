package setup

import (
	"github.com/ocelot-cloud/shared"
	"github.com/ocelot-cloud/shared/utils"
	"net/http"
	"ocelot/backend/apps"
	"ocelot/backend/tools"
	"strings"
)

var (
	logger = tools.Logger
)

func InitializeApplication() {
	tools.Router.HandleFunc("/api/login", loginHandler)
	tools.Router.HandleFunc("/api/check-auth", checkAuthHandler)
	apps.InitializeAppService()
	// TODO I need RegisterProtectedRoutes and RegisterUnprotectedRoutes, also aggregated all routes in a single module, so it is immediately clear what is where.
	tools.RegisterRoutes([]tools.Route{
		{"/secret", SecretHandler},
	})

	initializeDockerNetwork()
	if tools.Config.IsGuiEnabled {
		initializeFrontendResourceDelivery()
	}

	var handler http.Handler = http.HandlerFunc(ApplyAuthMiddleware)
	if tools.Config.AreCrossOriginRequestsAllowed {
		handler = utils.GetCorsDisablingHandler(handler)
	}

	logger.Info("Starting server listening on port %s", tools.Config.BackendExecutablePort)
	err := http.ListenAndServe(":"+tools.Config.BackendExecutablePort, handler)
	if err != nil {
		logger.Fatal("Failed to start server: " + err.Error())
	}
}

func initializeDockerNetwork() {
	// TODO I remember that this is somewhere else used. So duplication? Maybe in ci-runner?
	_ = shared.ExecuteShellCommand("docker network ls | grep -q ocelot-net || docker network create ocelot-net")
}

// TODO When implementing users and groups, here should be a check whether the user is authorized or not to access the app.

func initializeFrontendResourceDelivery() {
	tools.Router.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
