package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"kosho/internal"
)

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Cleanup dangling worktree references",
	Long:  `Cleanup any dangling worktree references by running 'git worktree prune'.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Find git root
		repoRoot, err := internal.FindGitRoot()
		if err != nil {
			return fmt.Errorf("failed to find git repository: %w", err)
		}

		fmt.Println("Pruning dangling worktree references...")

		// Run git worktree prune
		gitCmd := exec.Command("git", "worktree", "prune")
		gitCmd.Dir = repoRoot

		output, err := gitCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to prune worktrees: %w\nOutput: %s", err, string(output))
		}

		if len(output) > 0 {
			fmt.Printf("Output: %s", string(output))
		} else {
			fmt.Println("No dangling worktree references found")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pruneCmd)
}