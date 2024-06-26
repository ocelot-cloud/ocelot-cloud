package internal

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httputil"
	"net/url"
	"ocelot/security"
	"ocelot/tools"
	"strings"
)

type ApplicationInitializer struct {
	securityModule     *security.SecurityModule
	router             *mux.Router
	stackService       StackService
	config             *tools.GlobalConfig
	stackConfigService StackConfigService
}

func ProvideAppInitializer(router *mux.Router, config *tools.GlobalConfig, securityModule *security.SecurityModule) ApplicationInitializer {
	return ApplicationInitializer{securityModule, router, nil, config, nil}
}

func (a *ApplicationInitializer) InitializeApplicationInternally() {
	StackFileDir = a.getStackFileDir()
	a.stackConfigService = ProvideStackConfigService(StackFileDir)
	a.stackService = a.getStackService(a.stackConfigService)
	a.initializeDockerNetwork()
	a.initializeHandlers()
}

func (a *ApplicationInitializer) getStackFileDir() string {
	if a.config.UseDummyStacks {
		return "stacks/dummy"
	} else {
		return "stacks/local"
	}
}

func (a *ApplicationInitializer) getStackService(stackConfigService StackConfigService) StackService {
	if a.config.AreMocksEnabled {
		Logger.Debug("Using mock DockerService")
		return ProvideStackServiceMocked(stackConfigService)
	} else {
		Logger.Debug("Using real DockerService")
		return ProvideStackServiceReal(stackConfigService)
	}
}

func (a *ApplicationInitializer) initializeDockerNetwork() {
	_ = tools.ExecuteShellCommand("docker network ls | grep -q ocelot-net || docker network create ocelot-net")
}

func (a *ApplicationInitializer) initializeHandlers() {
	a.initializeFunctionalEndpoints()
	proxyHandler := a.buildProxyHandler()
	Logger.Info("Starting server listening on port 8080")
	err := http.ListenAndServe(":8080", http.HandlerFunc(proxyHandler))
	if err != nil {
		Logger.Fatal("Failed to start server: " + err.Error())
	}
}

func (a *ApplicationInitializer) buildProxyHandler() func(w http.ResponseWriter, r *http.Request) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Host == "ocelot-cloud.localhost" || r.Host == "localhost:8080" {
			a.router.ServeHTTP(w, r)
		} else {
			a.proxyRequestToTheDockerContainer(w, r)
		}
	}
	return handler
}

func (a *ApplicationInitializer) proxyRequestToTheDockerContainer(w http.ResponseWriter, r *http.Request) {
	Logger.Trace("Proxying request with target host %s", r.Host)
	targetContainer := strings.TrimSuffix(r.Host, ".localhost")
	targetPort := a.stackConfigService.GetStackConfig(targetContainer).Port
	targetURL, err := url.Parse("http://" + targetContainer + ":" + targetPort)
	if err != nil {
		Logger.Error("error when parsing URL, %s", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	r.URL.Host = targetContainer
	r.URL.Scheme = "http"
	r.Header.Set("X-Forwarded-Host", r.Host)
	proxy.ServeHTTP(w, r)
}

func (a *ApplicationInitializer) initializeFunctionalEndpoints() {
	api := a.router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/login", login).Methods("POST")
	api.HandleFunc("/check-session", checkSessionHandler).Methods("GET")
	api.HandleFunc("/hello", a.helloHandler)

	a.registerSecuredEndpoint("/stacks/read", createReadHandler(a.stackService))
	a.registerSecuredEndpoint("/stacks/deploy", createDeployHandler(a.stackService))
	a.registerSecuredEndpoint("/stacks/stop", createStopHandler(a.stackService))

	if a.config.IsGuiEnabled {
		a.InitializeFrontendResourceDelivery()
	}
}

func (a *ApplicationInitializer) InitializeFrontendResourceDelivery() {
	a.router.PathPrefix("/").Handler(a.securityModule.ApplyAuthMiddlewares(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Attempt to open the requested file within the ./dist directory.
		_, err := http.Dir("./dist").Open(r.URL.Path)

		// If the requested file does not exist (err is not nil) and the path does not seem to
		// refer to a static file (no extension), serve index.html. This caters to SPA routing needs,
		// allowing frontend routes to be handled by index.html.
		// This means that users can directly access pages with paths such as "example.com/some/path".
		if err != nil && !strings.Contains(r.URL.Path, ".") {
			Logger.Debug("Serving index.html for SPA route ''", r.URL.Path)
			http.ServeFile(w, r, "./dist/index.html")
			return
		}

		// If the request is for a static file or if the file exists, serve it directly.
		// This handles requests for JS, CSS, images, etc.
		Logger.Debug("Serving static content at '%s'", r.URL.Path)
		http.FileServer(http.Dir("./dist")).ServeHTTP(w, r)
	})))
}

func (a *ApplicationInitializer) registerSecuredEndpoint(path string, handlerFunc http.HandlerFunc) {
	a.router.Handle("/api"+path, a.securityModule.ApplyAuthMiddlewares(handlerFunc))
}
