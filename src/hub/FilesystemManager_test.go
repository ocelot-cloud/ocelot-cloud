package main

import (
	"fmt"
	"github.com/ocelot-cloud/shared"
	"os"
	"testing"
)

// TODO Implement tests and corresponding functions.
// TODO In the middle -> Create/Delete tag, tar.gz file from in-memory

var sampleUser = "myuser"
var sampleApp = "myapp"
var sampleTag = "v0.0.1"
var sampleFile, _ = fmt.Printf("%s.tar.gz", sampleTag)

func TestFilesystemManager(t *testing.T) {
	shared.AssertNil(t, deleteIfExist("users"))
	shared.AssertFalse(t, doesFolderExist("users"))
	CreateUser(sampleUser)
	shared.AssertTrue(t, doesFolderExist("users/myuser"))
	shared.AssertTrue(t, isFolderEmpty("users/myuser"))
	CreateApp(sampleUser, sampleApp)
	shared.AssertTrue(t, doesFolderExist("users/myuser/myapp"))
	shared.AssertTrue(t, isFolderEmpty("users/myuser/myapp"))

	DeleteApp(sampleUser, sampleApp)
	shared.AssertTrue(t, doesFolderExist("users/myuser"))
	shared.AssertFalse(t, doesFolderExist("users/myuser/myapp"))
	DeleteUser(sampleUser)
	shared.AssertTrue(t, doesFolderExist("users"))
	shared.AssertFalse(t, doesFolderExist("users/myuser"))
	shared.AssertNil(t, deleteIfExist("users"))
}

func doesFolderExist(relativePath string) bool {
	info, err := os.Stat(relativePath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func isFolderEmpty(relativePath string) bool {
	files, err := os.ReadDir(relativePath)
	if err != nil {
		return false
	}
	return len(files) == 0
}

func getFileNamesFromFolder(relativePath string) []string {
	var fileNames []string
	files, err := os.ReadDir(relativePath)
	if err != nil {
		return fileNames // return an empty slice if there's an error
	}
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames
}
