package src

func TestHub() {
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
	StartDaemon(hubDir, "go run .")
	WaitUntilPortIsReady("localhost:8082")
	ExecuteInDir(hubDir, "go test -tags=acceptance ./...")
}
