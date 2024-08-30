package src

import (
	"os"
	"path/filepath"
)

const ocelotContainerRunCommand = "docker compose -p ocelot-cloud up"
const ocelotContainerRunCommandDetached = "docker compose -p ocelot-cloud up -d"
const cypressCommand = "npx cypress run --spec cypress/e2e/cloud.cy.ts --headless"

var projectDir = GetProjectDir()
var scriptsDir = projectDir + "/scripts"
var srcDir = projectDir + "/src"
var cloudDir = srcDir + "/cloud"
var backendDir = cloudDir + "/backend"
var backendComponentTestsDir = backendDir + "/component-tests"
var frontendDir = cloudDir + "/frontend"
var acceptanceTestsDir = cloudDir + "/acceptance-tests"
var ocelotStackDir = backendDir + "/stacks/ocelot-cloud"
var backendAppsDir = backendDir + "/apps"
var backendSecurityDir = backendDir + "/security"
var backendToolsDir = backendDir + "/tools"

var hubDir = srcDir + "/hub"

// TestProfile There is also the "PROD" profile, but it should be used automatically if no profile is given.
const TestProfile = "TEST"

func GetProjectDir() string {
	devopsRunnerDir, _ := os.Getwd()
	src := filepath.Dir(devopsRunnerDir)
	return filepath.Dir(src)
}

func getTestProfileEnv() string {
	return "PROFILE=" + TestProfile
}

func getEnableDummyStacksEnv(enabled bool) string {
	prefix := "USE_DUMMY_STACKS="
	if enabled {
		return prefix + "true"
	} else {
		return prefix + "false"
	}
}
