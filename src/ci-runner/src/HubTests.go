package src

func TestHub() {
	TestFilesystemManager()
	TestUploadAndDownload()
}

func TestUploadAndDownload() {
	printTestDescription("Testing file upload and download")
	defer Cleanup()
	deleteArtifacts()
	StartDaemon(hubDir, "go run .")
	WaitUntilPortIsReady("localhost:8082")
	ExecuteInDir(hubDir, "go test -run TestFileUploadDownload")
	deleteArtifacts()
}

func TestFilesystemManager() {
	printTestDescription("Testing filesystem manager")
	ExecuteInDir(hubDir, "go test -run TestFilesystemManager")
}

func deleteArtifacts() {
	ExecuteInDir(hubDir, "bash -c 'rm -rf users *.tar.gz'")
}
