package main

import (
	"bytes"
	"fmt"
	"github.com/ocelot-cloud/shared"
	"os"
	"testing"
)

// TODO Implement tests and corresponding functions.
// TODO Store all stuff in a "data" folder.

var (
	sampleUser    = "myuser"
	sampleApp     = "myapp"
	sampleTag     = "v0.0.1"
	singleUserDir = usersDir + "/" + sampleUser
	appDir        = singleUserDir + "/" + sampleApp
	sampleFile    = appDir + fmt.Sprintf("/%s.tar.gz", sampleTag)
)

func TestFilesystemManager(t *testing.T) {
	setup()
	defer cleanup()
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
	exists, isDir := pathExists(relativePath)
	return exists && !isDir
}

func pathExists(path string) (exists bool, isDir bool) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, false
	}
	return true, info.IsDir()
}

func doesFolderExist(relativePath string) bool {
	exists, isDir := pathExists(relativePath)
	return exists && isDir
}

func isFolderEmpty(relativePath string) bool {
	files, err := os.ReadDir(relativePath)
	if err != nil {
		return false
	}
	return len(files) == 0
}

// TODO Create methods:
// getRepoList matching regex '*search-term*'
// getTagList
// limit to 100 elements? Allow search terms?
func TestReading(t *testing.T) {
	setup()
	defer cleanup()
	shared.AssertEqual(t, 0, len(GetUserList()))
	shared.AssertNil(t, CreateUser(sampleUser))
	users := GetUserList()
	shared.AssertEqual(t, 1, len(users))
	shared.AssertEqual(t, sampleUser, users[0])

	sampleUser2 := sampleUser + "2"
	shared.AssertNil(t, CreateUser(sampleUser2))
	users = GetUserList()
	shared.AssertEqual(t, 2, len(users))
	shared.AssertEqual(t, sampleUser, users[0])
	shared.AssertEqual(t, sampleUser2, users[1])
}

func TestSetup(t *testing.T) {
	cleanup()
	shared.AssertFalse(t, doesFolderExist(usersDir))
	setup()
	shared.AssertTrue(t, doesFolderExist(usersDir))
	cleanup()
}

func cleanup() {
	deleteIfExist(dataDir)
}
