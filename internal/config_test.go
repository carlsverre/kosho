package internal

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadSettings_ValidJSON(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "kosho-config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create .kosho directory
	koshoDir := filepath.Join(tempDir, ".kosho")
	if err := os.MkdirAll(koshoDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create valid settings.json
	settings := Settings{
		WorktreeInit: []string{"pnpm i", "echo 'initialized'"},
	}
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatal(err)
	}

	settingsPath := filepath.Join(koshoDir, "settings.json")
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		t.Fatal(err)
	}

	// Test LoadSettings
	result, err := LoadSettings(tempDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result.WorktreeInit) != 2 {
		t.Fatalf("Expected 2 commands, got %d", len(result.WorktreeInit))
	}

	if result.WorktreeInit[0] != "pnpm i" {
		t.Errorf("Expected 'pnpm i', got '%s'", result.WorktreeInit[0])
	}

	if result.WorktreeInit[1] != "echo 'initialized'" {
		t.Errorf("Expected 'echo 'initialized'', got '%s'", result.WorktreeInit[1])
	}
}

func TestLoadSettings_MissingFile(t *testing.T) {
	// Create temp directory without settings.json
	tempDir, err := os.MkdirTemp("", "kosho-config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Test LoadSettings - should return default/zero-value Settings
	result, err := LoadSettings(tempDir)
	if err != nil {
		t.Fatalf("Expected no error for missing file, got: %v", err)
	}

	if len(result.WorktreeInit) != 0 {
		t.Errorf("Expected empty WorktreeInit slice, got %d commands", len(result.WorktreeInit))
	}
}

func TestLoadSettings_InvalidJSON(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "kosho-config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create .kosho directory
	koshoDir := filepath.Join(tempDir, ".kosho")
	if err := os.MkdirAll(koshoDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create invalid JSON file
	settingsPath := filepath.Join(koshoDir, "settings.json")
	invalidJSON := `{
		"worktree_init": ["pnpm i"
	}` // Missing closing bracket
	if err := os.WriteFile(settingsPath, []byte(invalidJSON), 0644); err != nil {
		t.Fatal(err)
	}

	// Test LoadSettings - should return error
	_, err = LoadSettings(tempDir)
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}

	expectedPrefix := "kosho: invalid .kosho/settings.json –"
	if len(err.Error()) < len(expectedPrefix) || err.Error()[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("Expected error message to start with '%s', got: %s", expectedPrefix, err.Error())
	}
}

func TestLoadSettings_InvalidWorktreeInitFieldTypes(t *testing.T) {
	tests := []struct {
		name        string
		jsonContent string
		expectError bool
	}{
		{
			name:        "worktree_init as string instead of array",
			jsonContent: `{"worktree_init": "pnpm i"}`,
			expectError: true,
		},
		{
			name:        "worktree_init as number",
			jsonContent: `{"worktree_init": 123}`,
			expectError: true,
		},
		{
			name:        "worktree_init with non-string array elements",
			jsonContent: `{"worktree_init": [123, "pnpm i"]}`,
			expectError: true,
		},
		{
			name:        "worktree_init as object",
			jsonContent: `{"worktree_init": {"cmd": "pnpm i"}}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir, err := os.MkdirTemp("", "kosho-config-test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			// Create .kosho directory
			koshoDir := filepath.Join(tempDir, ".kosho")
			if err := os.MkdirAll(koshoDir, 0755); err != nil {
				t.Fatal(err)
			}

			// Create settings.json with test content
			settingsPath := filepath.Join(koshoDir, "settings.json")
			if err := os.WriteFile(settingsPath, []byte(tt.jsonContent), 0644); err != nil {
				t.Fatal(err)
			}

			// Test LoadSettings
			_, err = LoadSettings(tempDir)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for %s, got nil", tt.name)
			} else if !tt.expectError && err != nil {
				t.Errorf("Expected no error for %s, got: %v", tt.name, err)
			}

			if tt.expectError && err != nil {
				expectedPrefix := "kosho: invalid .kosho/settings.json –"
				if len(err.Error()) < len(expectedPrefix) || err.Error()[:len(expectedPrefix)] != expectedPrefix {
					t.Errorf("Expected error message to start with '%s', got: %s", expectedPrefix, err.Error())
				}
			}
		})
	}
}

func TestLoadSettings_EmptyWorktreeInitArray(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "kosho-config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create .kosho directory
	koshoDir := filepath.Join(tempDir, ".kosho")
	if err := os.MkdirAll(koshoDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create settings.json with empty worktree_init array
	settingsPath := filepath.Join(koshoDir, "settings.json")
	jsonContent := `{"worktree_init": []}`
	if err := os.WriteFile(settingsPath, []byte(jsonContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Test LoadSettings - empty array should be valid
	result, err := LoadSettings(tempDir)
	if err != nil {
		t.Fatalf("Expected no error for empty worktree_init array, got: %v", err)
	}

	if len(result.WorktreeInit) != 0 {
		t.Errorf("Expected empty WorktreeInit slice, got %d commands", len(result.WorktreeInit))
	}
}

func TestLoadSettings_WorktreeInitWithEmptyStrings(t *testing.T) {
	tests := []struct {
		name        string
		commands    []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "single empty string",
			commands:    []string{""},
			expectError: true,
			errorMsg:    "worktree_init[0] cannot be empty or whitespace-only",
		},
		{
			name:        "empty string in middle",
			commands:    []string{"pnpm i", "", "echo done"},
			expectError: true,
			errorMsg:    "worktree_init[1] cannot be empty or whitespace-only",
		},
		{
			name:        "empty string at end",
			commands:    []string{"pnpm i", "echo done", ""},
			expectError: true,
			errorMsg:    "worktree_init[2] cannot be empty or whitespace-only",
		},
		{
			name:        "multiple empty strings",
			commands:    []string{"", "pnpm i", ""},
			expectError: true,
			errorMsg:    "worktree_init[0] cannot be empty or whitespace-only",
		},
		{
			name:        "whitespace-only string",
			commands:    []string{"   "},
			expectError: true,
			errorMsg:    "worktree_init[0] cannot be empty or whitespace-only",
		},
		{
			name:        "tabs and spaces",
			commands:    []string{"\t  \n  "},
			expectError: true,
			errorMsg:    "worktree_init[0] cannot be empty or whitespace-only",
		},
		{
			name:        "mixed whitespace in middle",
			commands:    []string{"pnpm i", "  \t  ", "echo done"},
			expectError: true,
			errorMsg:    "worktree_init[1] cannot be empty or whitespace-only",
		},
		{
			name:        "valid commands",
			commands:    []string{"pnpm i", "echo 'initialized'"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir, err := os.MkdirTemp("", "kosho-config-test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			// Create .kosho directory
			koshoDir := filepath.Join(tempDir, ".kosho")
			if err := os.MkdirAll(koshoDir, 0755); err != nil {
				t.Fatal(err)
			}

			// Create settings.json
			settings := Settings{WorktreeInit: tt.commands}
			data, err := json.Marshal(settings)
			if err != nil {
				t.Fatal(err)
			}

			settingsPath := filepath.Join(koshoDir, "settings.json")
			if err := os.WriteFile(settingsPath, data, 0644); err != nil {
				t.Fatal(err)
			}

			// Test LoadSettings
			result, err := LoadSettings(tempDir)
			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected error for %s, got nil", tt.name)
				}
				expectedMsg := "kosho: invalid .kosho/settings.json – " + tt.errorMsg
				if err.Error() != expectedMsg {
					t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error for %s, got: %v", tt.name, err)
				}
				if len(result.WorktreeInit) != len(tt.commands) {
					t.Errorf("Expected %d commands, got %d", len(tt.commands), len(result.WorktreeInit))
				}
				for i, cmd := range tt.commands {
					if result.WorktreeInit[i] != cmd {
						t.Errorf("Expected command[%d] = '%s', got '%s'", i, cmd, result.WorktreeInit[i])
					}
				}
			}
		})
	}
}

func TestLoadSettings_MissingKoshoDirectory(t *testing.T) {
	// Create temp directory without .kosho subdirectory
	tempDir, err := os.MkdirTemp("", "kosho-config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Test LoadSettings - should return default Settings (not error) for missing .kosho dir
	result, err := LoadSettings(tempDir)
	if err != nil {
		t.Fatalf("Expected no error for missing .kosho directory, got: %v", err)
	}

	if len(result.WorktreeInit) != 0 {
		t.Errorf("Expected empty WorktreeInit slice, got %d commands", len(result.WorktreeInit))
	}
}

func TestLoadSettings_OtherTopLevelKeys(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "kosho-config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create .kosho directory
	koshoDir := filepath.Join(tempDir, ".kosho")
	if err := os.MkdirAll(koshoDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create settings.json with extra top-level keys (future-proofing test)
	settingsPath := filepath.Join(koshoDir, "settings.json")
	jsonContent := `{
		"worktree_init": ["pnpm i"],
		"unknown_setting": "should be ignored",
		"another_setting": 42
	}`
	if err := os.WriteFile(settingsPath, []byte(jsonContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Capture stderr to check warnings
	originalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Test LoadSettings - should successfully parse and warn about unknown keys
	result, err := LoadSettings(tempDir)
	
	// Restore stderr and read captured output
	w.Close()
	os.Stderr = originalStderr
	var buf bytes.Buffer
	buf.ReadFrom(r)
	stderrOutput := buf.String()

	if err != nil {
		t.Fatalf("Expected no error with extra top-level keys, got: %v", err)
	}

	if len(result.WorktreeInit) != 1 {
		t.Fatalf("Expected 1 command, got %d", len(result.WorktreeInit))
	}

	if result.WorktreeInit[0] != "pnpm i" {
		t.Errorf("Expected 'pnpm i', got '%s'", result.WorktreeInit[0])
	}

	// Check that warnings were printed for unknown keys
	if !strings.Contains(stderrOutput, `kosho: warning: unknown configuration key "unknown_setting" in .kosho/settings.json`) {
		t.Errorf("Expected warning for 'unknown_setting', got stderr: %s", stderrOutput)
	}
	if !strings.Contains(stderrOutput, `kosho: warning: unknown configuration key "another_setting" in .kosho/settings.json`) {
		t.Errorf("Expected warning for 'another_setting', got stderr: %s", stderrOutput)
	}
}
