package setup

import (
	"fmt"
	"ocelot/backend/security"
	"ocelot/backend/tools"
	"os"
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

	err := createAdminUserIfNotExistent(os.Getenv(initialAdminNameEnv), os.Getenv(initialAdminPasswordEnv), config.CreateDefaultAdminUser)
	if err != nil {
		logger.Fatal("Admin user initialization failed: %v", err)
	}
}

// TODO Add Input validation to env credentials
func createAdminUserIfNotExistent(adminNameEnv string, adminPasswordEnv string, createDefaultAdminUser bool) error {
	// TODO That means I can remove the ENV variable from the TEST profile backend start in ci-runner
	if createDefaultAdminUser {
		return repo.CreateUser("admin", "password", true)
	}

	if repo.DoesAnyAdminUserExist() {
		logger.Info("There is at least one admin user in the database, so admin initialization via env variables will not be conducted.")
		return nil
	} else {
		logger.Info("Application needs at least one admin user, but none was found in database. Trying to create the admin user from env variables.")
		return createAdminsUserFromEnvs(adminNameEnv, adminPasswordEnv)
	}
}

func createAdminsUserFromEnvs(adminNameEnv string, adminPasswordEnv string) error {
	if adminNameEnv == "" {
		return fmt.Errorf("necessary env variable '%s' is not set", initialAdminNameEnv)
	} else if adminPasswordEnv == "" {
		return fmt.Errorf("necessary env variable '%s' is not set", initialAdminPasswordEnv)
	} else {
		err := repo.CreateUser(adminNameEnv, adminPasswordEnv, true)
		if err != nil {
			return fmt.Errorf("initial admin user creation from env variables failed: %v", err)
		}
		logger.Info("Initial admin user '%s' created", adminNameEnv)
		return nil
	}
}
