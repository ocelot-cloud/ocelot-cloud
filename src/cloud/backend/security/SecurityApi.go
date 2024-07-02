package security

import (
	"github.com/gorilla/mux"
	"github.com/ocelot-cloud/shared"
	"net/http"
	"ocelot/backend/config"
	"strings"
)

var Logger = shared.ProvideLogger()

type SecurityModule struct {
	router *mux.Router
	config *tools.GlobalConfig
}

func ProvideSecurityModule(router *mux.Router, config *tools.GlobalConfig) *SecurityModule {
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
		if strings.HasPrefix(r.URL.Path, "/api/") {
			cookie, err := r.Cookie("auth")
			// TODO Not secure.
			if err != nil || cookie.Value != "valid" {
				Logger.Debug("requests cookie is invalid")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
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
