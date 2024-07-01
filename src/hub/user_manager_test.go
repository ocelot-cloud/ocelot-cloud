package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

var samplePassword = "mypassword"
var um UserManager = &UserManagerSqlite{}

func init() {
	resetDatabase()
}

// TODO Finalize functionality: findApp()

func TestUserCreation(t *testing.T) {
	defer resetDatabase()
	assert.False(t, um.DoesUserExist(sampleUser))
	assert.Nil(t, um.CreateRepoUser(sampleUser, samplePassword))
	assert.True(t, um.DoesUserExist(sampleUser))

	assert.Nil(t, um.DeleteRepoUser(sampleUser))
	assert.False(t, um.DoesUserExist(sampleUser))
}

func TestCreateApp(t *testing.T) {
	defer resetDatabase()
	assert.Nil(t, um.CreateRepoUser(sampleUser, samplePassword))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(sampleUser, sampleApp))
}

func TestDeleteAppCascadingThroughUser(t *testing.T) {
	defer resetDatabase()
	assert.Nil(t, um.CreateRepoUser(sampleUser, samplePassword))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.Nil(t, um.DeleteApp(sampleUser, sampleApp))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
}

func TestDeleteAppDirectly(t *testing.T) {
	defer resetDatabase()
	assert.Nil(t, um.CreateRepoUser(sampleUser, samplePassword))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.Nil(t, um.DeleteRepoUser(sampleUser))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
}

func TestCantCreateUserTwice(t *testing.T) {
	defer resetDatabase()
	assert.Nil(t, um.CreateRepoUser(sampleUser, samplePassword))
	assert.NotNil(t, um.CreateRepoUser(sampleUser, samplePassword))
}

func TestCantCreateAppTwiceForSameUser(t *testing.T) {
	defer resetDatabase()
	assert.Nil(t, um.CreateRepoUser(sampleUser, samplePassword))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	assert.NotNil(t, um.CreateApp(sampleUser, sampleApp))
}

func TestCantCreateAppWithoutUser(t *testing.T) {
	defer resetDatabase()
	assert.NotNil(t, um.CreateApp(sampleUser, sampleApp))
}

func TestTolerateSamePasswordForTwoUsers(t *testing.T) {
	defer resetDatabase()
	user2 := sampleUser + "2"
	assert.Nil(t, um.CreateRepoUser(sampleUser, samplePassword))
	assert.Nil(t, um.CreateRepoUser(user2, samplePassword))
	assert.True(t, um.IsPasswordCorrect(sampleUser, samplePassword))
	assert.True(t, um.IsPasswordCorrect(user2, samplePassword))
}

func TestTolerateSameAppsForTwoUsers(t *testing.T) {
	defer resetDatabase()
	user2 := sampleUser + "2"
	assert.Nil(t, um.CreateRepoUser(sampleUser, samplePassword))
	assert.Nil(t, um.CreateRepoUser(user2, samplePassword))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	assert.Nil(t, um.CreateApp(user2, sampleApp))

	assert.True(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(user2, sampleApp))

	assert.Nil(t, um.DeleteApp(sampleUser, sampleApp))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(user2, sampleApp))
}

func TestPasswordVerification(t *testing.T) {
	defer resetDatabase()
	assert.Nil(t, um.CreateRepoUser(sampleUser, samplePassword))
	assert.True(t, um.IsPasswordCorrect(sampleUser, samplePassword))
	assert.False(t, um.IsPasswordCorrect(sampleUser, samplePassword+"x"))
}

func TestSearch(t *testing.T) {
	defer resetDatabase()
	um.CreateRepoUser(sampleUser, samplePassword)
	um.CreateApp(sampleUser, "prefix_myapp_suffix")
	um.CreateApp(sampleUser, "prefix_another-app_suffix")
	a, err := um.FindApps("myapp")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(a))
	assert.Equal(t, sampleUser, a[0].Username)
	assert.Equal(t, "prefix_myapp_suffix", a[0].AppName)

	a, err = um.FindApps("app")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(a))
	assert.Equal(t, sampleUser, a[0].Username)
	assert.Equal(t, sampleUser, a[1].Username)
	assert.Equal(t, "prefix_myapp_suffix", a[0].AppName)
	assert.Equal(t, "prefix_another-app_suffix", a[1].AppName)
}

func resetDatabase() {
	deleteIfExist(databaseFile)
	initializeDatabase()
}
