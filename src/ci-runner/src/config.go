package src

import (
	"os"
	"path/filepath"
)

const Scheme = "http"
const RootDomain = "localhost"
const ocelotUrl = Scheme + "://ocelot-cloud." + RootDomain
const frontendServerUrl = Scheme + "://localhost:8081"

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
var ocelotStackDir = backendDir + "/assets/ocelot-cloud"
var backendAppsDir = backendDir + "/apps"
var backendSecurityDir = backendDir + "/security"
var backendToolsDir = backendDir + "/tools"

var hubDir = srcDir + "/hub"

const TestProfile = "TEST"
const ProdProfile = "PROD"

func GetProjectDir() string {
	devopsRunnerDir, _ := os.Getwd()
	src := filepath.Dir(devopsRunnerDir)
	return filepath.Dir(src)
}

func getTestProfileEnv() string {
	return "PROFILE=" + TestProfile
}

func getProdProfileEnv() string {
	return "PROFILE=" + ProdProfile
}

func getEnableDummyStacksEnv(enabled bool) string {
	prefix := "USE_DUMMY_STACKS="
	if enabled {
		return prefix + "true"
	} else {
		return prefix + "false"
	}
}
