package src

func TestHub() {
	printTestDescription("Testing file upload and download")
	defer Cleanup()
	StartDaemon(hubDir, "go run .")
	WaitUntilPortIsReady("localhost:8082")
	ExecuteInDir(hubDir, "go test .")
}
