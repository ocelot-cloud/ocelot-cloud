package src

// TODO add to CI runner CLI interface
func TestHub() {
	defer Cleanup()
	ExecuteInDir(hubDir, "rm -rf users")
	StartDaemon(hubDir, "go run .")
	ExecuteInDir(hubDir, "bash test.sh")
}
