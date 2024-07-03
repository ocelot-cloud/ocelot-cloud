package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"sort"
	"testing"
)

var samplePassword = "mypassword"
var um Repository = &UserManagerSqlite{}

func init() {
	resetDatabase()
}

func TestUserCreation(t *testing.T) {
	defer resetDatabase()
	assert.False(t, um.DoesUserExist(sampleUser))
	assert.Nil(t, um.CreateUser(sampleUser, samplePassword))
	assert.True(t, um.DoesUserExist(sampleUser))

	assert.Nil(t, um.DeleteUser(sampleUser))
	assert.False(t, um.DoesUserExist(sampleUser))
}

func TestCreateApp(t *testing.T) {
	defer resetDatabase()
	assert.Nil(t, um.CreateUser(sampleUser, samplePassword))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(sampleUser, sampleApp))
}

func TestDeleteAppCascadingThroughUser(t *testing.T) {
	defer resetDatabase()
	assert.Nil(t, um.CreateUser(sampleUser, samplePassword))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.Nil(t, um.DeleteApp(sampleUser, sampleApp))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
}

func TestDeleteAppDirectly(t *testing.T) {
	defer resetDatabase()
	assert.Nil(t, um.CreateUser(sampleUser, samplePassword))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.Nil(t, um.DeleteUser(sampleUser))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
}

func TestCantCreateUserTwice(t *testing.T) {
	defer resetDatabase()
	assert.Nil(t, um.CreateUser(sampleUser, samplePassword))
	assert.NotNil(t, um.CreateUser(sampleUser, samplePassword))
}

func TestCantCreateAppTwiceForSameUser(t *testing.T) {
	defer resetDatabase()
	assert.Nil(t, um.CreateUser(sampleUser, samplePassword))
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
	assert.Nil(t, um.CreateUser(sampleUser, samplePassword))
	assert.Nil(t, um.CreateUser(user2, samplePassword))
	assert.True(t, um.IsPasswordCorrect(sampleUser, samplePassword))
	assert.True(t, um.IsPasswordCorrect(user2, samplePassword))
}

func TestTolerateSameAppsForTwoUsers(t *testing.T) {
	defer resetDatabase()
	user2 := sampleUser + "2"
	assert.Nil(t, um.CreateUser(sampleUser, samplePassword))
	assert.Nil(t, um.CreateUser(user2, samplePassword))
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
	assert.Nil(t, um.CreateUser(sampleUser, samplePassword))
	assert.True(t, um.IsPasswordCorrect(sampleUser, samplePassword))
	assert.False(t, um.IsPasswordCorrect(sampleUser, samplePassword+"x"))
}

func TestSearch(t *testing.T) {
	defer resetDatabase()
	um.CreateUser(sampleUser, samplePassword)
	app1 := "prefix_myapp_suffix"
	app2 := "prefix_another-app_suffix"
	um.CreateApp(sampleUser, app1)
	um.CreateApp(sampleUser, app2)

	a, err := um.FindApps("app")
	assert.Nil(t, err)

	sort.Slice(a, func(i, j int) bool {
		return a[i].AppName < a[j].AppName
	})

	assert.Equal(t, 2, len(a))
	assert.Equal(t, sampleUser, a[0].Username)
	assert.Equal(t, sampleUser, a[1].Username)
	assert.Equal(t, app2, a[0].AppName)
	assert.Equal(t, app1, a[1].AppName)
}

func TestSearchNegative(t *testing.T) {
	defer resetDatabase()
	um.CreateUser(sampleUser, samplePassword)
	app := "prefix_myapp_suffix"
	um.CreateApp(sampleUser, app)

	a, err := um.FindApps("some")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(a))
}

func resetDatabase() {
	deleteIfExist(databaseFile)
	initializeDatabase()
}
