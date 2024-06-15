package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	dataDir  = "data"
	usersDir = dataDir + "/users"
)

func init() {
	setup()
}

func setup() {
	if _, err := os.Stat(usersDir); os.IsNotExist(err) {
		if err := os.MkdirAll(usersDir, os.ModePerm); err != nil {
			// TODO duplication
			logger.Error("Error creating users directory: %v", err)
			os.Exit(1)
		}
	}
}

func CreateUser(username string) error {
	userDir := filepath.Join(usersDir, username)
	if err := os.MkdirAll(userDir, os.ModePerm); err != nil {
		// TODO duplication
		logger.Error("Error creating user directory: %v", err)
		return fmt.Errorf("Error creating user directory: %v", err)
	}
	return nil
}

func DeleteUser(username string) {
	userDir := filepath.Join(usersDir, username)
	if err := deleteIfExist(userDir); err != nil {
		logger.Error("Error deleting user directory: %v", err)
	}
}

func CreateApp(username, app string) error {
	appDir := filepath.Join(usersDir, username, app)
	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		// TODO duplication
		logger.Error("Error creating app directory: %v", err)
		return fmt.Errorf("Error creating app directory: %v", err)
	}
	return nil
}

func DeleteApp(username, app string) {
	appDir := filepath.Join(usersDir, username, app)
	if err := deleteIfExist(appDir); err != nil {
		logger.Error("Error deleting app directory: %v", err)
	}
}

func deleteIfExist(path string) error {
	if exists(path) {
		err := os.RemoveAll(path)
		if err != nil {
			logger.Error("Error deleting file: %v", err)
			return fmt.Errorf("failed to delete %s: %v", path, err)
		}
	}
	return nil
}

func exists(relativePath string) bool {
	_, err := os.Stat(relativePath)
	return !os.IsNotExist(err)
}

func GetUserList() []string {
	return getSubFolderNamesFromFolder(usersDir)
}

func getSubFolderNamesFromFolder(relativePath string) []string {
	var fileNames []string
	files, err := os.ReadDir(relativePath)
	if err != nil {
		return nil
	}
	for _, file := range files {
		if file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames
}

func CreateTag(user string, app string, tag string, buffer *bytes.Buffer) {
	tagFilePath := filepath.Join(usersDir, user, app, fmt.Sprintf("%s.tar.gz", tag))
	file, err := os.Create(tagFilePath)
	if err != nil {
		logger.Error("Error creating tag file: %v", err)
		return
	}
	defer file.Close()

	if _, err := io.Copy(file, buffer); err != nil {
		logger.Error("Error writing to tag file: %v", err)
	}
}

func DeleteTag(user string, app string, tag string) {
	tagFilePath := filepath.Join(usersDir, user, app, fmt.Sprintf("%s.tar.gz", tag))
	if err := deleteIfExist(tagFilePath); err != nil {
		logger.Error("Error deleting tag file: %v", err)
	}
}

func getTagFileContent(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		logger.Error("Error reading file content: %v", err)
		return ""
	}
	return string(data)
}
