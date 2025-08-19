package internal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lithammer/dedent"
)

var (
	KoshoGitIgnore = []byte(dedent.Dedent(`
	    *
		!/.gitignore
		!/kosho_config/
		!/kosho_config/**
	`))
)

func SetupKoshoDir() error {
	repoRoot, err := FindGitRoot()
	if err != nil {
		return fmt.Errorf("failed to find git repository: %w", err)
	}

	// if the root .gitignore contains .kosho, remove it
	rootGitIgnorePath := filepath.Join(repoRoot, ".gitignore")
	if file, err := os.Open(rootGitIgnorePath); err == nil {
		defer func() { _ = file.Close() }()
		scanner := bufio.NewScanner(file)
		var lines []string
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.Contains(line, ".kosho") {
				lines = append(lines, line)
			}
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading .gitignore: %w", err)
		}
		// Rewrite .gitignore without the .kosho line
		if err := os.WriteFile(rootGitIgnorePath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
			return fmt.Errorf("failed to rewrite .gitignore: %w", err)
		}
	}

	// Create the .kosho directory if it doesn't exist
	koshoDir := filepath.Join(repoRoot, ".kosho")
	if _, err := os.Stat(koshoDir); os.IsNotExist(err) {
		if err := os.Mkdir(koshoDir, 0755); err != nil {
			return fmt.Errorf("failed to create .kosho directory: %w", err)
		}
	}

	// Create .kosho/.gitignore if it doesn't exist
	koshoGitIgnorePath := filepath.Join(koshoDir, ".gitignore")
	if _, err := os.Stat(koshoGitIgnorePath); os.IsNotExist(err) {
		if err := os.WriteFile(koshoGitIgnorePath, KoshoGitIgnore, 0644); err != nil {
			return fmt.Errorf("failed to create .kosho/.gitignore: %w", err)
		}
	}

	return nil
}
