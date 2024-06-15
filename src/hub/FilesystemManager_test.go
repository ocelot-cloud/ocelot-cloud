package main

import (
	"bytes"
	"fmt"
	"github.com/ocelot-cloud/shared"
	"os"
	"testing"
)

// TODO Implement tests and corresponding functions.
// TODO In the middle -> Create/Delete tag, tar.gz file from in-memory

var (
	sampleUser = "myuser"
	sampleApp  = "myapp"
	sampleTag  = "v0.0.1"
	sampleFile = fmt.Sprintf("%s.tar.gz", sampleTag)
)

func TestFilesystemManager(t *testing.T) {
	shared.AssertNil(t, deleteIfExist("users"))
	shared.AssertFalse(t, doesFolderExist("users"))
	shared.AssertNil(t, CreateUser(sampleUser))
	shared.AssertTrue(t, doesFolderExist("users/myuser"))
	shared.AssertTrue(t, isFolderEmpty("users/myuser"))
	CreateApp(sampleUser, sampleApp)
	shared.AssertTrue(t, doesFolderExist("users/myuser/myapp"))
	shared.AssertTrue(t, isFolderEmpty("users/myuser/myapp"))

	data := []byte("hello")
	buffer := bytes.NewBuffer(data)
	CreateTag(sampleUser, sampleApp, sampleTag, buffer) // TODO Should return error?
	shared.AssertTrue(t, doesFolderExist("users/myuser/myapp"))
	shared.AssertEqual(t, "hello", getTagFileContent("users/myuser/myapp/v0.0.1"))
	DeleteTag(sampleUser, sampleApp, sampleTag)
	shared.AssertFalse(t, doesFileExist("users/myuser/myapp/v0.0.1"))

	DeleteApp(sampleUser, sampleApp)
	shared.AssertTrue(t, doesFolderExist("users/myuser"))
	shared.AssertFalse(t, doesFolderExist("users/myuser/myapp"))
	DeleteUser(sampleUser)
	shared.AssertTrue(t, doesFolderExist("users"))
	shared.AssertFalse(t, doesFolderExist("users/myuser"))
	shared.AssertNil(t, deleteIfExist("users"))
}

func doesFileExist(relativePath string) bool {
	info, err := os.Stat(relativePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
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

// TODO Create methods:
// getUserList
// getRepoList matching regex '*search-term*'
// getTagList
// limit to 100 elements? Allow search terms?
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
