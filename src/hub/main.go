package main

import (
	"github.com/ocelot-cloud/shared"
	"net/http"
	"os"
)

func init() {
	shared.SetLogLevel(os.Getenv("LOG_LEVEL"))
	Logger = shared.ProvideLogger()
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
	mux.HandleFunc(downloadPath, downloadHandler)
	mux.HandleFunc(tagUploadPath, tagUploadHandler)
	mux.HandleFunc(tagDeletePath, tagDeleteHandler)
	mux.HandleFunc(getTagsPath, getTagsHandler)
	mux.HandleFunc(changePasswordPath, changePasswordHandler)
	mux.HandleFunc(appCreationPath, appHandler)
	mux.HandleFunc(appGetListPath, appGetListHandler)
	mux.HandleFunc(appDeletePath, appDeleteHandler)
	mux.HandleFunc(searchAppsPath, searchAppsHandler)
	mux.HandleFunc(deleteUserPath, userDeleteHandler)
	mux.HandleFunc(logoutPath, logoutHandler)
	mux.HandleFunc(loginPath, loginHandler)
	mux.HandleFunc(authCheckPath, authCheckHandler)
	mux.HandleFunc(registrationPath, registrationHandler)

	if profile == TEST {
		Logger.Warn("opening unprotected full data wipe endpoint meant for testing only")
		mux.HandleFunc(wipeDataPath, wipeDataHandler)

		sampleUser := "sample"
		repo.CreateUser(&RegistrationForm{sampleUser, "password", "admin@admin.com"})
		Logger.Warn("Created '%s' user with weak password for manual testing", sampleUser)
	}
}

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
