package internal

import (
	"os/exec"
	"runtime"
)

// ShellCommand returns the shell command and argument for executing shell commands
// cross-platform. On Unix-like systems uses exec.LookPath to find "sh" and returns "-c", 
// on Windows returns "cmd.exe" and "/C".
func ShellCommand() (string, string) {
	if runtime.GOOS == "windows" {
		return "cmd.exe", "/C"
	}
	
	shellPath, err := exec.LookPath("sh")
	if err != nil {
		// Fallback to /bin/sh if LookPath fails
		shellPath = "/bin/sh"
	}
	return shellPath, "-c"
}
