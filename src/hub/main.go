package main

import (
	"github.com/ocelot-cloud/shared"
	"net/http"
)

var (
	Logger       = shared.ProvideLogger()
	uploadPath   = "/api/upload"
	downloadPath = "/api/download/"
	port         = "8082"
	rootUrl      = "http://localhost:" + port
)

func main() {
	http.HandleFunc(uploadPath, uploadHandler)
	http.HandleFunc(downloadPath, downloadHandler)

	Logger.Info("Server started on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		// TODO Is server stop sometimes normal, e.g. when gracefully shutdown?
		Logger.Fatal("Server stopped: %v", err)
	}
}
