package repo

import (
	"fmt"
	"ocelot/backend/tools"
	"os"
)

const (
	initialAdminNameEnv     = "INITIAL_ADMIN_NAME"
	initialAdminPasswordEnv = "INITIAL_ADMIN_PASSWORD"
)

// TODO Maybe put that stuff in the security module? Also this isn't just security, but also other stuff. Maybe create a "repository" package?
func InitializeDatabase() {
	if tools.Config.UseRealDatabase {
		InitializeDatabaseWithSource(DatabaseFile)
	} else {
		InitializeDatabaseWithSource(":memory:")
	}

	err := createAdminUserIfNotExistent(os.Getenv(initialAdminNameEnv), os.Getenv(initialAdminPasswordEnv), tools.Config.CreateDefaultAdminUser)
	if err != nil {
		Logger.Fatal("Admin user initialization failed: %v", err)
	}
}

// TODO Add Input validation to env credentials
func createAdminUserIfNotExistent(adminNameEnv string, adminPasswordEnv string, createDefaultAdminUser bool) error {
	// TODO That means I can remove the ENV variable from the TEST profile backend start in ci-runner
	if createDefaultAdminUser {
		return UserRepo.CreateUser("admin", "password", true)
	}

	if UserRepo.DoesAnyAdminUserExist() {
		Logger.Info("There is at least one admin user in the database, so admin initialization via env variables will not be conducted.")
		return nil
	} else {
		Logger.Info("Application needs at least one admin user, but none was found in database. Trying to create the admin user from env variables.")
		return createAdminsUserFromEnvs(adminNameEnv, adminPasswordEnv)
	}
}

func createAdminsUserFromEnvs(adminNameEnv string, adminPasswordEnv string) error {
	if adminNameEnv == "" {
		return fmt.Errorf("necessary env variable '%s' is not set", initialAdminNameEnv)
	} else if adminPasswordEnv == "" {
		return fmt.Errorf("necessary env variable '%s' is not set", initialAdminPasswordEnv)
	} else {
		err := UserRepo.CreateUser(adminNameEnv, adminPasswordEnv, true)
		if err != nil {
			return fmt.Errorf("initial admin user creation from env variables failed: %v", err)
		}
		Logger.Info("Initial admin user '%s' created", adminNameEnv)
		return nil
	}
}
