package main

import (
	"github.com/ocelot-cloud/shared"
	"net/http"
)

var (
	Logger       = shared.ProvideLogger()
	tagPath      = "/tags"
	uploadPath   = tagPath + "/upload"
	downloadPath = tagPath + "/download/"
	userPath     = "/users"
	appPath      = "/apps"
	port         = "8082"
	rootUrl      = "http://localhost:" + port
)

func main() {
	// TODO Maybe wrap gorilla/mux like in backend, apply a common security policy and put it shared module.
	http.HandleFunc(uploadPath, uploadHandler)
	http.HandleFunc(downloadPath, downloadHandler)
	http.HandleFunc(tagPath, tagHandler)
	http.HandleFunc(appPath, appHandler)
	http.HandleFunc(userPath, userHandler)

	Logger.Info("Server started on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		// TODO Is server stop sometimes normal, e.g. when gracefully shutdown?
		Logger.Fatal("Server stopped: %v", err)
	}
}
