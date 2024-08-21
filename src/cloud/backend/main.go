package main

import (
	"github.com/gorilla/mux"
	"github.com/ocelot-cloud/shared"
	"ocelot/backend/business"
	"ocelot/backend/config"
	"ocelot/backend/security"
	"os/exec"
	"strings"
)

// TODO Make tests with real containers only when using REST API, for GUI/Acceptance tests always use the mock since it makes trouble in CI otherwise
// TODO Simplify profiles: DEV + PROD, no mocked frontend anymore, no security disabling anymore.
// TODO Due to implementation of the hub I can delete alls the stacks in the cloud. Acceptance tests need to integrate hub and need to implement download of stacks at the beginning?
// TODO In the end, add deploy script which only works on my device, since I have the correct SSH keys and config.

var logger = shared.ProvideLogger()

func main() {
	verifyCliToolInstallations()
	config := tools.GenerateGlobalConfiguration()
	router := mux.NewRouter()
	securityModule := security.ProvideSecurityModule(router, config)
	businessModule := business.ProvideBusinessModule(router, config, securityModule)
	businessModule.InitializeApplication()
}

func verifyCliToolInstallations() {
	cliTools := []string{
		"sqlite3 --version",
		"docker version",
		"docker compose version",
	}

	for _, fullCmd := range cliTools {
		parts := strings.Split(fullCmd, " ")
		toolName := parts[0]
		cmdArgs := parts[1:]

		crashIfToolIsNotInstalled(toolName, cmdArgs)
	}
	logger.Info("All required CLI tools seem to be installed.")
}

func crashIfToolIsNotInstalled(toolName string, args []string) {
	cmd := exec.Command(toolName, args...)
	if err := cmd.Run(); err != nil {
		logger.Fatal("Error, tried command '%s %s' but CLI tool seems not to be installed properly.", toolName, strings.Join(args, " "))
	}
}
