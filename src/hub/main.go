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

	handlerWithCors := applyCorsPolicy(mux)
	Logger.Info("Server starting on port %s", port)
	err := http.ListenAndServe(":"+port, handlerWithCors)
	if err != nil {
		Logger.Fatal("Server stopped: %v", err)
	}
}

func initializeHandlers(mux *http.ServeMux) {
	registerUnprotectedHandler(mux, loginPath, loginHandler)
	registerUnprotectedHandler(mux, registrationPath, registrationHandler)
	registerUnprotectedHandler(mux, downloadPath, downloadHandler)
	registerUnprotectedHandler(mux, getTagsPath, getTagsHandler)
	registerUnprotectedHandler(mux, searchAppsPath, searchAppsHandler)

	registerProtectedHandler(mux, authCheckPath, authCheckHandler)
	registerProtectedHandler(mux, tagUploadPath, tagUploadHandler)
	registerProtectedHandler(mux, tagDeletePath, tagDeleteHandler)
	registerProtectedHandler(mux, changePasswordPath, changePasswordHandler)
	registerProtectedHandler(mux, appCreationPath, appCreationHandler)
	registerProtectedHandler(mux, appGetListPath, appGetListHandler)
	registerProtectedHandler(mux, appDeletePath, appDeleteHandler)
	registerProtectedHandler(mux, deleteUserPath, userDeleteHandler)
	registerProtectedHandler(mux, logoutPath, logoutHandler)

	if profile == TEST {
		Logger.Warn("opening unprotected full data wipe endpoint meant for testing only")
		registerUnprotectedHandler(mux, wipeDataPath, wipeDataHandler)

		sampleUser := "sample"
		// TODO Handle error -> logger.Fatal
		repo.CreateUser(&RegistrationForm{sampleUser, "password", "admin@admin.com"})
		Logger.Warn("Created '%s' user with weak password for manual testing", sampleUser)
	}
}

// TODO Handle error of "repo.CreateUser(&RegistrationForm..." -> logger.Fatal

func initializeDatabase() {
	// Strange phenomenon: When I run ./hub via terminal and run tests in separate terminal, everything works
	// as expected. But when I run hub as a daemon process, via bash (&) or ci-runner, the tests fail with
	// this DB error: "attempt to write readonly database". So I use in-memory database for all tests.
	if profile == TEST {
		initializeDatabaseWithSource(":memory:")
		Logger.Warn("initializing database only in-memory - when application stops, all data will be deleted")
	} else {
		initializeDatabaseWithSource(databaseFile)
	}
}

func applyCorsPolicy(next http.Handler) http.Handler {
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

// getUserFromContext Since only authenticated users are added to the context, it only works in protected handlers.
func getUserFromContext(r *http.Request) string {
	return r.Context().Value("user").(string)
}

func registerProtectedHandler(mux *http.ServeMux, path string, handler http.HandlerFunc) {
	mux.Handle(path, authMiddleware(handler))
}

func registerUnprotectedHandler(mux *http.ServeMux, path string, handler http.HandlerFunc) {
	mux.HandleFunc(path, handler)
}
