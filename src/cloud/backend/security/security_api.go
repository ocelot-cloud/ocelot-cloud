package security

import (
	"github.com/gorilla/mux"
	"github.com/ocelot-cloud/shared"
	"net/http"
	"ocelot/backend/config"
	"ocelot/backend/security/internal"
	"strings"
)

var Logger = shared.ProvideLogger()

type SecurityModule struct {
	router *mux.Router
	config *tools.GlobalConfig
}

func ProvideSecurityModule(router *mux.Router, config *tools.GlobalConfig) *SecurityModule {
	router.HandleFunc("/api/login", internal.LoginHandler).Methods("POST")
	return &SecurityModule{router, config}
}

func (s *SecurityModule) ApplyAuthMiddlewares(h http.Handler) http.Handler {
	if s.config.IsSecurityEnabled {
		return s.applyAuthMiddleware(h)
	} else {
		return s.applyCorsPolicy(h)
	}
}

func (s *SecurityModule) applyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO Add "Origin" header check to prevent CSRF attacks.
		// 1) Scheme must be the same
		// 2) Domain must be the same (example.com) or a subdomain (gitea.example.com)
		// 3) I think port can be ignored since I used the standard ports.
		// TODO In Production mode, when security is enabled, there must be a environment variable called "HOST" (aka Origin) of the form http(s)://*(:[0-9]*), so a URL with http or https, with or without port(?) etc. This is for security to fulfill the origin policy to prevent CSRF attacks.
		if strings.HasPrefix(r.URL.Path, "/api/") {
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
		} else {
			Logger.Debug("a user requested the frontend resources")
			next.ServeHTTP(w, r)
		}
	})
}

func (s *SecurityModule) applyCorsPolicy(next http.Handler) http.Handler {
	if s.config.AreCrossOriginRequestsAllowed {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	} else {
		return next
	}
}
