package src

import (
	"os"
	"path/filepath"
)

const Scheme = "http"
const RootDomain = "localhost"
const ocelotUrl = Scheme + "://ocelot-cloud." + RootDomain
const frontendServerUrl = Scheme + "://localhost:8081"

const ocelotContainerRunCommandDetached = "docker compose -p ocelot-cloud up -d"
const cypressCommand = "npx cypress run --spec cypress/e2e/cloud.cy.ts --headless"

var projectDir = GetProjectDir()
var scriptsDir = projectDir + "/scripts"
var srcDir = projectDir + "/src"

var backendDir = srcDir + "/backend"
var frontendDir = srcDir + "/frontend"
var acceptanceTestsDir = srcDir + "/acceptance-tests"

var ocelotStackDir = backendDir + "/assets/ocelot-cloud"
var backendAppsDir = backendDir + "/apps"
var backendRepoDir = backendDir + "/repo"
var backendToolsDir = backendDir + "/tools"
var cloudHubClientDir = backendDir + "/apps_new"
var backendComponentTestsDir = backendDir + "/component-tests"

var hubBackendDir = srcDir + "/hub/backend"

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
