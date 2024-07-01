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

// TODO Finalize functionality
// TODO Add cases: "does not exist", "wrong password", "already existing"
// TODO "app already existing" applies only if the creating user has an app with that name
// TODO Duplications among multiple users are allowed for passwords and apps
// TODO Add "DeleteApp"
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
	assert.Nil(t, um.AddApp(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(sampleUser, sampleApp))
	// TODO add: assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
}

func TestPasswordVerification(t *testing.T) {
	defer resetDatabase()
	assert.Nil(t, um.CreateRepoUser(sampleUser, samplePassword))
	assert.True(t, um.IsPasswordCorrect(sampleUser, samplePassword))
	assert.False(t, um.IsPasswordCorrect(sampleUser, samplePassword+"x"))
}

func resetDatabase() {
	deleteIfExist(databaseFile)
	initializeDatabase()
}
