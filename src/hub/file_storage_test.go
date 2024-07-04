package main

import (
	"bytes"
	"fmt"
	"github.com/ocelot-cloud/shared/assert"
	"os"
	"testing"
)

var (
	sampleUser                    = "myuser"
	sampleApp                     = "myapp"
	sampleTag                     = "v0.0.1"
	singleUserDir                 = usersDir + "/" + sampleUser
	appDir                        = singleUserDir + "/" + sampleApp
	sampleFile                    = appDir + fmt.Sprintf("/%s.tar.gz", sampleTag)
	sampleTaggedFileContentBuffer = bytes.NewBuffer([]byte("hello"))
	sampleFileInfo                = &FileInfo{sampleUser, sampleApp, sampleTag, sampleFile}
)

func TestFilesystemManager(t *testing.T) {
	createDataDir()
	defer cleanup()
	assert.Nil(t, fs.CreateUser(sampleUser))
	assert.True(t, doesFolderExist(singleUserDir))
	assert.True(t, isFolderEmpty(singleUserDir))

	err := fs.CreateUser(sampleUser)
	assert.NotNil(t, err)
	expectedErrorMessage := fmt.Sprintf("User already exists: %s", sampleUser)
	assert.Equal(t, expectedErrorMessage, err.Error())

	assert.Nil(t, fs.CreateApp(sampleUser, sampleApp))
	assert.True(t, doesFolderExist(appDir))
	assert.True(t, isFolderEmpty(appDir))

	err = fs.CreateApp(sampleUser, sampleApp)
	assert.NotNil(t, err)
	expectedErrorMessage = fmt.Sprintf("App '%s' of user '%s' already exists", sampleApp, sampleUser)
	assert.Equal(t, expectedErrorMessage, err.Error())

	assert.Nil(t, fs.CreateTag(sampleFileInfo, sampleTaggedFileContentBuffer))
	assert.True(t, doesFolderExist(appDir))
	assert.Equal(t, "hello", getTagFileContent(sampleFile))

	err = fs.CreateTag(sampleFileInfo, sampleTaggedFileContentBuffer)
	assert.NotNil(t, err)
	expectedErrorMessage = fmt.Sprintf("Tag '%s' of app '%s' of user '%s' already exists", sampleTag, sampleApp, sampleUser)
	assert.Equal(t, expectedErrorMessage, err.Error())

	fs.DeleteTag(sampleUser, sampleApp, sampleTag)
	assert.False(t, doesFileExist(sampleFile))

	fs.DeleteApp(sampleUser, sampleApp)
	assert.True(t, doesFolderExist(singleUserDir))
	assert.False(t, doesFolderExist(appDir))

	fs.DeleteUser(sampleUser)
	assert.True(t, doesFolderExist(usersDir))
	assert.False(t, doesFolderExist(singleUserDir))
	assert.Nil(t, deleteIfExist(usersDir))
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

func TestReadingUsers(t *testing.T) {
	createDataDir()
	defer cleanup()
	assert.Equal(t, 0, len(GetUserList()))
	assert.Nil(t, fs.CreateUser(sampleUser))
	users := GetUserList()
	assert.Equal(t, 1, len(users))
	assert.Equal(t, sampleUser, users[0])

	sampleUser2 := sampleUser + "2"
	assert.Nil(t, fs.CreateUser(sampleUser2))
	users = GetUserList()
	assert.Equal(t, 2, len(users))
	assert.Equal(t, sampleUser, users[0])
	assert.Equal(t, sampleUser2, users[1])
}

func TestSetup(t *testing.T) {
	cleanup()
	assert.False(t, doesFolderExist(usersDir))
	createDataDir()
	assert.True(t, doesFolderExist(usersDir))
	cleanup()
}

func TestReadingApps(t *testing.T) {
	defer cleanup()
	list, err := GetAppList(sampleUser)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(list))

	assert.Nil(t, fs.CreateUser(sampleUser))
	list, err = GetAppList(sampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(list))

	assert.Nil(t, fs.CreateApp(sampleUser, sampleApp))
	list, err = GetAppList(sampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, sampleApp, list[0])
}

func TestReadingTags(t *testing.T) {
	defer cleanup()
	list, err := GetTagList(sampleUser, sampleApp)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(list))

	assert.Nil(t, fs.CreateUser(sampleUser))
	assert.Nil(t, fs.CreateApp(sampleUser, sampleApp))
	list, err = GetTagList(sampleUser, sampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(list))

	assert.Nil(t, fs.CreateTag(sampleFileInfo, sampleTaggedFileContentBuffer))
	list, err = GetTagList(sampleUser, sampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, sampleTag, list[0])
}

func cleanup() {
	err := deleteIfExist(dataDir)
	if err != nil {
		Logger.Error("Cleanup: Could not delete dir: %s", dataDir)
		os.Exit(1)
	}
}

func TestParentNotFound(t *testing.T) {
	createDataDir()
	defer cleanup()

	err := fs.CreateApp(sampleUser, sampleApp)
	assert.NotNil(t, err)
	expectedErrorMessage := fmt.Sprintf("User '%s' does not exist", sampleUser)
	assert.Equal(t, expectedErrorMessage, err.Error())
	assert.Nil(t, fs.DeleteApp(sampleUser, sampleApp))

	assert.Nil(t, fs.CreateUser(sampleUser))
	err = fs.CreateTag(sampleFileInfo, sampleTaggedFileContentBuffer)
	assert.NotNil(t, err)
	expectedErrorMessage = fmt.Sprintf("App '%s' of user '%s' does not exist", sampleApp, sampleUser)
	assert.Equal(t, expectedErrorMessage, err.Error())
	assert.Nil(t, fs.DeleteApp(sampleUser, sampleApp))
}

// TODO Apply consistent naming to packages, files and types.
