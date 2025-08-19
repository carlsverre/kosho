package internal

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRunInitHooks_EmptyCommands(t *testing.T) {
	tempDir, cleanup := setupTempWorktree(t)
	defer cleanup()

	kw := NewKoshoWorktree(tempDir, "test-worktree")

	// Test with nil slice
	err := kw.RunInitHooks(nil)
	if err != nil {
		t.Errorf("RunInitHooks with nil commands should return nil, got: %v", err)
	}

	// Test with empty slice
	err = kw.RunInitHooks([]string{})
	if err != nil {
		t.Errorf("RunInitHooks with empty commands should return nil, got: %v", err)
	}

	// Test with slice containing only empty strings
	err = kw.RunInitHooks([]string{"", "  ", "\t"})
	if err != nil {
		t.Errorf("RunInitHooks with only empty/whitespace commands should return nil, got: %v", err)
	}
}

func TestRunInitHooks_ValidCommands(t *testing.T) {
	tempDir, cleanup := setupTempWorktree(t)
	defer cleanup()

	kw := NewKoshoWorktree(tempDir, "test-worktree")

	// Test with echo command
	err := kw.RunInitHooks([]string{"echo 'test message'"})
	if err != nil {
		t.Errorf("RunInitHooks with valid echo command should succeed, got: %v", err)
	}

	// Test with touch command to create a file
	testFile := filepath.Join(kw.WorktreePath(), "test-file.txt")
	var touchCmd string
	if runtime.GOOS == "windows" {
		// On Windows, use type nul > filename to create an empty file
		touchCmd = "type nul > test-file.txt"
	} else {
		touchCmd = "touch test-file.txt"
	}

	err = kw.RunInitHooks([]string{touchCmd})
	if err != nil {
		t.Errorf("RunInitHooks with touch command should succeed, got: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("Expected file %s to be created by touch command", testFile)
	}
}

func TestRunInitHooks_FailingCommand(t *testing.T) {
	tempDir, cleanup := setupTempWorktree(t)
	defer cleanup()

	kw := NewKoshoWorktree(tempDir, "test-worktree")

	// Test with command that always fails
	var failCmd string
	if runtime.GOOS == "windows" {
		failCmd = "exit 1"
	} else {
		failCmd = "false"
	}

	err := kw.RunInitHooks([]string{failCmd})
	if err == nil {
		t.Error("RunInitHooks with failing command should return an error")
	}

	// Check that error message contains the command
	if !strings.Contains(err.Error(), "init hook") || !strings.Contains(err.Error(), "failed") {
		t.Errorf("Error message should mention init hook failure, got: %v", err)
	}
}

func TestRunInitHooks_MultipleCommands_StopOnFailure(t *testing.T) {
	tempDir, cleanup := setupTempWorktree(t)
	defer cleanup()

	kw := NewKoshoWorktree(tempDir, "test-worktree")

	// Create commands: touch file1, fail, touch file2
	file1 := filepath.Join(kw.WorktreePath(), "file1.txt")
	file2 := filepath.Join(kw.WorktreePath(), "file2.txt")

	var touchCmd1, failCmd, touchCmd2 string
	if runtime.GOOS == "windows" {
		touchCmd1 = "type nul > file1.txt"
		failCmd = "exit 1"
		touchCmd2 = "type nul > file2.txt"
	} else {
		touchCmd1 = "touch file1.txt"
		failCmd = "false"
		touchCmd2 = "touch file2.txt"
	}

	commands := []string{touchCmd1, failCmd, touchCmd2}
	err := kw.RunInitHooks(commands)

	// Should return an error
	if err == nil {
		t.Error("RunInitHooks should return error when a command fails")
	}

	// First file should exist (command executed before failure)
	if _, err := os.Stat(file1); os.IsNotExist(err) {
		t.Error("First command should have been executed successfully")
	}

	// Second file should not exist (command after failure should not execute)
	if _, err := os.Stat(file2); !os.IsNotExist(err) {
		t.Error("Commands after failure should not be executed")
	}
}

func TestRunInitHooks_MultipleValidCommands(t *testing.T) {
	tempDir, cleanup := setupTempWorktree(t)
	defer cleanup()

	kw := NewKoshoWorktree(tempDir, "test-worktree")

	// Create multiple successful commands
	file1 := filepath.Join(kw.WorktreePath(), "file1.txt")
	file2 := filepath.Join(kw.WorktreePath(), "file2.txt")

	var touchCmd1, touchCmd2 string
	if runtime.GOOS == "windows" {
		touchCmd1 = "type nul > file1.txt"
		touchCmd2 = "type nul > file2.txt"
	} else {
		touchCmd1 = "touch file1.txt"
		touchCmd2 = "touch file2.txt"
	}

	commands := []string{touchCmd1, touchCmd2}
	err := kw.RunInitHooks(commands)

	// Should succeed
	if err != nil {
		t.Errorf("RunInitHooks with all valid commands should succeed, got: %v", err)
	}

	// Both files should exist
	if _, err := os.Stat(file1); os.IsNotExist(err) {
		t.Error("First file should have been created")
	}
	if _, err := os.Stat(file2); os.IsNotExist(err) {
		t.Error("Second file should have been created")
	}
}

func TestShellCommand_CrossPlatform(t *testing.T) {
	shellCmd, shellArg := ShellCommand()

	if runtime.GOOS == "windows" {
		if shellCmd != "cmd.exe" {
			t.Errorf("On Windows, expected shell command 'cmd.exe', got: %s", shellCmd)
		}
		if shellArg != "/C" {
			t.Errorf("On Windows, expected shell argument '/C', got: %s", shellArg)
		}
	} else {
		// On Unix-like systems, should return a valid shell path
		if shellCmd == "" {
			t.Error("On Unix-like systems, expected non-empty shell command")
		}
		if shellArg != "-c" {
			t.Errorf("On Unix-like systems, expected shell argument '-c', got: %s", shellArg)
		}
		
		// The shell command should be executable
		if _, err := exec.LookPath(shellCmd); err != nil {
			t.Errorf("Shell command '%s' is not executable: %v", shellCmd, err)
		}
	}
}

func TestRunInitHooks_WorkingDirectory(t *testing.T) {
	tempDir, cleanup := setupTempWorktree(t)
	defer cleanup()

	kw := NewKoshoWorktree(tempDir, "test-worktree")

	// Create a command that creates a file in the current working directory
	var pwdCmd string
	if runtime.GOOS == "windows" {
		pwdCmd = "echo %cd% > current-dir.txt"
	} else {
		pwdCmd = "pwd > current-dir.txt"
	}

	err := kw.RunInitHooks([]string{pwdCmd})
	if err != nil {
		t.Errorf("RunInitHooks should succeed, got: %v", err)
	}

	// Read the file and verify it contains the worktree path
	outputFile := filepath.Join(kw.WorktreePath(), "current-dir.txt")
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Errorf("Failed to read output file: %v", err)
		return
	}

	expectedPath := kw.WorktreePath()
	actualPath := strings.TrimSpace(string(content))

	// Normalize paths for comparison (handle Windows path differences)
	expectedPath = filepath.Clean(expectedPath)
	actualPath = filepath.Clean(actualPath)

	if actualPath != expectedPath {
		t.Errorf("Command should run in worktree directory. Expected: %s, Got: %s", expectedPath, actualPath)
	}
}

// setupTempWorktree creates a temporary directory structure for testing
// and returns the temp dir path and a cleanup function
func setupTempWorktree(t *testing.T) (string, func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "kosho-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create the .kosho directory structure
	koshoDir := filepath.Join(tempDir, ".kosho")
	if err := os.MkdirAll(koshoDir, 0755); err != nil {
		t.Fatalf("Failed to create .kosho directory: %v", err)
	}

	// Create the worktree directory
	worktreeDir := filepath.Join(koshoDir, "test-worktree")
	if err := os.MkdirAll(worktreeDir, 0755); err != nil {
		t.Fatalf("Failed to create worktree directory: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}
