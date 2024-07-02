package internal

import (
	"os"
	"testing"
)

// TODO Finish SQLite Client Implementation
func TestSqliteClient(t *testing.T) {
	DoSomeDataBaseStuff()
	err := os.Remove(databaseFile)
	if err != nil && !os.IsNotExist(err) {
		t.Errorf("Error deleting %s: %v", databaseFile, err)
	}
}
