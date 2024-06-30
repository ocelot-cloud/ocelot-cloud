package src

func TestHub() {
	printTestDescription("Testing file upload and download")
	defer Cleanup()
	StartDaemon(hubDir, "go run .")
	WaitUntilPortIsReady("localhost:8082")
	ExecuteInDir(hubDir, "go test .")
}

// TODO Write tests for CI-Runner? Especially when It should fail, e.g. to no were tests found or so.
// TODO for this command: "go test -run TestFilesystemManager,..." comes the output -> "testing: warning: no tests to run", which should immediately fail.
