package security

import (
	"github.com/gorilla/mux"
	"net/http"
	"ocelot/backend/apps"
	"ocelot/backend/tools"
	"strings"
)

var (
	Logger = tools.Logger
	router *mux.Router
	config *tools.GlobalConfig
)

func InitializeSecurity(routerArg *mux.Router, configArg *tools.GlobalConfig) {
	router = routerArg
	config = configArg
	router.HandleFunc("/api/login", loginHandler)
	router.HandleFunc("/api/check-auth", checkAuthHandler)
}

// TODO Assert that you can't access any available when you dont have a valid cookie in the request.
func ApplyAuthMiddleware(w http.ResponseWriter, r *http.Request) {
	// TODO Add "Origin" header check to prevent CSRF attacks.

	/* TODO
	The secret is created and stored in the database and send to the ocelot frontend. When a requests provides
	a valid secret, it needs to be deleted from the database afterwards.
	For easy (but not secure) prototype, I can use the cookie as secret.
	The secret should expire after 10 seconds or so.
	Also remove the session cookie from the request when proxying it.
	*/

	// TODO Write a test for the domain check. All tests still pass if it is missing.

	ocelotDomain := "ocelot-cloud." + config.RootDomain // TODO Should be abstracted.
	// TODO Surprising, why would I need a localDomain? Remove or add an explanation
	localDomain := config.RootDomain + ":" + config.BackendPort // TODO Not sure when and why this is needed.

	if r.Host == ocelotDomain || r.Host == localDomain {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			Logger.Trace("accessing ocelot backend")
			applyBackendApiAuthMiddleware(w, r)
		} else {
			Logger.Debug("backend serves frontend resources")
			// TODO serve frontend without auth
			router.ServeHTTP(w, r)
		}
	} else {
		Logger.Debug("app redirect is called")
		// TODO check if header matches regex: "*." + config.RootDomain; if yes continue, else return error.
		apps.ProxyRequestToTheDockerContainer(w, r)
	}
}

func applyBackendApiAuthMiddleware(w http.ResponseWriter, r *http.Request) {
	// TODO Add a test that fails if one of the paths is removed?
	// TODO abstract paths
	// TODO This should be an outer check for the if-block below: if r.Host == "ocelot-cloud."+config.RootDomain
	if r.URL.Path == "/api/login" || r.URL.Path == "/api/check-auth" {
		Logger.Debug("unprotected ocelot-cloud endpoint is addressed: %s", r.URL.Path)
		router.ServeHTTP(w, r)
		return
	}

	// TODO store generated cookie in a repo and check if their value is correct.
	_, err := r.Cookie(tools.CookieName)
	if err != nil {
		Logger.Debug("requests cookie is invalid for request: %s%s", r.Host, r.URL.Path)
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		Logger.Trace("user has a valid cookie and is allowed to access protected backend functions")
		router.ServeHTTP(w, r)
	}
}
