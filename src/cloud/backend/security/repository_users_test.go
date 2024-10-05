package security

import (
	"github.com/ocelot-cloud/shared/assert"
	"os"
	"testing"
	"time"
)

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
	assert.Nil(t, UserRepo.CreateUser(sampleUser, samplePassword, false))
	assert.True(t, UserRepo.IsPasswordCorrect(sampleUser, samplePassword))
	assert.False(t, UserRepo.IsPasswordCorrect(sampleUser, samplePassword+"x"))

	assert.NotNil(t, UserRepo.CreateUser(sampleUser, samplePassword+"x", false))

	assert.Nil(t, UserRepo.DeleteUser(sampleUser))
	assert.False(t, UserRepo.IsPasswordCorrect(sampleUser, samplePassword))
}

func TestDoesUserExist(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.False(t, UserRepo.DoesUserExist(sampleUser))
	assert.Nil(t, UserRepo.CreateUser(sampleUser, samplePassword, false))
	assert.True(t, UserRepo.DoesUserExist(sampleUser))
	assert.Nil(t, UserRepo.DeleteUser(sampleUser))
	assert.False(t, UserRepo.DoesUserExist(sampleUser))
}

func TestGetUserWithCookie(t *testing.T) {
	defer dbRepo.WipeDatabase()
	_, err := UserRepo.GetUserViaCookie(sampleCookie)
	assert.NotNil(t, err)

	assert.Nil(t, UserRepo.CreateUser(sampleUser, samplePassword, false))
	assert.Nil(t, UserRepo.HashAndSaveCookie(sampleUser, sampleCookie, time.Now()))
	auth, err := UserRepo.GetUserViaCookie(sampleCookie)
	assert.Nil(t, err)
	assert.Equal(t, sampleUser, auth.User)
	assert.False(t, auth.IsAdmin)
	dbRepo.WipeDatabase()

	assert.Nil(t, UserRepo.CreateUser(sampleUser, samplePassword, true))
	assert.Nil(t, UserRepo.HashAndSaveCookie(sampleUser, sampleCookie, time.Now()))
	auth, err = UserRepo.GetUserViaCookie(sampleCookie)
	assert.Nil(t, err)
	assert.Equal(t, sampleUser, auth.User)
	assert.True(t, auth.IsAdmin)
}

func TestDoesAnyAdminUserExist(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.False(t, UserRepo.DoesAnyAdminUserExist())
	assert.Nil(t, UserRepo.CreateUser(sampleUser, samplePassword, false))
	assert.False(t, UserRepo.DoesAnyAdminUserExist())
	assert.Nil(t, UserRepo.CreateUser(sampleUser+"x", samplePassword, true))
	assert.True(t, UserRepo.DoesAnyAdminUserExist())
}

func TestLogout(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, UserRepo.CreateUser(sampleUser, samplePassword, false))
	assert.Nil(t, UserRepo.HashAndSaveCookie(sampleUser, sampleCookie, time.Now()))
	auth, err := UserRepo.GetUserViaCookie(sampleCookie)
	assert.Nil(t, err)
	assert.Equal(t, sampleUser, auth.User)

	assert.Nil(t, UserRepo.Logout(sampleUser))
	assert.True(t, UserRepo.DoesUserExist(sampleUser))
	auth, err = UserRepo.GetUserViaCookie(sampleCookie)
	assert.NotNil(t, err)
	assert.Nil(t, auth)
}

func TestChangePassword(t *testing.T) {
	defer dbRepo.WipeDatabase()
	oldPassword := samplePassword
	newPassword := samplePassword + "x"
	assert.Nil(t, UserRepo.CreateUser(sampleUser, oldPassword, false))
	assert.True(t, UserRepo.IsPasswordCorrect(sampleUser, oldPassword))

	assert.Nil(t, UserRepo.ChangePassword(sampleUser, newPassword))
	assert.False(t, UserRepo.IsPasswordCorrect(sampleUser, oldPassword))
	assert.True(t, UserRepo.IsPasswordCorrect(sampleUser, newPassword))
}

func TestSecrets(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, UserRepo.CreateUser(sampleUser, samplePassword, false))
	secret, err := UserRepo.GenerateSecret(sampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 64, len(secret))

	assert.Nil(t, UserRepo.CreateUser(sampleUser+"x", samplePassword, false))
	secret2, _ := UserRepo.GenerateSecret(sampleUser + "x")
	assert.NotEqual(t, secret, secret2)

	assert.False(t, UserRepo.IsSecretCorrect(sampleUser, secret+"x"))
	assert.False(t, UserRepo.IsSecretCorrect(sampleUser+"x", secret))
	assert.True(t, UserRepo.IsSecretCorrect(sampleUser, secret))
	assert.Nil(t, UserRepo.RemoveSecret(sampleUser))
	assert.False(t, UserRepo.IsSecretCorrect(sampleUser, secret))
}
