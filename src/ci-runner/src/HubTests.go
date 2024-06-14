package src

func TestHub() {
	defer Cleanup()
	deleteArtifacts()
	StartDaemon(hubDir, "go run .")
	WaitUntilPortIsReady("localhost:8082")
	ExecuteInDir(hubDir, "go test ./...")
	deleteArtifacts()
}

func deleteArtifacts() {
	ExecuteInDir(hubDir, "bash -c 'rm -rf users *.tar.gz'")
}
