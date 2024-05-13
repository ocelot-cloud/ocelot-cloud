package internal

import (
	"os"
	"testing"
)

// TODO Finish SQLite Client Implementation
func TestSqliteClient(t *testing.T) {
	DoSomeDataBaseStuff()
	err := os.Remove("ocelot-cloud.db")
	if err != nil && !os.IsNotExist(err) {
		t.Error("Error deleting ocelot-cloud.db:", err)
	}
}
