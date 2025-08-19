package internal

import (
	"path/filepath"
)

// Settings represents the configuration for kosho
type Settings struct {
	HooksDir string
}

// LoadSettings loads kosho configuration from the repository root
func LoadSettings(repoRoot string) (Settings, error) {
	hooksDir := filepath.Join(repoRoot, ".kosho", "_hooks")
	
	return Settings{
		HooksDir: hooksDir,
	}, nil
}
