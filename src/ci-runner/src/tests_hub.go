package src

import "github.com/ocelot-cloud/task-runner"

func TestHubAll() {
	tr.ExecuteInDir(hubBackendDir, "rm -rf data")
	TestHubUnits()
	TestHubBackend()
	TestHubPersistence()
	TestHubAcceptance()
}

func TestHubUnits() {
	tr.PrintTaskDescription("Testing hub units")
	defer tr.Cleanup()
	tr.ExecuteInDir(hubBackendDir, "go test -tags=unit ./...", "LOG_LEVEL=DEBUG")
}

func TestHubBackend() {
	tr.PrintTaskDescription("Testing hub backend")
	defer tr.Cleanup()
	tr.StartDaemon(hubBackendDir, "go run .", "PROFILE=TEST")
	tr.WaitUntilPortIsReady("8082")
	tr.ExecuteInDir(hubBackendDir, "go test -tags=acceptance ./...")
}

func TestHubAcceptance() {
	tr.PrintTaskDescription("Testing hub backend")
	defer tr.Cleanup()
	tr.ExecuteInDir(hubBackendDir, "rm -rf data")
	tr.StartDaemon(hubBackendDir, "bash run-development-setup.sh")
	tr.WaitUntilPortIsReady("8082")
	Build(Frontend)
	tr.StartDaemon(frontendDir, "npm run serve", "VITE_APP_PROFILE="+TestProfile)
	tr.WaitForWebPageToBeReady(frontendServerUrl)
	tr.ExecuteInDir(acceptanceTestsDir, "npx cypress run --spec cypress/e2e/hub.cy.ts --headless")
}

func TestHubPersistence() {
	tr.PrintTaskDescription("Testing hub persistence")
	defer tr.Cleanup()
	tr.ExecuteInDir(hubBackendDir, "rm -rf data")
	tr.ExecuteInDir(hubBackendDir, "go build")
	tr.StartDaemon(hubBackendDir, "./hub")
	tr.WaitUntilPortIsReady("8082")
	tr.ExecuteInDir(hubBackendDir, "[ -f ./data/sqlite.db ]")
	tr.ExecuteInDir(hubBackendDir, "[ -f ./data/logs.txt ]")
}
