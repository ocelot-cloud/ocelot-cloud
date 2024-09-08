package apps_new

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

// TODO Execute tests for both, real docker service and mock.
func TestStateTransitions(t *testing.T) {
	app := createAppState()
	assert.Equal(t, Uninitialized, app.state)
	app.Deploy()
	assert.Equal(t, Downloading, app.state)
	waitAndAssert(t, app.state, Starting)
	/* TODO
	waitAndAssert(t, app.state, Available)
	app.Stop()
	assert.Equal(t, Stopping, app.state)
	waitAndAssert(t, app.state, Uninitialized)
	*/
}

func waitAndAssert(t *testing.T, state State, uninitialized State) {
	// TODO
}
