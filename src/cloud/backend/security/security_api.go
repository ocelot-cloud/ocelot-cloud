package security

import (
	"github.com/gorilla/mux"
	"github.com/ocelot-cloud/shared/utils"
	"net/http"
	"ocelot/backend/tools"
	"strings"
)

const CookieName = "ocelot-auth"

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
func ApplyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO Add "Origin" header check to prevent CSRF attacks.

		/* TODO
		The secret is created and stored in the database and send to the ocelot frontend. When a requests provides
		a valid secret, it needs to be deleted from the database afterwards.
		For easy (but not secure) prototype, I can use the cookie as secret.
		The secret should expire after 10 seconds or so.
		Also remove the session cookie from the request when proxying it.
		*/

		// TODO Write a test for the domain check. All tests still pass if it is missing.
		if r.Header.Get(utils.OriginHeader) == "ocelot-cloud."+config.RootDomain && strings.HasPrefix(r.URL.Path, "/api/") {
			Logger.Trace("accessing ocelot backend")
			applyBackendApiAuthMiddleware(w, r, next)
		} else {
			// TODO check if header fits: "*." + config.RootDomain; else return error.
			next.ServeHTTP(w, r)
		}
	})
}

type Route struct {
	Path        string
	HandlerFunc http.HandlerFunc
}

func RegisterRoutes(routes []Route) {
	for _, r := range routes {
		router.Handle("/api"+r.Path, r.HandlerFunc)
	}
}

func applyBackendApiAuthMiddleware(w http.ResponseWriter, r *http.Request, next http.Handler) {
	// TODO Add a test that fails if one of the paths is removed?
	// TODO abstract paths
	// TODO This should be an outer check for the if-block below: if r.Host == "ocelot-cloud."+config.RootDomain
	if r.URL.Path == "/api/login" || r.URL.Path == "/api/check-auth" {
		Logger.Debug("unprotected ocelot-cloud endpoint is addressed: %s", r.URL.Path)
		next.ServeHTTP(w, r)
		return
	}

	// TODO store generated cookie in a repo and check if their value is correct.
	_, err := r.Cookie(CookieName)
	if err != nil {
		Logger.Debug("requests cookie is invalid for request: %s%s", r.Host, r.URL.Path)
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		Logger.Trace("user has a valid cookie and is allowed to access protected backend functions")
		next.ServeHTTP(w, r)
	}
}
