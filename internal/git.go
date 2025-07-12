package internal

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func FindGitRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("not a git repository (or any of the parent directories): %w", err)
	}

	gitRoot := strings.TrimSpace(string(output))
	return gitRoot, nil
}

func EnsureGitIgnored(glob string) error {
	repoRoot, err := FindGitRoot()
	if err != nil {
		return fmt.Errorf("failed to find git repository: %w", err)
	}
	gitignorePath := filepath.Join(repoRoot, ".gitignore")

	// Check if .gitignore exists and if /.kosho is already in it
	if file, err := os.Open(gitignorePath); err == nil {
		defer func() { _ = file.Close() }()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == glob {
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

	_, err = fmt.Fprintf(file, "%s\n", glob)
	if err != nil {
		return fmt.Errorf("failed to write to .gitignore: %w", err)
	}

	return nil
}
