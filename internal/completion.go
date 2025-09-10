package internal

import (
	"github.com/spf13/cobra"
)

// GetWorktreeNames returns a list of available worktree names for autocompletion
func GetWorktreeNames() ([]string, error) {
	koshoDir, err := NewKoshoDir()
	if err != nil {
		return nil, err
	}

	worktrees, err := koshoDir.ListWorktrees()
	if err != nil {
		return nil, err
	}

	var worktreeNames []string
	for _, worktree := range worktrees {
		worktreeNames = append(worktreeNames, worktree.Name())
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
