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
			if line == "/.kosho**" {
				return nil // Already present
			}
		}
	}

	// Update .gitignore
	file, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open .gitignore: %w", err)
	}
	defer func() { _ = file.Close() }()

	_, err = file.WriteString("/.kosho**\n")
	if err != nil {
		return fmt.Errorf("failed to write to .gitignore: %w", err)
	}

	return nil
}
