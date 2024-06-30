package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
			logger.Error("Error creating users directory: %v. Terminating application.", err)
			os.Exit(1)
		}
	}
}

func CreateUser(username string) error {
	userDir := filepath.Join(usersDir, username)

	if _, err := os.Stat(userDir); err == nil {
		return logger.LogAndReturnError("User already exists: %v", username)
	}

	if err := os.MkdirAll(userDir, os.ModePerm); err != nil {
		return logger.LogAndReturnError("Error creating user directory: %v", err)
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

	if _, err := os.Stat(appDir); err == nil {
		return logger.LogAndReturnError("App '%s' of user '%s' already exists", app, username)
	}

	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		return logger.LogAndReturnError("Error creating app directory: %v", err)
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
	return getFilesFromFolder(usersDir, true)
}

func getFilesFromFolder(relativePath string, isFolder bool) []string {
	var fileNames []string
	files, err := os.ReadDir(relativePath)
	if err != nil {
		return nil
	}
	for _, file := range files {
		if file.IsDir() == isFolder {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames
}

func GetAppList(username string) ([]string, error) {
	if doesUserExist(username) {
		return getFilesFromFolder(usersDir+"/"+username, true), nil
	} else {
		return nil, fmt.Errorf("user '%s' not found", username)
	}
}

func GetTagList(username string, app string) ([]string, error) {
	if doesUserExist(username) {
		if doesAppExist(username, app) {
			tagFileNames := getFilesFromFolder(usersDir+"/"+username+"/"+app, false)
			var tags []string
			for _, tagFileName := range tagFileNames {
				tags = append(tags, strings.TrimSuffix(tagFileName, ".tar.gz"))
			}
			return tags, nil
		}
		return nil, fmt.Errorf("app '%s' not found", app)
	} else {
		return nil, fmt.Errorf("user '%s' not found", username)
	}
}

func doesAppExist(username string, app string) bool {
	appList, err := GetAppList(username)
	if err != nil {
		return false
	}

	for _, v := range appList {
		if v == app {
			return true
		}
	}
	return false
}

func doesUserExist(username string) bool {
	userList := GetUserList()
	for _, v := range userList {
		if v == username {
			return true
		}
	}
	return false
}

func CreateTag(user string, app string, tag string, buffer *bytes.Buffer) error {
	tagFilePath := filepath.Join(usersDir, user, app, fmt.Sprintf("%s.tar.gz", tag))

	if _, err := os.Stat(tagFilePath); err == nil {
		return logger.LogAndReturnError("Tag '%s' of app '%s' of user '%s' already exists", tag, app, user)
	}

	file, err := os.Create(tagFilePath)

	if err != nil {
		return logger.LogAndReturnError("Error creating tag file: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, buffer); err != nil {
		return logger.LogAndReturnError("Error writing to tag file: %v", err)
	}
	return nil
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
