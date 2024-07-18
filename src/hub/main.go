package main

import (
	"github.com/ocelot-cloud/shared"
	"net/http"
)

func init() {
	Logger = shared.ProvideLogger()
}

func main() {
	initializeDatabase()

	http.HandleFunc(downloadPath, downloadHandler)
	http.HandleFunc(tagPath, tagHandler)
	http.HandleFunc(changePasswordPath, changePasswordHandler)
	http.HandleFunc(changeOriginPath, changeOriginHandler)

	http.HandleFunc(appPath, appHandler)
	http.HandleFunc(userPath, userHandler)
	http.HandleFunc(loginPath, loginHandler)

	http.HandleFunc(registrationPath, registrationHandler)

	if profile == TEST {
		Logger.Warn("opening unprotected full data wipe endpoint meant for testing only")
		http.HandleFunc(wipeDataPath, wipeDataHandler)
	}

	Logger.Info("Server started on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
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
