package internal

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRunPostCreateHook_NoHookFile(t *testing.T) {
	tempDir, cleanup := setupTempWorktree(t)
	defer cleanup()

	kw := NewKoshoWorktree(tempDir, "test-worktree")
	hooksDir := filepath.Join(tempDir, ".kosho", "_hooks")

	// Test with non-existent hooks directory
	err := kw.RunPostCreateHook(hooksDir)
	if err != nil {
		t.Errorf("RunPostCreateHook with no hooks directory should return nil, got: %v", err)
	}

	// Create hooks directory but no post-create script
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		t.Fatalf("Failed to create hooks directory: %v", err)
	}

	err = kw.RunPostCreateHook(hooksDir)
	if err != nil {
		t.Errorf("RunPostCreateHook with no post-create script should return nil, got: %v", err)
	}
}

func TestRunPostCreateHook_ValidHook(t *testing.T) {
	tempDir, cleanup := setupTempWorktree(t)
	defer cleanup()

	kw := NewKoshoWorktree(tempDir, "test-worktree")
	hooksDir := filepath.Join(tempDir, ".kosho", "_hooks")

	// Create hooks directory
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		t.Fatalf("Failed to create hooks directory: %v", err)
	}

	// Create a simple post-create hook script
	hookScript := filepath.Join(hooksDir, "post-create")
	var scriptContent string
	var fileMode os.FileMode = 0755
	if runtime.GOOS == "windows" {
		fileMode = 0644
		scriptContent = `@echo off
echo Post-create hook executed
echo %KOSHO_WORKTREE_NAME% > hook-output.txt
`
	} else {
		scriptContent = `#!/bin/bash
echo "Post-create hook executed"
echo "$KOSHO_WORKTREE_NAME" > hook-output.txt
`
	}

	if err := os.WriteFile(hookScript, []byte(scriptContent), fileMode); err != nil {
		t.Fatalf("Failed to create hook script: %v", err)
	}

	// Run the hook
	err := kw.RunPostCreateHook(hooksDir)
	if err != nil {
		t.Errorf("RunPostCreateHook with valid script should succeed, got: %v", err)
	}

	// Verify the hook created the expected output file
	outputFile := filepath.Join(kw.WorktreePath(), "hook-output.txt")
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Errorf("Expected hook to create output file, got error: %v", err)
	}

	expectedContent := "test-worktree"
	if strings.TrimSpace(string(content)) != expectedContent {
		t.Errorf("Expected hook output to contain %q, got %q", expectedContent, strings.TrimSpace(string(content)))
	}
}

func TestRunPostCreateHook_NonExecutableHook(t *testing.T) {
	tempDir, cleanup := setupTempWorktree(t)
	defer cleanup()

	kw := NewKoshoWorktree(tempDir, "test-worktree")
	hooksDir := filepath.Join(tempDir, ".kosho", "_hooks")

	// Create hooks directory
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		t.Fatalf("Failed to create hooks directory: %v", err)
	}

	// Create a non-executable post-create hook
	hookScript := filepath.Join(hooksDir, "post-create")
	scriptContent := "#!/bin/bash\necho 'This should not run'\n"

	if err := os.WriteFile(hookScript, []byte(scriptContent), 0644); err != nil {
		t.Fatalf("Failed to create hook script: %v", err)
	}

	// Run the hook - should fail due to lack of execute permission
	err := kw.RunPostCreateHook(hooksDir)
	if err == nil {
		t.Error("RunPostCreateHook with non-executable script should return an error")
	}

	if !strings.Contains(err.Error(), "not executable") {
		t.Errorf("Error should mention non-executable hook, got: %v", err)
	}
}

func TestRunPostCreateHook_FailingHook(t *testing.T) {
	tempDir, cleanup := setupTempWorktree(t)
	defer cleanup()

	kw := NewKoshoWorktree(tempDir, "test-worktree")
	hooksDir := filepath.Join(tempDir, ".kosho", "_hooks")

	// Create hooks directory
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		t.Fatalf("Failed to create hooks directory: %v", err)
	}

	// Create a hook script that exits with error
	hookScript := filepath.Join(hooksDir, "post-create")
	var scriptContent string
	var fileMode os.FileMode = 0755
	if runtime.GOOS == "windows" {
		fileMode = 0644
		scriptContent = `@echo off
echo This hook will fail
exit 1
`
	} else {
		scriptContent = `#!/bin/bash
echo "This hook will fail"
exit 1
`
	}

	if err := os.WriteFile(hookScript, []byte(scriptContent), fileMode); err != nil {
		t.Fatalf("Failed to create hook script: %v", err)
	}

	// Run the hook - should fail
	err := kw.RunPostCreateHook(hooksDir)
	if err == nil {
		t.Error("RunPostCreateHook with failing script should return an error")
	}

	if !strings.Contains(err.Error(), "post-create hook failed") {
		t.Errorf("Error should mention post-create hook failure, got: %v", err)
	}
}

func TestRunPostCreateHook_EnvironmentVariables(t *testing.T) {
	tempDir, cleanup := setupTempWorktree(t)
	defer cleanup()

	kw := NewKoshoWorktree(tempDir, "test-worktree")
	hooksDir := filepath.Join(tempDir, ".kosho", "_hooks")

	// Create hooks directory
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		t.Fatalf("Failed to create hooks directory: %v", err)
	}

	// Create a hook script that uses environment variables
	hookScript := filepath.Join(hooksDir, "post-create")
	var scriptContent string
	var fileMode os.FileMode = 0755
	if runtime.GOOS == "windows" {
		fileMode = 0644
		scriptContent = `@echo off
echo %KOSHO_WORKTREE_NAME% > name.txt
echo %KOSHO_WORKTREE_PATH% > path.txt
`
	} else {
		scriptContent = `#!/bin/bash
echo "$KOSHO_WORKTREE_NAME" > name.txt
echo "$KOSHO_WORKTREE_PATH" > path.txt
`
	}

	if err := os.WriteFile(hookScript, []byte(scriptContent), fileMode); err != nil {
		t.Fatalf("Failed to create hook script: %v", err)
	}

	// Run the hook
	err := kw.RunPostCreateHook(hooksDir)
	if err != nil {
		t.Errorf("RunPostCreateHook should succeed, got: %v", err)
	}

	// Verify environment variables were passed correctly
	nameFile := filepath.Join(kw.WorktreePath(), "name.txt")
	nameContent, err := os.ReadFile(nameFile)
	if err != nil {
		t.Errorf("Expected hook to create name.txt, got error: %v", err)
	}
	if strings.TrimSpace(string(nameContent)) != "test-worktree" {
		t.Errorf("Expected KOSHO_WORKTREE_NAME to be 'test-worktree', got %q", strings.TrimSpace(string(nameContent)))
	}

	pathFile := filepath.Join(kw.WorktreePath(), "path.txt")
	pathContent, err := os.ReadFile(pathFile)
	if err != nil {
		t.Errorf("Expected hook to create path.txt, got error: %v", err)
	}
	expectedPath := kw.WorktreePath()
	if strings.TrimSpace(string(pathContent)) != expectedPath {
		t.Errorf("Expected KOSHO_WORKTREE_PATH to be %q, got %q", expectedPath, strings.TrimSpace(string(pathContent)))
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
