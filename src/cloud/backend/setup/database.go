package setup

import (
	"fmt"
	"ocelot/backend/security"
	"ocelot/backend/tools"
)

const (
	initialAdminNameEnv     = "INITIAL_ADMIN_NAME"
	initialAdminPasswordEnv = "INITIAL_ADMIN_PASSWORD"
)

var repo = security.Repo

// TODO Maybe put that stuff in the security module? Also this isn't just security, but also other stuff. Maybe create a "repository" package?
func InitializeDatabase(config *tools.GlobalConfig) {
	if config.UseRealDatabase {
		security.InitializeDatabaseWithSource(security.DatabaseFile)
	} else {
		security.InitializeDatabaseWithSource(":memory:")
	}

	/* TODO Uncomment and adapt tests/handlers
	err := createAdminUserIfNotExistent(os.Getenv(initialAdminNameEnv), os.Getenv(initialAdminPasswordEnv))
	if err != nil {
		logger.Fatal("Admin user initialization failed: %v", err)
	}
	*/
}

// TODO Add Input validation to env credentials
func createAdminUserIfNotExistent(adminNameEnv string, adminPasswordEnv string) error {
	if repo.DoesAnyAdminUserExist() {
		logger.Info("There is at least one admin user in the database, so env admin initialization env variables will be ignored")
		return nil
	} else {
		logger.Info("Application needs at least one admin user, but none was found in database. Trying to create the admin user from env variables.")
		if adminNameEnv == "" {
			return fmt.Errorf("necessary env variable '%s' is not set", initialAdminNameEnv)
		} else if adminPasswordEnv == "" {
			return fmt.Errorf("necessary env variable '%s' is not set", initialAdminPasswordEnv)
		} else {
			err := repo.CreateUser(adminPasswordEnv, adminPasswordEnv, true)
			if err != nil {
				return fmt.Errorf("initial admin user creation from env variables failed: %v", err)
			}
			return nil
		}
	}
}
