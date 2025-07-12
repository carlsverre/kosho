package dirs

import (
	"os"
	"path/filepath"
)

func GetDataDir() string {
	if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		return xdgDataHome
	}
	
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	
	return filepath.Join(homeDir, ".local", "share")
}

func GetWorktreeDir(repoName string) string {
	dataDir := GetDataDir()
	if dataDir == "" {
		return ""
	}
	
	return filepath.Join(dataDir, repoName)
}