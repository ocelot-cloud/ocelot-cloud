package apps_new

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
	"time"
)

// TODO Execute tests for both, real docker service and mock.
func TestStateTransitions(t *testing.T) {
	app := createAppState()
	assert.Equal(t, Uninitialized, app.state)
	app.Deploy()
	assert.Equal(t, Downloading, app.state)
	waitAndAssert(t, Starting, app)
	waitAndAssert(t, Available, app)
	app.Stop()
	assert.Equal(t, Stopping, app.state)
	waitAndAssert(t, Uninitialized, app)
}

func waitAndAssert(t *testing.T, expected State, app *AppContext) {
	maxAttempts := 20
	for i := 0; i < maxAttempts; i++ {
		if app.state == expected {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fail()
}
