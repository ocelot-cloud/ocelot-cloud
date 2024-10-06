package repo

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
	"time"
)

func TestEmptyAdminInitializationEnvsShouldFail(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.NotNil(t, createAdminUserIfNotExistent("", "", false))
	assert.NotNil(t, createAdminUserIfNotExistent("admin", "", false))
	assert.NotNil(t, createAdminUserIfNotExistent("", "password", false))
}

func TestAdminInitializationWithCorrectEnvs(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, createAdminUserIfNotExistent("admin", "password", false))
}

func TestAdminInitializationIsIgnoredWhenAlreadyExistsInDatabase(t *testing.T) {
	defer dbRepo.WipeDatabase()
	err := UserRepo.CreateUser("admin", "password", true)
	assert.Nil(t, err)
	assert.Nil(t, createAdminUserIfNotExistent("", "", false))
}

func TestDefaultAdminCreation(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.False(t, UserRepo.DoesAnyAdminUserExist())
	assert.Nil(t, createAdminUserIfNotExistent("", "", true))
	assert.True(t, UserRepo.DoesAnyAdminUserExist())
	assert.True(t, UserRepo.IsPasswordCorrect("admin", "password"))

	assert.Nil(t, UserRepo.HashAndSaveCookie("admin", "some-cookie", time.Now()))
	auth, err := UserRepo.GetUserViaCookie("some-cookie")
	assert.Nil(t, err)
	assert.Equal(t, "admin", auth.User)
	assert.True(t, auth.IsAdmin)
}
