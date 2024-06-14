package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// TODO Add logger, also shared logic.
// TODO upload and download files, client logic and tests located in cloud, create repo, read repos and files, delete repos and files
// TODO Cloud + Hub: add accounts (sqlite?), GUI to self-register, login and handler logic, maybe put logic in a shared folder/module?, delete account
// TODO security: auth, tokens, upload only for logged in users and only to their repos, download is possible anonymously
// TODO structure: https://hub.ocelot-cloud.org/myuser_myapp_v1.0.tar.gz
// TODO Integration with cloud: acceptance test starts hub and cloud, cloud is told network location of hub, cloud initially has not a single app, but downloads it from hub during test

const uploadPath = "./users" // TODO Create folder if not exist

func main() {
	err := os.MkdirAll(uploadPath, os.ModePerm) // TODO Make permissions 600
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download/", downloadHandler)

	log.Println("Server started on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// TODO Add test
	// TODO Make security test that user and repo are in the name correctly, and that both exist.
	if !strings.HasSuffix(header.Filename, ".tar.gz") {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	out, err := os.Create(filepath.Join(uploadPath, header.Filename))
	if err != nil {
		log.Printf("Failed to save file: %v\n", err.Error())
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Printf("Failed to save file: %v\n", err.Error())
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "File uploaded successfully")
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/download/")
	if fileName == "" {
		http.Error(w, "File name is missing", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(uploadPath, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}
