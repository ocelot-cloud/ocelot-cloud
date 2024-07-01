package main

import (
	"bytes"
	"fmt"
	"github.com/ocelot-cloud/shared"
	"io"
	"net/http"
	"os"
	"strings"
)

// TODO General notes: Hub is only backend. Can be used via Cloud GUI, which directly addresses the hub API.

// TODO upload and download files, client logic and tests located in cloud, create repo, read repos and files, delete repos and files
// TODO Cloud + Hub: add accounts (sqlite?), GUI to self-register, login and handler logic, maybe put logic in a shared folder/module?, delete account
// TODO security: auth, tokens, upload only for logged in users and only to their repos, download is possible anonymously
// TODO structure: https://hub.ocelot-cloud.org/myuser_myapp_v1.0.tar.gz
// TODO Integration with cloud: acceptance test starts hub and cloud, cloud is told network location of hub, cloud initially has not a single app, but downloads it from hub during test
// TODO In "users" should be subdirectories like "users/myuser/myapp/v1.0"
// TODO Combine a complete story like: User registers account, logs in, uploads file etc...
// TODO Email verification fpr accounts? -> Maybe shift to production release issues?
// TODO How should uploads work? I imagine that the user simply-drags and drops stuff.
// TODO tar.gz should be unpacked on ocelot server for viewing the content. Only packed for transport?
// TODO protect against zip-bomb attack.
// TODO Introduce sqlite for user database: username, password-hash, salt, email, email verified -> maybe shared logic?
// TODO Can be deployed together with traefik to generate certs. Add "deploy hub" to ci-runner, also add docker-compose.yml. Maybe add a test server?
// TODO At the beginning always login in the local cloud. On first use of upload, login to hub. Cloud gets a token for future automatic logins.
// TODO When upload is implemented in hub, then I can delete alls the stacks in the cloud. Acceptance tests need to integrate hub and need to implement download of stacks at the beginning?

// TODO REST API
// security relevant:
//   create user/app/tag
//   delete user/app/tag
//   upload(app, tag) -> upload goes to the currently logged in user repo
// not security relevant:
//   search(app) -> may return many entries of the same app from different users
//   getTags(user, app)
//   download(user, app, tag)
// Design question: In my ocelot-cloud I want to add my credentials for the hub only once. So store the credentials there?
// Then implement the client in the cloud. Maybe run usage tests against it for simple scenarios.
// Implement input validation? only allow lowercase letters and underscores

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

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logAndRespondError(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		logAndRespondError(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// TODO Add test
	// TODO Make security test that user and repo are in the name correctly, and that both exist.
	if !strings.HasSuffix(header.Filename, ".tar.gz") {
		logAndRespondError(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// TODO Duplication
	fileInfo := strings.Split(header.Filename, "_")
	if len(fileInfo) != 3 {
		logAndRespondError(w, "Invalid file name", http.StatusBadRequest)
		return
	}
	username := fileInfo[0]
	app := fileInfo[1]
	tag, _ := strings.CutSuffix(fileInfo[2], ".tar.gz") // TODO Check error.

	var fileBuffer bytes.Buffer
	_, err = io.Copy(&fileBuffer, file)
	if err != nil {
		logAndRespondError(w, "Failed to read file content", http.StatusInternalServerError)
		return
	}

	err = CreateTag(username, app, tag, &fileBuffer)
	if err != nil {
		logAndRespondError(w, "Failed to write content to local file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	Logger.Info("File uploaded successfully: %s", header.Filename)
}

func logAndRespondError(w http.ResponseWriter, msg string, httpStatus int) {
	Logger.Error(msg)
	http.Error(w, msg, httpStatus)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	uploadName := strings.TrimPrefix(r.URL.Path, downloadPath)
	if uploadName == "" {
		logAndRespondError(w, "File name is missing", http.StatusBadRequest)
		return
	}

	fileInfo := strings.Split(uploadName, "_")
	if len(fileInfo) != 3 {
		logAndRespondError(w, "Invalid file name", http.StatusBadRequest)
		return
	}
	username := fileInfo[0]
	app := fileInfo[1]
	fileName := fileInfo[2]

	// TODO Should be returned by external function. Too low level here.
	path := fmt.Sprintf("data/users/%s/%s/%s", username, app, fileName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logAndRespondError(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, path)
}

type FileInfo struct {
	FileName string
	App      string
	Tag      string
}
