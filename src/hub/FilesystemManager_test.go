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
	sampleUser                    = "myuser"
	sampleApp                     = "myapp"
	sampleTag                     = "v0.0.1"
	singleUserDir                 = usersDir + "/" + sampleUser
	appDir                        = singleUserDir + "/" + sampleApp
	sampleFile                    = appDir + fmt.Sprintf("/%s.tar.gz", sampleTag)
	sampleTaggedFileContentBuffer = bytes.NewBuffer([]byte("hello"))
)

func TestFilesystemManager(t *testing.T) {
	setup()
	defer cleanup()
	shared.AssertNil(t, CreateUser(sampleUser))
	shared.AssertTrue(t, doesFolderExist(singleUserDir))
	err := CreateUser(sampleUser)
	shared.AssertNotNil(t, err)
	expectedErrorMessage := fmt.Sprintf("User already exists: %s", sampleUser)
	shared.AssertEqual(t, expectedErrorMessage, err.Error())

	shared.AssertTrue(t, isFolderEmpty(singleUserDir))
	shared.AssertNil(t, CreateApp(sampleUser, sampleApp))
	shared.AssertTrue(t, doesFolderExist(appDir))
	shared.AssertTrue(t, isFolderEmpty(appDir))

	shared.AssertNil(t, CreateTag(sampleUser, sampleApp, sampleTag, sampleTaggedFileContentBuffer))
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

func doesFileExist(path string) bool {
	exists, isDir := pathExists(path)
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

// TODO Handle cases for "already existing" and "not found". e.g. user/app/tag does not exist

func TestReadingUsers(t *testing.T) {
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

func TestReadingApps(t *testing.T) {
	defer cleanup()
	list, err := GetAppList(sampleUser)
	shared.AssertNotNil(t, err)
	shared.AssertEqual(t, 0, len(list))

	shared.AssertNil(t, CreateUser(sampleUser))
	list, err = GetAppList(sampleUser)
	shared.AssertNil(t, err)
	shared.AssertEqual(t, 0, len(list))

	shared.AssertNil(t, CreateApp(sampleUser, sampleApp))
	list, err = GetAppList(sampleUser)
	shared.AssertNil(t, err)
	shared.AssertEqual(t, 1, len(list))
	shared.AssertEqual(t, sampleApp, list[0])
}

func TestReadingTags(t *testing.T) {
	defer cleanup()
	list, err := GetTagList(sampleUser, sampleApp)
	shared.AssertNotNil(t, err)
	shared.AssertEqual(t, 0, len(list))

	shared.AssertNil(t, CreateUser(sampleUser))
	shared.AssertNil(t, CreateApp(sampleUser, sampleApp))
	list, err = GetTagList(sampleUser, sampleApp)
	shared.AssertNil(t, err)
	shared.AssertEqual(t, 0, len(list))

	shared.AssertNil(t, CreateTag(sampleUser, sampleApp, sampleTag, sampleTaggedFileContentBuffer))
	list, err = GetTagList(sampleUser, sampleApp)
	shared.AssertNil(t, err)
	shared.AssertEqual(t, 1, len(list))
	shared.AssertEqual(t, sampleTag, list[0])
}

func cleanup() {
	err := deleteIfExist(dataDir)
	if err != nil {
		logger.Error("Cleanup: Could not delete dir: %s", dataDir)
		os.Exit(1)
	}
}
