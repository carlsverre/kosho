package internal

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// GetWorktreeNames returns a list of available worktree names for autocompletion
func GetWorktreeNames() ([]string, error) {
	repoRoot, err := FindGitRoot()
	if err != nil {
		return nil, err
	}

	koshoDir := filepath.Join(repoRoot, ".kosho")

	// Check if .kosho directory exists
	if _, err := os.Stat(koshoDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// List directories in .kosho
	entries, err := os.ReadDir(koshoDir)
	if err != nil {
		return nil, err
	}

	var worktreeNames []string
	for _, entry := range entries {
		if entry.IsDir() {
			worktreeNames = append(worktreeNames, entry.Name())
		}
	}

	return worktreeNames, nil
}

// WorktreeCompletionFunc provides autocompletion for worktree names
func WorktreeCompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	worktreeNames, err := GetWorktreeNames()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	return worktreeNames, cobra.ShellCompDirectiveNoFileComp
}
