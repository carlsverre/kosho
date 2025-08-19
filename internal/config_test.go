package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSettings(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "kosho-config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Test LoadSettings
	result, err := LoadSettings(tempDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedHooksDir := filepath.Join(tempDir, ".kosho", "_hooks")
	if result.HooksDir != expectedHooksDir {
		t.Errorf("Expected HooksDir to be %q, got %q", expectedHooksDir, result.HooksDir)
	}
}
