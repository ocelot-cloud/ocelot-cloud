package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateUser(username string) {
	userDir := filepath.Join("users", username)
	if err := os.MkdirAll(userDir, os.ModePerm); err != nil {
		logger.Error("Error creating user directory: %v", err)
	}
}

func DeleteUser(username string) {
	userDir := filepath.Join("users", username)
	if err := deleteIfExist(userDir); err != nil {
		logger.Error("Error deleting user directory: %v", err)
	}
}

func CreateApp(username, app string) {
	appDir := filepath.Join("users", username, app)
	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		logger.Error("Error creating app directory: %v", err)
	}
}

func DeleteApp(username, app string) {
	appDir := filepath.Join("users", username, app)
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
