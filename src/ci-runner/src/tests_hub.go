package src

func TestHubAll() {
	ExecuteInDir(hubDir, "rm -rf data")
	TestHubUnits()
	TestHubBackend()
	TestHubPersistence()
	TestHubAcceptance()
}

func TestHubUnits() {
	printTaskDescription("Testing hub units")
	defer Cleanup()
	ExecuteInDir(hubDir, "go test -tags=unit ./...", "LOG_LEVEL=DEBUG")
}

func TestHubBackend() {
	printTaskDescription("Testing hub backend")
	defer Cleanup()
	StartDaemon(hubDir, "go run .", "PROFILE=TEST")
	WaitUntilPortIsReady("localhost:8082")
	ExecuteInDir(hubDir, "go test -tags=acceptance ./...")
}

func TestHubAcceptance() {
	printTaskDescription("Testing hub backend")
	defer Cleanup()
	ExecuteInDir(hubDir, "rm -rf data")
	StartDaemon(hubDir, "bash run-development-setup.sh")
	WaitUntilPortIsReady("localhost:8082")
	Build(Frontend)
	StartDaemon(frontendDir, "npm run serve", "VITE_APP_PROFILE="+TestProfile)
	WaitForIndexPageToBeReady(frontendServerUrl)
	ExecuteInDir(acceptanceTestsDir, "npx cypress run --spec cypress/e2e/hub.cy.ts --headless")
}

func TestHubPersistence() {
	printTaskDescription("Testing hub persistence")
	defer Cleanup()
	ExecuteInDir(hubDir, "rm -rf data")
	ExecuteInDir(hubDir, "go build")
	StartDaemon(hubDir, "./hub")
	WaitUntilPortIsReady("localhost:8082")
	ExecuteInDir(hubDir, "[ -f ./data/sqlite.db ]")
	ExecuteInDir(hubDir, "[ -f ./data/logs.txt ]")
}
