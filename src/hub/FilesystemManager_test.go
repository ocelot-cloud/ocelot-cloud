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
// TODO Store all stuff in a "data" folder.

var (
	sampleUser    = "myuser"
	sampleApp     = "myapp"
	sampleTag     = "v0.0.1"
	usersDir      = "users"
	singleUserDir = usersDir + "/" + sampleUser
	appDir        = singleUserDir + "/" + sampleApp
	sampleFile    = appDir + fmt.Sprintf("/%s.tar.gz", sampleTag)
)

func TestFilesystemManager(t *testing.T) {
	shared.AssertNil(t, deleteIfExist(usersDir))
	shared.AssertFalse(t, doesFolderExist(usersDir))
	shared.AssertNil(t, CreateUser(sampleUser))
	shared.AssertTrue(t, doesFolderExist(singleUserDir))
	shared.AssertTrue(t, isFolderEmpty(singleUserDir))
	shared.AssertNil(t, CreateApp(sampleUser, sampleApp))
	shared.AssertTrue(t, doesFolderExist(appDir))
	shared.AssertTrue(t, isFolderEmpty(appDir))

	data := []byte("hello")
	buffer := bytes.NewBuffer(data)
	CreateTag(sampleUser, sampleApp, sampleTag, buffer) // TODO Should return error?
	shared.AssertTrue(t, doesFolderExist(appDir))
	shared.AssertEqual(t, "hello", getTagFileContent(sampleFile))
	DeleteTag(sampleUser, sampleApp, sampleTag)
	shared.AssertFalse(t, doesFileExist(sampleFile))

	DeleteApp(sampleUser, sampleApp)
	shared.AssertTrue(t, doesFolderExist(singleUserDir))
	shared.AssertFalse(t, doesFolderExist(appDir))
	DeleteUser(sampleUser)
	shared.AssertTrue(t, doesFolderExist(usersDir))
	shared.AssertFalse(t, doesFolderExist(singleUserDir))
	shared.AssertNil(t, deleteIfExist(usersDir))
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
