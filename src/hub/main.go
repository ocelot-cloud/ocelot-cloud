package main

import (
	"github.com/ocelot-cloud/shared"
	"net/http"
	"os"
)

func init() {
	createDataDir()
	Logger = shared.ProvideLogger() // TODO dataDir should be moved to "shared". ProvideLogger should create the logs.txt in dataDir
}

func main() {
	initializeDatabase()

	// TODO Maybe wrap gorilla/mux like in backend, apply a common security policy and put it in shared module.
	// TODO apply middleware?
	http.HandleFunc(downloadPath, downloadHandler)
	http.HandleFunc(tagPath, tagHandler)
	http.HandleFunc(changePasswordPath, changePasswordHandler)

	http.HandleFunc(appPath, appHandler)
	http.HandleFunc(userPath, userHandler)
	http.HandleFunc(loginPath, loginHandler)

	// Registration process is excluded from security, so no middleware is used.
	http.HandleFunc(registrationPath, registrationHandler)

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
	useInMemoryDB := os.Getenv("USE_IN_MEMORY_DB")
	if useInMemoryDB == "true" {
		initializeDatabaseWithSource(":memory:")
	} else {
		initializeDatabaseWithSource(databaseFile)
	}
}

func applyMiddleware(handler func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler(w, r)
		} else {
			cookie, err := r.Cookie(cookieName)
			if err != nil {
				// TODO
			} else if cookie == nil {
				// TODO
			} else if repo.IsCookieValid(cookie.Value) {
				// TODO
			} else {
				// TODO Take request -> get user -> check if origin is correct.
				handler(w, r)
			}
		}
	}
}
