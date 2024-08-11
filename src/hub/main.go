package main

import (
	"github.com/ocelot-cloud/shared"
	"net/http"
	"os"
)

// TODO Allow all log levels. This should be implemented in the shared module. 'shared.setLogLevel("DEBUG")'
func init() {
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		shared.LogLevel = shared.DEBUG
	} else {
		shared.LogLevel = shared.INFO
	}
	Logger = shared.ProvideLogger()
	if shared.LogLevel == shared.DEBUG {
		Logger.Debug("log level set to DEBUG")
	} else {
		Logger.Info("log level set to INFO")
	}
}

func main() {
	initializeDatabase()

	mux := http.NewServeMux()

	mux.HandleFunc(downloadPath, downloadHandler)
	mux.HandleFunc(tagPath, tagHandler)
	mux.HandleFunc(changePasswordPath, changePasswordHandler)

	mux.HandleFunc(appPath, appHandler)
	mux.HandleFunc(searchAppsPath, searchAppsHandler)
	mux.HandleFunc(userPath, userHandler)
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

	handlerWithCors := applyCorsPolicy(mux)

	Logger.Info("Server started on port %s", port)
	err := http.ListenAndServe(":"+port, handlerWithCors)
	if err != nil {
		Logger.Fatal("Server stopped: %v", err)
	}
}

func initializeDatabase() {
	// Strange phenomenon: When I run ./hub via terminal and run tests in separate terminal, everything works
	// as expected. But when I run hub as a daemon process, via bash (&) or ci-runner, the tests fail with
	// this DB error: "attempt to write readonly database". So I use in-memory database for all tests.
	if profile == TEST {
		initializeDatabaseWithSource(":memory:")
		Logger.Warn("initializing database only in-memory")
	} else {
		initializeDatabaseWithSource(databaseFile)
	}
}

// TODO Duplication with cloud code
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
