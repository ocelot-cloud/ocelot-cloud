package setup

import (
	"ocelot/backend/security"
	"ocelot/backend/tools"
)

// TODO Maybe put that stuff in the security module? Also this isn't just security, but also other stuff. Maybe create a "repository" package?
func InitializeDatabase(config *tools.GlobalConfig) {
	if config.UseRealDatabase {
		security.InitializeDatabaseWithSource(security.DatabaseFile)
	} else {
		security.InitializeDatabaseWithSource(":memory:")
	}

	err := createAdminUserIfNotExistent()
	if err != nil {
		logger.Fatal("Admin user initialization failed: %v", err)
	}
}

func createAdminUserIfNotExistent() error {
	// TODO Check if admin user exists. If not take it from the env variables. If not existent, crash.
	// TODO Add tests: 1) neither admin in repo nor in envs -> crash, 2) no admin in repo, but in envs -> no crash, 3) admin in repo, but not in envs -> no crash
	return nil
}
