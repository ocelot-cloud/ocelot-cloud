package security

import (
	"github.com/gorilla/mux"
	"net/http"
	"ocelot/backend/tools"
	"strings"
)

var (
	Logger = tools.Logger
	router *mux.Router
)

func InitializeSecurity(routerArg *mux.Router) {
	router = routerArg
	router.HandleFunc("/api/login", loginHandler)
	router.HandleFunc("/api/check-auth", checkAuthHandler)
}

func ApplyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO Add "Origin" header check to prevent CSRF attacks.
		// 1) Scheme must be the same
		// 2) Domain must be the same (example.com) or a subdomain (gitea.example.com)
		// 3) I think port can be ignored since I used the standard ports.
		// TODO In Production mode, when security is enabled, there must be a environment variable called "HOST" (aka Origin) of the form http(s)://*(:[0-9]*), so a URL with http or https, with or without port(?) etc. This is for security to fulfill the origin policy to prevent CSRF attacks.
		// TODO The logic seems weird here, doesn't it? Before the AuthMiddleware there should be a check for the "/api" prefix path, right?
		if strings.HasPrefix(r.URL.Path, "/api/") {
			handleBackendApiRequest(w, r, next)
		} else {
			Logger.Debug("a user requested the frontend resources")
			next.ServeHTTP(w, r)
		}
	})
}

func handleBackendApiRequest(w http.ResponseWriter, r *http.Request, next http.Handler) {
	cookie, err := r.Cookie("auth")
	// TODO Not secure.
	if err != nil || cookie.Value != "valid" {
		Logger.Debug("requests cookie is invalid")
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		Logger.Debug("user has a valid cookie and is allowed to access protected backend functions")
		next.ServeHTTP(w, r)
	}
}
