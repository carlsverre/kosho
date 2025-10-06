package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/carlsverre/kosho/internal"

	"github.com/spf13/cobra"
)

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Cleanup clean worktrees and dangling worktree references",
	RunE: func(cmd *cobra.Command, args []string) error {
		koshoDir, err := internal.LoadKoshoDir()
		if err != nil {
			return fmt.Errorf("failed to load Kosho dir: %w", err)
		}

		worktrees, err := koshoDir.ListWorktrees()
		if err != nil {
			return fmt.Errorf("failed to list worktrees: %w", err)
		}

		// iterate through worktrees, removing any clean worktrees
		for _, worktree := range worktrees {
			clean, err := worktree.IsClean()
			if err != nil {
				return fmt.Errorf("failed to check worktree status: %w", err)
			}
			if clean {
				err := worktree.Remove(false)
				if err != nil {
					return fmt.Errorf("failed to remove worktree %s: %w", worktree.Name(), err)
				}
			}
		}

		// Run git worktree prune
		gitCmd := exec.Command("git", "worktree", "prune", "--verbose")
		gitCmd.Dir = koshoDir.RepoPath()
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr
		err = gitCmd.Run()
		if err != nil {
			return fmt.Errorf("failed to prune worktrees: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pruneCmd)
}
