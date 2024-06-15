package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CreateUser(username string) error {
	usersDir := "users"
	if _, err := os.Stat(usersDir); os.IsNotExist(err) {
		if err := os.MkdirAll(usersDir, os.ModePerm); err != nil {
			// TODO duplication
			logger.Error("Error creating users directory: %v", err)
			return fmt.Errorf("Error creating users directory: %v", err)
		}
	}

	userDir := filepath.Join(usersDir, username)
	if err := os.MkdirAll(userDir, os.ModePerm); err != nil {
		// TODO duplication
		logger.Error("Error creating user directory: %v", err)
		return fmt.Errorf("Error creating user directory: %v", err)
	}
	return nil
}

func DeleteUser(username string) {
	userDir := filepath.Join("users", username)
	if err := deleteIfExist(userDir); err != nil {
		logger.Error("Error deleting user directory: %v", err)
	}
}

func CreateApp(username, app string) error {
	appDir := filepath.Join("users", username, app)
	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		// TODO duplication
		logger.Error("Error creating app directory: %v", err)
		return fmt.Errorf("Error creating app directory: %v", err)
	}
	return nil
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

func CreateTag(user string, app string, tag string, buffer *bytes.Buffer) {
	tagFilePath := filepath.Join("users", user, app, fmt.Sprintf("%s.tar.gz", tag))
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
	tagFilePath := filepath.Join("users", user, app, fmt.Sprintf("%s.tar.gz", tag))
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
