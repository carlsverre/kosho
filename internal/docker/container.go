package docker

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func StartInteractiveShell(worktreePath string) error {
	fmt.Printf("Starting interactive bash shell in %s\n", worktreePath)

	// Change to the worktree directory
	err := os.Chdir(worktreePath)
	if err != nil {
		return fmt.Errorf("failed to change to worktree directory: %w", err)
	}

	// Execute bash shell
	cmd := exec.Command("bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		return fmt.Errorf("failed to run bash shell: %w", err)
	}

	return nil
}

func StartDetachedShell(worktreePath string) error {
	fmt.Printf("Starting detached bash shell in %s\n", worktreePath)
	fmt.Printf("Note: Detached mode is stubbed - would run container in background\n")
	return nil
}

// Stubbed Docker functions for future implementation

func CreateContainer(repoName, worktreeName, worktreePath string) error {
	configVolume := fmt.Sprintf("%s-%s-config", repoName, worktreeName)
	historyVolume := fmt.Sprintf("%s-%s-history", repoName, worktreeName)

	fmt.Printf("Would create Docker container with:\n")
	fmt.Printf("  Config volume: %s\n", configVolume)
	fmt.Printf("  History volume: %s\n", historyVolume)
	fmt.Printf("  Workspace mount: %s:/workspace\n", worktreePath)
	fmt.Printf("  Image: kosho-runtime\n")
	fmt.Printf("  Capabilities: NET_ADMIN, NET_RAW\n")

	return nil
}

func StartContainer(containerID string) error {
	fmt.Printf("Would start container: %s\n", containerID)
	return nil
}

func StopContainer(containerID string) error {
	fmt.Printf("Would stop container: %s\n", containerID)
	return nil
}

func CreateNamedVolume(volumeName string) error {
	fmt.Printf("Would create named volume: %s\n", volumeName)
	return nil
}
