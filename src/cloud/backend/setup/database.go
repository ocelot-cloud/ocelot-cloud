package setup

import (
	"ocelot/backend/security"
	"ocelot/backend/tools"
	"os"
)

// TODO Maybe put that stuff in the security module? Also this isn't just security, but also other stuff. Maybe create a "repository" package?
func InitializeDatabase(config *tools.GlobalConfig) {
	if config.UseRealDatabase {
		security.InitializeDatabaseWithSource(security.DatabaseFile)
	} else {
		security.InitializeDatabaseWithSource(":memory:")
	}

	err := createAdminUserIfNotExistent(os.Getenv("INITIAL_ADMIN_NAME"), os.Getenv("INITIAL_ADMIN_PASSWORD"))
	if err != nil {
		logger.Fatal("Admin user initialization failed: %v", err)
	}
}

func createAdminUserIfNotExistent(adminNameEnv string, adminPasswordEnv string) error {
	return nil
}
