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

// TODO Make CI pipeline running again
// TODO Make tests with real containers only when using REST API, for GUI/Acceptance tests always use the mock since it makes trouble in CI otherwise

// TODO Update "shared" module version
// TODO Consider reusing stuff from the hub, like security (potential clash with cloud package "security"), sql logic, hub client (search for apps, download, maybe upload to keep them private?)
// TODO Implement security, there should be a policy that Origin from request header == initially defined Origin as ENV variable or default ("http://localhost:8080")
// TODO refactor table: list apps with state, but make them selectable, so that there is only a single start/stop button.
// TODO Simplify profiles: DEV + PROD, no mocked frontend anymore, no security disabling anymore.
// TODO Due to implementation of the hub I can delete alls the stacks in the cloud. Acceptance tests need to integrate hub and need to implement download of stacks at the beginning? Hub should have those default files included? -> Dummies stay in cloud, sample apps like gitea go to the hub
// TODO In the end, add deploy script which only works on my device, since I have the correct SSH keys and config.
// TODO Drop the folder structure for the stacks and store everything in an sqlite. When using dummies, just load them into database at the start if not present.

var logger = shared.ProvideLogger()

func main() {
	// TODO shouldn't here be all the "read CLI args/env variable/profiles/config/loglevel/db setup stuff" step?
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
