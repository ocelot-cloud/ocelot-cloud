package setup

import (
	"github.com/ocelot-cloud/shared/assert"
	"ocelot/backend/security"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	security.InitializeDatabaseWithSource(":memory:")
	repo.WipeDatabase()
	code := m.Run()
	os.Exit(code)
}

func TestEmptyAdminInitializationEnvsShouldFail(t *testing.T) {
	assert.NotNil(t, createAdminUserIfNotExistent("", "", false))
	assert.NotNil(t, createAdminUserIfNotExistent("admin", "", false))
	assert.NotNil(t, createAdminUserIfNotExistent("", "password", false))
}

func TestAdminInitializationWithCorrectEnvs(t *testing.T) {
	assert.Nil(t, createAdminUserIfNotExistent("admin", "password", false))
}

func TestAdminInitializationIsIgnoredWhenAlreadyExistsInDatabase(t *testing.T) {
	err := repo.CreateUser("admin", "password", true)
	assert.Nil(t, err)
	assert.Nil(t, createAdminUserIfNotExistent("", "", false))
}

func TestDefaultAdminCreation(t *testing.T) {
	assert.False(t, repo.DoesAnyAdminUserExist())
	assert.Nil(t, createAdminUserIfNotExistent("", "", true))
	assert.True(t, repo.DoesAnyAdminUserExist())
	assert.True(t, repo.IsPasswordCorrect("admin", "password"))
	// TODO Also assert isAdmin == true
}
