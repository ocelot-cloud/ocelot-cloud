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
	assert.NotNil(t, createAdminUserIfNotExistent("", ""))
	assert.NotNil(t, createAdminUserIfNotExistent("admin", ""))
	assert.NotNil(t, createAdminUserIfNotExistent("", "password"))
}

func TestAdminInitializationWithCorrectEnvs(t *testing.T) {
	assert.Nil(t, createAdminUserIfNotExistent("admin", "password"))
}

func TestAdminInitializationIsIgnoredWhenAlreadyExistsInDatabase(t *testing.T) {
	err := repo.CreateUser("admin", "password", true)
	assert.Nil(t, err)
	assert.Nil(t, createAdminUserIfNotExistent("", ""))
}
