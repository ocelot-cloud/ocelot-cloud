package main

import (
	"fmt"
	"github.com/ocelot-cloud/shared"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// TODO upload and download files, client logic and tests located in cloud, create repo, read repos and files, delete repos and files
// TODO Cloud + Hub: add accounts (sqlite?), GUI to self-register, login and handler logic, maybe put logic in a shared folder/module?, delete account
// TODO security: auth, tokens, upload only for logged in users and only to their repos, download is possible anonymously
// TODO structure: https://hub.ocelot-cloud.org/myuser_myapp_v1.0.tar.gz
// TODO Integration with cloud: acceptance test starts hub and cloud, cloud is told network location of hub, cloud initially has not a single app, but downloads it from hub during test
// TODO In "users" should be subdirectories like "users/myuser/myapp/v1.0"
// TODO Combine a complete story like: User registers account, logs in, uploads file etc...
// TODO Email verification fpr accounts?
// TODO How should uploads work? I imagine that the user simply-drags and drops stuff.
// TODO tar.gz should be unpacked on ocelot server for viewing the content. Only packed for transport?
// TODO protect against zip-bomb attack.
// TODO Introduce sqlite for user database: username, password-hash, salt, email, email verified -> maybe shared logic?

const uploadPath = "./users" // TODO Create folder if not exist

var logger = shared.ProvideLogger()

// TODO use paths that start with "/api/"
func main() {
	err := os.MkdirAll(uploadPath, os.ModePerm) // TODO Make permissions 600
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download/", downloadHandler)

	logger.Info("Server started on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
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

	out, err := os.Create(filepath.Join(uploadPath, header.Filename))
	if err != nil {
		logAndRespondError(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		logAndRespondError(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "File uploaded successfully") // TODO to be logged
}

func logAndRespondError(w http.ResponseWriter, msg string, httpStatus int) {
	logger.Error(msg)
	http.Error(w, msg, httpStatus)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/download/")
	if fileName == "" {
		logAndRespondError(w, "File name is missing", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(uploadPath, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logAndRespondError(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}
