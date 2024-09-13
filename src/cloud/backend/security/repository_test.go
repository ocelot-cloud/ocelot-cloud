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
	_, err := repo.GetUserViaCookie(sampleCookie)
	assert.NotNil(t, err)

	assert.Nil(t, repo.CreateUser(sampleUser, samplePassword, false))
	assert.Nil(t, repo.HashAndSaveCookie(sampleUser, sampleCookie, time.Now()))
	auth, err := repo.GetUserViaCookie(sampleCookie)
	assert.Nil(t, err)
	assert.Equal(t, sampleUser, auth.User)
	assert.False(t, auth.IsAdmin)
	repo.WipeDatabase()

	assert.Nil(t, repo.CreateUser(sampleUser, samplePassword, true))
	assert.Nil(t, repo.HashAndSaveCookie(sampleUser, sampleCookie, time.Now()))
	auth, err = repo.GetUserViaCookie(sampleCookie)
	assert.Nil(t, err)
	assert.Equal(t, sampleUser, auth.User)
	assert.True(t, auth.IsAdmin)
}

func TestDoesAnyAdminUserExist(t *testing.T) {
	defer repo.WipeDatabase()
	assert.False(t, repo.DoesAnyAdminUserExist())
	assert.Nil(t, repo.CreateUser(sampleUser, samplePassword, false))
	assert.False(t, repo.DoesAnyAdminUserExist())
	assert.Nil(t, repo.CreateUser(sampleUser+"x", samplePassword, true))
	assert.True(t, repo.DoesAnyAdminUserExist())
}

func TestLogout(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleUser, samplePassword, false))
	assert.Nil(t, repo.HashAndSaveCookie(sampleUser, sampleCookie, time.Now()))
	auth, err := repo.GetUserViaCookie(sampleCookie)
	assert.Nil(t, err)
	assert.Equal(t, sampleUser, auth.User)

	assert.Nil(t, repo.Logout(sampleUser))
	assert.True(t, repo.DoesUserExist(sampleUser))
	auth, err = repo.GetUserViaCookie(sampleCookie)
	assert.NotNil(t, err)
	assert.Nil(t, auth)
}

func TestChangePassword(t *testing.T) {
	defer repo.WipeDatabase()
	oldPassword := samplePassword
	newPassword := samplePassword + "x"
	assert.Nil(t, repo.CreateUser(sampleUser, oldPassword, false))
	assert.True(t, repo.IsPasswordCorrect(sampleUser, oldPassword))

	assert.Nil(t, repo.ChangePassword(sampleUser, newPassword))
	assert.False(t, repo.IsPasswordCorrect(sampleUser, oldPassword))
	assert.True(t, repo.IsPasswordCorrect(sampleUser, newPassword))
}

// TODO check if expiration is working
// TODO can't set a cookie without user
// TODO all inconsistencies should be handled in this layer -> user does not exist, user already existing etc.
// TODO error: user already exists
// TODO SetCookie, DeleteCookie, IsCookieValid
// TODO the DB interface appears to grow quite large when all all use cases are implemented. Check if could be split up.
// TODO Test deletion cascading, e.g. deleting user should also delete his group memberships etc.
