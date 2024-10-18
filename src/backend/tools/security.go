package tools

import (
	"fmt"
	"net/http"
)

// TODO Not sure, but this should maybe be put to "tools" for later reuse? Maybe also put the router and global config there for simplification.
type Route struct {
	Path        string
	HandlerFunc http.HandlerFunc
}

func RegisterRoutes(routes []Route) {
	for _, r := range routes {
		Router.Handle("/api"+r.Path, r.HandlerFunc)
	}
}

type Authorization struct {
	User    string
	IsAdmin bool
}

func GetAuthFromContext(w http.ResponseWriter, r *http.Request) (*Authorization, error) {
	auth, ok := r.Context().Value("auth").(*Authorization)
	if !ok {
		Logger.Error("auth not found in context, but protected handlers must always have an auth context")
		w.WriteHeader(http.StatusInternalServerError)
		return nil, fmt.Errorf("")
	}
	return auth, nil
}
