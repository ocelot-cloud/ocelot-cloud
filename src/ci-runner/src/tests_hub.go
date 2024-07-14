package src

func TestHub() {
	ExecuteInDir(hubDir, "rm -rf data")
	testHubUnits()
	testHubAcceptance()
}

func testHubUnits() {
	printTestDescription("Testing hub units")
	defer Cleanup()
	ExecuteInDir(hubDir, "go test -tags=unit ./...")
}

func testHubAcceptance() {
	printTestDescription("Testing hub acceptance")
	defer Cleanup()
	StartDaemon(hubDir, "go run .", "PROFILE=TEST")
	WaitUntilPortIsReady("localhost:8082")
	ExecuteInDir(hubDir, "go test -tags=acceptance ./...")
}
