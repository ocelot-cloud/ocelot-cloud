package src

import "ocelot/ci-runner/cli"

func TestHubAll() {
	cli.ExecuteInDir(hubDir, "rm -rf data")
	TestHubUnits()
	TestHubBackend()
	TestHubPersistence()
	TestHubAcceptance()
}

func TestHubUnits() {
	printTaskDescription("Testing hub units")
	defer cli.Cleanup()
	cli.ExecuteInDir(hubDir, "go test -tags=unit ./...", "LOG_LEVEL=DEBUG")
}

func TestHubBackend() {
	printTaskDescription("Testing hub backend")
	defer cli.Cleanup()
	cli.StartDaemon(hubDir, "go run .", "PROFILE=TEST")
	cli.WaitUntilPortIsReady("8082")
	cli.ExecuteInDir(hubDir, "go test -tags=acceptance ./...")
}

func TestHubAcceptance() {
	printTaskDescription("Testing hub backend")
	defer cli.Cleanup()
	cli.ExecuteInDir(hubDir, "rm -rf data")
	cli.StartDaemon(hubDir, "bash run-development-setup.sh")
	cli.WaitUntilPortIsReady("8082")
	Build(Frontend)
	cli.StartDaemon(frontendDir, "npm run serve", "VITE_APP_PROFILE="+TestProfile)
	cli.WaitForIndexPageToBeReady(frontendServerUrl)
	cli.ExecuteInDir(acceptanceTestsDir, "npx cypress run --spec cypress/e2e/hub.cy.ts --headless")
}

func TestHubPersistence() {
	printTaskDescription("Testing hub persistence")
	defer cli.Cleanup()
	cli.ExecuteInDir(hubDir, "rm -rf data")
	cli.ExecuteInDir(hubDir, "go build")
	cli.StartDaemon(hubDir, "./hub")
	cli.WaitUntilPortIsReady("8082")
	cli.ExecuteInDir(hubDir, "[ -f ./data/sqlite.db ]")
	cli.ExecuteInDir(hubDir, "[ -f ./data/logs.txt ]")
}
