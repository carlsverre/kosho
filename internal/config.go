package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Settings represents the configuration stored in .kosho/settings.json
type Settings struct {
	WorktreeInit []string `json:"worktree_init"`
}

// LoadSettings reads and parses .kosho/settings.json from the repository root
func LoadSettings(repoRoot string) (Settings, error) {
	settingsPath := filepath.Join(repoRoot, ".kosho", "settings.json")

	// Return zero-value Settings if file doesn't exist
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		return Settings{}, nil
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return Settings{}, fmt.Errorf("failed to read .kosho/settings.json: %w", err)
	}

	// First unmarshal into a map to detect unknown keys
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return Settings{}, fmt.Errorf("kosho: invalid .kosho/settings.json – %w", err)
	}

	// Check for unknown keys and warn about them
	knownKeys := map[string]bool{
		"worktree_init": true,
	}
	for key := range rawData {
		if !knownKeys[key] {
			fmt.Fprintf(os.Stderr, "kosho: warning: unknown configuration key %q in .kosho/settings.json\n", key)
		}
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return Settings{}, fmt.Errorf("kosho: invalid .kosho/settings.json – %w", err)
	}

	// Validate worktree_init array
	if err := validateWorktreeInit(settings.WorktreeInit); err != nil {
		return Settings{}, fmt.Errorf("kosho: invalid .kosho/settings.json – %w", err)
	}

	return settings, nil
}

func validateWorktreeInit(commands []string) error {
	for i, cmd := range commands {
		if strings.TrimSpace(cmd) == "" {
			return fmt.Errorf("worktree_init[%d] cannot be empty or whitespace-only", i)
		}
	}
	return nil
}
