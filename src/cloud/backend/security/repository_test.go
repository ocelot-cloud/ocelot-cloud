package security

import (
	"github.com/ocelot-cloud/shared/assert"
	"os"
	"testing"
)

var repo Repository = &MyRepository{}

func TestMain(m *testing.M) {
	initializeDatabaseWithSource(":memory:")
	repo.WipeDatabase()
	code := m.Run()
	os.Exit(code)
}

const (
	sampleUser     = "user"
	samplePassword = "password"
)

// TODO Finish SQLite Client Implementation And Tests
func TestSqliteClient(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleUser, samplePassword, true))
	assert.True(t, repo.IsPasswordCorrect(sampleUser, samplePassword))
}
