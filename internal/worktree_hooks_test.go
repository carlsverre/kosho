package internal

import (
	"os"
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

	// Test with echo command (simple binary execution)
	err := kw.RunInitHooks([]string{"echo test message"})
	if err != nil {
		t.Errorf("RunInitHooks with valid echo command should succeed, got: %v", err)
	}

	// Test with touch command to create a file
	testFile := filepath.Join(kw.WorktreePath(), "test-file.txt")
	touchCmd := "touch test-file.txt"

	err = kw.RunInitHooks([]string{touchCmd})
	if err != nil {
		t.Errorf("RunInitHooks with touch command should succeed, got: %v", err)
	}

	// Verify the file was created (skip on Windows since touch might not be available)
	if runtime.GOOS != "windows" {
		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Errorf("Expected file %s to be created by touch command", testFile)
		}
	}
}

func TestRunInitHooks_FailingCommand(t *testing.T) {
	tempDir, cleanup := setupTempWorktree(t)
	defer cleanup()

	kw := NewKoshoWorktree(tempDir, "test-worktree")

	// Test with command that always fails (nonexistent command)
	failCmd := "nonexistent-command-12345"

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

	touchCmd1 := "touch file1.txt"
	failCmd := "nonexistent-command-12345"
	touchCmd2 := "touch file2.txt"

	commands := []string{touchCmd1, failCmd, touchCmd2}
	err := kw.RunInitHooks(commands)

	// Should return an error
	if err == nil {
		t.Error("RunInitHooks should return error when a command fails")
	}

	// Skip file checks on Windows since touch might not be available
	if runtime.GOOS != "windows" {
		// First file should exist (command executed before failure)
		if _, err := os.Stat(file1); os.IsNotExist(err) {
			t.Error("First command should have been executed successfully")
		}

		// Second file should not exist (command after failure should not execute)
		if _, err := os.Stat(file2); !os.IsNotExist(err) {
			t.Error("Commands after failure should not be executed")
		}
	}
}

func TestRunInitHooks_MultipleValidCommands(t *testing.T) {
	tempDir, cleanup := setupTempWorktree(t)
	defer cleanup()

	kw := NewKoshoWorktree(tempDir, "test-worktree")

	// Create multiple successful commands (using echo which is cross-platform)
	commands := []string{"echo first command", "echo second command"}
	err := kw.RunInitHooks(commands)

	// Should succeed
	if err != nil {
		t.Errorf("RunInitHooks with all valid commands should succeed, got: %v", err)
	}
}



func TestRunInitHooks_WorkingDirectory(t *testing.T) {
	tempDir, cleanup := setupTempWorktree(t)
	defer cleanup()

	kw := NewKoshoWorktree(tempDir, "test-worktree")

	// Simple test - just verify commands run without shell-specific syntax
	err := kw.RunInitHooks([]string{"echo working directory test"})
	if err != nil {
		t.Errorf("RunInitHooks should succeed, got: %v", err)
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
