package src

func TestHub() {
	ExecuteInDir(hubDir, "rm -rf data")
	TestHubUnits()
	TestHubBackend()
}

func TestHubUnits() {
	printTestDescription("Testing hub units")
	defer Cleanup()
	ExecuteInDir(hubDir, "go test -tags=unit ./...", "LOG_LEVEL=DEBUG")
}

func TestHubBackend() {
	printTestDescription("Testing hub backend")
	defer Cleanup()
	StartDaemon(hubDir, "go run .", "PROFILE=TEST", "LOG_LEVEL=DEBUG")
	WaitUntilPortIsReady("localhost:8082")
	ExecuteInDir(hubDir, "go test -tags=acceptance ./...")
}

func TestHubAcceptance() {
	printTestDescription("Testing hub backend")
	defer Cleanup()
	ExecuteInDir(hubDir, "rm -rf data")
	StartDaemon(hubDir, "go run .")
	WaitUntilPortIsReady("localhost:8082")
	StartDaemon(frontendDir, "bash run-development-setup.sh")
	WaitUntilPortIsReady("localhost:8081")
	ExecuteInDir(acceptanceTestsDir, "npx cypress run --spec cypress/e2e/hub.cy.ts --headless")
}

func TestHubAll() {
	TestHubBackend()
	TestHubAcceptance()
}
