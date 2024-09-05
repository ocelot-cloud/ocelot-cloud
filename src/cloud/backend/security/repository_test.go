package security

import (
	"github.com/ocelot-cloud/shared/assert"
	"os"
	"testing"
	"time"
)

var repo Repository = &MyRepository{}

func TestMain(m *testing.M) {
	InitializeDatabaseWithSource(":memory:")
	repo.WipeDatabase()
	code := m.Run()
	os.Exit(code)
}

const (
	sampleUser     = "user"
	samplePassword = "password"
	sampleCookie   = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
)

// TODO Finish SQLite Client Implementation And Tests
func TestSqliteClient(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleUser, samplePassword, false))
	assert.True(t, repo.IsPasswordCorrect(sampleUser, samplePassword))
	assert.False(t, repo.IsPasswordCorrect(sampleUser, samplePassword+"x"))

	assert.NotNil(t, repo.CreateUser(sampleUser, samplePassword+"x", false))

	assert.Nil(t, repo.DeleteUser(sampleUser))
	assert.False(t, repo.IsPasswordCorrect(sampleUser, samplePassword))
}

func TestDoesUserExist(t *testing.T) {
	defer repo.WipeDatabase()
	assert.False(t, repo.DoesUserExist(sampleUser))
	assert.Nil(t, repo.CreateUser(sampleUser, samplePassword, false))
	assert.True(t, repo.DoesUserExist(sampleUser))
	assert.Nil(t, repo.DeleteUser(sampleUser))
	assert.False(t, repo.DoesUserExist(sampleUser))
}

func TestGetUserWithCookie(t *testing.T) {
	defer repo.WipeDatabase()
	_, err := repo.GetUserWithCookie(sampleCookie)
	assert.NotNil(t, err)
	assert.Nil(t, repo.CreateUser(sampleUser, samplePassword, false))
	assert.Nil(t, repo.HashAndSaveCookie(sampleUser, sampleCookie, time.Now()))
	user, err := repo.GetUserWithCookie(sampleCookie)
	assert.Nil(t, err)
	assert.Equal(t, sampleUser, user)
}

func TestDoesAnyAdminUserExist(t *testing.T) {
	defer repo.WipeDatabase()
	assert.False(t, repo.DoesAnyAdminUserExist())
	assert.Nil(t, repo.CreateUser(sampleUser, samplePassword, false))
	assert.False(t, repo.DoesAnyAdminUserExist())
	assert.Nil(t, repo.CreateUser(sampleUser+"x", samplePassword, true))
	assert.True(t, repo.DoesAnyAdminUserExist())
}

// TODO check if expiration is working
// TODO can't set a cookie without user
// TODO all inconsistencies should be handled in this layer -> user does not exist, user already existing etc.
// TODO error: user already exists
// TODO SetCookie, DeleteCookie, IsCookieValid
