package main

import (
	"github.com/ocelot-cloud/shared"
	"testing"
)

var samplePassword = "mypassword"

// TODO Can "shared" module provider an "assert" package? -> assert.Equal() ...

// TODO Finalize.
func TestStuff(t *testing.T) {
	var a UserManager
	a = &UserManagerSqlite{}
	defer resetDatabase()

	shared.AssertFalse(t, a.DoesUserExist(sampleUser))
	err := a.CreateRepoUser(sampleUser, samplePassword)
	shared.AssertNil(t, err)
	// TODO shared.AssertTrue(t, a.DoesUserExist(sampleUser))
}

func resetDatabase() {
	// TODO
}
