package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

var samplePassword = "mypassword"
var um UserManager

// TODO Finalize functionality
// TODO Add cases: "does not exist", "wrong password", "already existing"
// TODO "app already existing" applies only if the creating user has an app with that name
// TODO Duplications among multiple users are allowed for passwords and apps
// TODO Add "DeleteApp"
func TestStuff(t *testing.T) {
	initializeDatabase()
	um = &UserManagerSqlite{}
	defer resetDatabase(t)
	assert.False(t, um.DoesUserExist(sampleUser))
	err := um.CreateRepoUser(sampleUser, samplePassword)
	assert.Nil(t, err)
	assert.True(t, um.DoesUserExist(sampleUser))

	assert.True(t, um.IsPasswordCorrect(sampleUser, samplePassword))
	assert.False(t, um.IsPasswordCorrect(sampleUser, samplePassword+"x"))

	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.Nil(t, um.AddApp(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(sampleUser, sampleApp))

	err = um.DeleteRepoUser(sampleUser)
	assert.Nil(t, err)
	assert.False(t, um.DoesUserExist(sampleUser))
	// TODO add: assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
}

func resetDatabase(t *testing.T) {
	err := deleteIfExist(databaseFile)
	if err != nil {
		Logger.Error("Failed to delete database: %s, error: %v", databaseFile, err)
		t.Fail()
	}
	initializeDatabase()
}
