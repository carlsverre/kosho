package git

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func EnsureWorktreesInGitignore(repoRoot string) error {
	gitignorePath := filepath.Join(repoRoot, ".gitignore")
	
	// Check if .gitignore exists and if .worktrees is already in it
	if file, err := os.Open(gitignorePath); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == ".worktrees" || line == ".worktrees/" {
				return nil // Already present
			}
		}
	}
	
	// Append .worktrees to .gitignore
	file, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open .gitignore: %w", err)
	}
	defer file.Close()
	
	// Add newline before if file exists and doesn't end with newline
	if stat, err := file.Stat(); err == nil && stat.Size() > 0 {
		// Check if file ends with newline
		file.Seek(-1, 2)
		lastByte := make([]byte, 1)
		file.Read(lastByte)
		if lastByte[0] != '\n' {
			file.WriteString("\n")
		}
	}
	
	_, err = file.WriteString(".worktrees\n")
	if err != nil {
		return fmt.Errorf("failed to write to .gitignore: %w", err)
	}
	
	return nil
}