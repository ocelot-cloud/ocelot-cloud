package main

import (
	"context"
	"github.com/ocelot-cloud/shared"
	"net/http"
	"os"
)

func init() {
	Logger = shared.ProvideLogger(os.Getenv("LOG_LEVEL"))
	Logger.Info("log level set to: %s", shared.GetLogLevel())
}

func main() {
	initializeDatabase()
	mux := http.NewServeMux()
	initializeHandlers(mux)

	Logger.Info("Server starting on port %s", port)
	err := http.ListenAndServe(":"+port, getCorsDisablingHandler(mux))
	if err != nil {
		Logger.Fatal("Server stopped: %v", err)
	}
}

func initializeDatabase() {
	if profile == TEST {
		initializeDatabaseWithSource(":memory:")
		Logger.Warn("initializing database only in-memory - when application stops, all data will be deleted")
	} else {
		initializeDatabaseWithSource(databaseFile)
	}
}

type route struct {
	path    string
	handler http.HandlerFunc
}

func initializeHandlers(mux *http.ServeMux) {
	unprotectedRoutes := []route{
		{loginPath, loginHandler},
		{registrationPath, registrationHandler},
		{downloadPath, downloadHandler},
		{getTagsPath, getTagsHandler},
		{searchAppsPath, searchAppsHandler},
	}

	protectedRoutes := []route{
		{authCheckPath, authCheckHandler},
		{tagUploadPath, tagUploadHandler},
		{tagDeletePath, tagDeleteHandler},
		{changePasswordPath, changePasswordHandler},
		{appCreationPath, appCreationHandler},
		{appGetListPath, appGetListHandler},
		{appDeletePath, appDeleteHandler},
		{deleteUserPath, userDeleteHandler},
		{logoutPath, logoutHandler},
	}

	if profile == TEST {
		Logger.Warn("opening unprotected full data wipe endpoint meant for testing only")
		unprotectedRoutes = append(unprotectedRoutes, route{wipeDataPath, wipeDataHandler})

		sampleUser := "sample"
		err := repo.CreateUser(&RegistrationForm{sampleUser, "password", "admin@admin.com"})
		if err != nil {
			Logger.Fatal("Failed to create '%s' user: %v.", sampleUser, err)
		}
		Logger.Warn("Created '%s' user with weak password for manual testing", sampleUser)
	}

	registerUnprotectedRoutes(mux, unprotectedRoutes)
	registerProtectedRoutes(mux, protectedRoutes)
}

// getUserFromContext Since only authenticated users are added to the context, it only works in protected handlers.
func getUserFromContext(r *http.Request) string {
	return r.Context().Value("user").(string)
}

func registerUnprotectedRoutes(mux *http.ServeMux, routes []route) {
	for _, r := range routes {
		mux.HandleFunc(r.path, r.handler)
	}
}

func registerProtectedRoutes(mux *http.ServeMux, routes []route) {
	for _, r := range routes {
		mux.Handle(r.path, authMiddleware(r.handler))
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := checkAuthentication(w, r)
		if err != nil {
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// getCorsDisablingHandler This is necessary to allow cross-origin requests from the ocelot-cloud GUI to the hub.
// The "Origin" header is managed and checked with custom logic to prevent CSRF attacks.
func getCorsDisablingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
