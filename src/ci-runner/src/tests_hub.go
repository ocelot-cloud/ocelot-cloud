package src

import "github.com/ocelot-cloud/task-runner"

func TestHubAll() {
	tr.ExecuteInDir(hubDir, "rm -rf data")
	TestHubUnits()
	TestHubBackend()
	TestHubPersistence()
	TestHubAcceptance()
}

func TestHubUnits() {
	tr.PrintTaskDescription("Testing hub units")
	defer tr.Cleanup()
	tr.ExecuteInDir(hubDir, "go test -tags=unit ./...", "LOG_LEVEL=DEBUG")
}

func TestHubBackend() {
	tr.PrintTaskDescription("Testing hub backend")
	defer tr.Cleanup()
	tr.StartDaemon(hubDir, "go run .", "PROFILE=TEST")
	tr.WaitUntilPortIsReady("8082")
	tr.ExecuteInDir(hubDir, "go test -tags=acceptance ./...")
}

func TestHubAcceptance() {
	tr.PrintTaskDescription("Testing hub backend")
	defer tr.Cleanup()
	tr.ExecuteInDir(hubDir, "rm -rf data")
	tr.StartDaemon(hubDir, "bash run-development-setup.sh")
	tr.WaitUntilPortIsReady("8082")
	Build(Frontend)
	tr.StartDaemon(frontendDir, "npm run serve", "VITE_APP_PROFILE="+TestProfile)
	tr.WaitForWebPageToBeReady(frontendServerUrl)
	tr.ExecuteInDir(acceptanceTestsDir, "npx cypress run --spec cypress/e2e/hub.cy.ts --headless")
}

func TestHubPersistence() {
	tr.PrintTaskDescription("Testing hub persistence")
	defer tr.Cleanup()
	tr.ExecuteInDir(hubDir, "rm -rf data")
	tr.ExecuteInDir(hubDir, "go build")
	tr.StartDaemon(hubDir, "./hub")
	tr.WaitUntilPortIsReady("8082")
	tr.ExecuteInDir(hubDir, "[ -f ./data/sqlite.db ]")
	tr.ExecuteInDir(hubDir, "[ -f ./data/logs.txt ]")
}
