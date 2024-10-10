package setup

import (
	"net/http"
	"ocelot/backend/apps"
	"ocelot/backend/tools"
	"strings"
)

var (
	Logger = tools.Logger
)

// TODO Assert that you can't access any available app when you dont have a valid cookie in the request.
func ApplyAuthMiddleware(w http.ResponseWriter, r *http.Request) {
	// TODO Add "Origin" header check to prevent CSRF attacks.
	// TODO	remove the session cookie from the request when proxying it.
	// TODO Networking logic becomes quite complex. I should create a  separate unit and create unit tests.
	// TODO Write a test for the domain check. All tests still pass if it is missing.
	Logger.Debug("Request path: %s", r.URL.Path)
	if isAddressedToOcelotHost(r) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			Logger.Trace("accessing ocelot backend")
			applyBackendApiAuthMiddleware(w, r)
		} else {
			Logger.Debug("backend serves frontend resources")
			router.ServeHTTP(w, r)
		}
	} else {
		Logger.Debug("app redirect is called")
		// TODO check if header matches regex: "*." + config.RootDomain; if yes continue, else return error.
		apps.ProxyRequestToTheDockerContainer(w, r)
	}
}

// TODO Can some stuff be removed?
func isAddressedToOcelotHost(r *http.Request) bool {
	Logger.Debug("Checking host of request: %s", r.Host)
	if tools.Config.Profile == tools.PROD {
		ocelotHost := "ocelot-cloud." + tools.Config.RootDomain // TODO Should maybe be abstracted?
		Logger.Debug("ocelot host: %s", ocelotHost)
		Logger.Debug("request host: %s", r.Host)
		return r.Host == ocelotHost
	} else {
		host := tools.Config.RootDomain + ":" + tools.Config.PubliclyAvailablePort
		Logger.Debug("host: %s", host)
		return r.Host == host
	}
}

func applyBackendApiAuthMiddleware(w http.ResponseWriter, r *http.Request) {
	// TODO Add a test that fails if one of the paths is removed?
	// TODO abstract paths
	// TODO This should be an outer check for the if-block below: if r.Host == "ocelot-cloud."+config.RootDomain
	if r.URL.Path == "/api/login" || r.URL.Path == "/api/check-auth" || r.URL.Path == "/api/stacks/wipe-data" {
		Logger.Debug("unprotected ocelot-cloud endpoint is addressed: %s", r.URL.Path)
		router.ServeHTTP(w, r)
		return
	}

	r, err := GetRequestWithAuthContext(w, r)
	if err != nil {
		return
	}

	Logger.Trace("user has a valid cookie and is allowed to access protected backend functions")
	router.ServeHTTP(w, r)
}
