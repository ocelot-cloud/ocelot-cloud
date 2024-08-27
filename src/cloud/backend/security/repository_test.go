package security

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	initializeDatabaseWithSource(":memory:")
	code := m.Run()
	os.Exit(code)
}

// TODO Finish SQLite Client Implementation
func TestSqliteClient(t *testing.T) {
	println("Hello!")
}
