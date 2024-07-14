package main

import (
	"github.com/ocelot-cloud/shared"
	"net/http"
	"os"
)

func init() {
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
			Logger.Error("Error creating data directory: %v. Terminating application.", err)
			os.Exit(1)
		}
	}
	Logger = shared.ProvideLogger() // TODO dataDir should be moved to "shared". ProvideLogger should create the logs.txt in dataDir
}

// TODO Should be put in shared folder. Also necessary for logger files.

func main() {
	initializeDatabase()

	// TODO Maybe wrap gorilla/mux like in backend, apply a common security policy and put it in shared module.
	// TODO apply middleware?
	http.HandleFunc(downloadPath, downloadHandler)
	http.HandleFunc(tagPath, tagHandler)
	http.HandleFunc(changePasswordPath, changePasswordHandler)
	http.HandleFunc(changeOriginPath, changeOriginHandler)

	http.HandleFunc(appPath, appHandler)
	http.HandleFunc(userPath, userHandler)
	http.HandleFunc(loginPath, loginHandler)

	// Registration process is excluded from security, so no middleware is used.
	http.HandleFunc(registrationPath, registrationHandler)

	if profile == TEST {
		Logger.Warn("as") //TODO
		http.HandleFunc(wipeDataPath, wipeDataHandler)
	}

	Logger.Info("Server started on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		// TODO Is server stop sometimes normal, e.g. when gracefully shutdown?
		Logger.Fatal("Server stopped: %v", err)
	}
}

func initializeDatabase() {
	// Strange phenomenon: When I run ./hub via terminal and run tests in separate terminal, everything works
	// as expected. But when I run hub as a daemon process, via bash or ci-runner, the tests fail with
	// this DB error: "attempt to write readonly database". So I use in-memory database for all tests.
	if profile == TEST {
		initializeDatabaseWithSource(":memory:")
		Logger.Warn("initializing database only in-memory")
	} else {
		initializeDatabaseWithSource(databaseFile)
	}
}
