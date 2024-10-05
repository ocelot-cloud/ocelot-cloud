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
	dbRepo.WipeDatabase()
	code := m.Run()
	os.Exit(code)
}

var (
	sampleUser     = "user"
	samplePassword = "password"
	sampleCookie   = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
)

// TODO Finish SQLite Client Implementation And Tests
func TestSqliteClient(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, userRepo.CreateUser(sampleUser, samplePassword, false))
	assert.True(t, userRepo.IsPasswordCorrect(sampleUser, samplePassword))
	assert.False(t, userRepo.IsPasswordCorrect(sampleUser, samplePassword+"x"))

	assert.NotNil(t, userRepo.CreateUser(sampleUser, samplePassword+"x", false))

	assert.Nil(t, userRepo.DeleteUser(sampleUser))
	assert.False(t, userRepo.IsPasswordCorrect(sampleUser, samplePassword))
}

func TestDoesUserExist(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.False(t, userRepo.DoesUserExist(sampleUser))
	assert.Nil(t, userRepo.CreateUser(sampleUser, samplePassword, false))
	assert.True(t, userRepo.DoesUserExist(sampleUser))
	assert.Nil(t, userRepo.DeleteUser(sampleUser))
	assert.False(t, userRepo.DoesUserExist(sampleUser))
}

func TestGetUserWithCookie(t *testing.T) {
	defer dbRepo.WipeDatabase()
	_, err := userRepo.GetUserViaCookie(sampleCookie)
	assert.NotNil(t, err)

	assert.Nil(t, userRepo.CreateUser(sampleUser, samplePassword, false))
	assert.Nil(t, userRepo.HashAndSaveCookie(sampleUser, sampleCookie, time.Now()))
	auth, err := userRepo.GetUserViaCookie(sampleCookie)
	assert.Nil(t, err)
	assert.Equal(t, sampleUser, auth.User)
	assert.False(t, auth.IsAdmin)
	dbRepo.WipeDatabase()

	assert.Nil(t, userRepo.CreateUser(sampleUser, samplePassword, true))
	assert.Nil(t, userRepo.HashAndSaveCookie(sampleUser, sampleCookie, time.Now()))
	auth, err = userRepo.GetUserViaCookie(sampleCookie)
	assert.Nil(t, err)
	assert.Equal(t, sampleUser, auth.User)
	assert.True(t, auth.IsAdmin)
}

func TestDoesAnyAdminUserExist(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.False(t, userRepo.DoesAnyAdminUserExist())
	assert.Nil(t, userRepo.CreateUser(sampleUser, samplePassword, false))
	assert.False(t, userRepo.DoesAnyAdminUserExist())
	assert.Nil(t, userRepo.CreateUser(sampleUser+"x", samplePassword, true))
	assert.True(t, userRepo.DoesAnyAdminUserExist())
}

func TestLogout(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, userRepo.CreateUser(sampleUser, samplePassword, false))
	assert.Nil(t, userRepo.HashAndSaveCookie(sampleUser, sampleCookie, time.Now()))
	auth, err := userRepo.GetUserViaCookie(sampleCookie)
	assert.Nil(t, err)
	assert.Equal(t, sampleUser, auth.User)

	assert.Nil(t, userRepo.Logout(sampleUser))
	assert.True(t, userRepo.DoesUserExist(sampleUser))
	auth, err = userRepo.GetUserViaCookie(sampleCookie)
	assert.NotNil(t, err)
	assert.Nil(t, auth)
}

func TestChangePassword(t *testing.T) {
	defer dbRepo.WipeDatabase()
	oldPassword := samplePassword
	newPassword := samplePassword + "x"
	assert.Nil(t, userRepo.CreateUser(sampleUser, oldPassword, false))
	assert.True(t, userRepo.IsPasswordCorrect(sampleUser, oldPassword))

	assert.Nil(t, userRepo.ChangePassword(sampleUser, newPassword))
	assert.False(t, userRepo.IsPasswordCorrect(sampleUser, oldPassword))
	assert.True(t, userRepo.IsPasswordCorrect(sampleUser, newPassword))
}
