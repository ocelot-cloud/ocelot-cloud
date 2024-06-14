package src

func TestHub() {
	defer Cleanup()
	deleteArtifacts()
	StartDaemon(hubDir, "go run .")
	WaitUntilPortIsReady("localhost:8082")
	ExecuteInDir(hubDir, "bash test.sh")
	deleteArtifacts()
}

func deleteArtifacts() {
	ExecuteInDir(hubDir, "bash -c 'rm -rf users *.tar.gz'")
}
