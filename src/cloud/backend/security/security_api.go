package security

import (
	"github.com/gorilla/mux"
	"net/http"
	"ocelot/backend/tools"
	"strings"
)

const cookieName = "auth"

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

func ApplyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO Add "Origin" header check to prevent CSRF attacks.
		if strings.HasPrefix(r.URL.Path, "/api/") {
			applyBackendApiAuthMiddleware(w, r, next)
		} else {
			Logger.Debug("a user requested the frontend resources")
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
	_, err := r.Cookie(cookieName)
	if err != nil {
		Logger.Debug("requests cookie is invalid for request: %s%s", r.Host, r.URL.Path)
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		Logger.Trace("user has a valid cookie and is allowed to access protected backend functions")
		next.ServeHTTP(w, r)
	}
}
