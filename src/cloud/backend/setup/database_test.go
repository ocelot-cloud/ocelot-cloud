package setup

import "testing"

func TestAsdf(t *testing.T) {
	// TODO Check if admin user exists. If not take it from the env variables. If not existent, crash.
	// TODO Add tests: 1) neither admin in repo nor in envs -> crash, 2) no admin in repo, but in envs -> no crash, 3) admin in repo, but not in envs -> no crash

	// TODO createAdminUserIfNotExistent()
}
