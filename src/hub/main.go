package main

import (
	"github.com/ocelot-cloud/shared"
	"net/http"
)

func init() {
	createDataDir()
	Logger = shared.ProvideLogger() // TODO dataDir should be moved to "shared". ProvideLogger should create the logs.txt in dataDir
}

func main() {
	initializeDatabase(databaseFile)
	// TODO Maybe wrap gorilla/mux like in backend, apply a common security policy and put it in shared module.
	http.HandleFunc(uploadPath, uploadHandler)     // TODO apply middleware
	http.HandleFunc(downloadPath, downloadHandler) // TODO apply middleware

	http.HandleFunc(tagPath, applyMiddleware(tagHandler))
	http.HandleFunc(appPath, applyMiddleware(appHandler))
	http.HandleFunc(userPath, applyMiddleware(userHandler))
	http.HandleFunc(loginPath, applyMiddleware(loginHandler))

	// Registration process is excluded from security, so no middleware is used.
	http.HandleFunc(registrationPath, registrationHandler)

	Logger.Info("Server started on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		// TODO Is server stop sometimes normal, e.g. when gracefully shutdown?
		Logger.Fatal("Server stopped: %v", err)
	}
}

func applyMiddleware(handler func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler(w, r)
		} else {
			// TODO Take cookie -> get user -> check if origin is correct.
			handler(w, r)
		}
	}
}
