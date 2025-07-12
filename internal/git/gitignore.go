package git

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func EnsureKoshoInGitignore(repoRoot string) error {
	gitignorePath := filepath.Join(repoRoot, ".gitignore")

	// Check if .gitignore exists and if /.kosho is already in it
	if file, err := os.Open(gitignorePath); err == nil {
		defer func() { _ = file.Close() }()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "/.kosho" || line == ".kosho" || line == ".kosho/" {
				return nil // Already present
			}
		}
	}

	// Append /.kosho to .gitignore
	file, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open .gitignore: %w", err)
	}
	defer func() { _ = file.Close() }()

	// Add newline before if file exists and doesn't end with newline
	if stat, err := file.Stat(); err == nil && stat.Size() > 0 {
		// Check if file ends with newline
		_, _ = file.Seek(-1, 2)
		lastByte := make([]byte, 1)
		_, _ = file.Read(lastByte)
		if lastByte[0] != '\n' {
			_, _ = file.WriteString("\n")
		}
	}

	_, err = file.WriteString("/.kosho\n")
	if err != nil {
		return fmt.Errorf("failed to write to .gitignore: %w", err)
	}

	return nil
}
