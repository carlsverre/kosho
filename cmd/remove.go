package cmd

import (
	"fmt"
	"os"

	"github.com/carlsverre/kosho/internal"

	"github.com/spf13/cobra"
)

var forceFlag bool

var removeCmd = &cobra.Command{
	Use:   "remove [flags] NAME",
	Short: "Remove a kosho worktree",
	Long: `Remove a kosho worktree.
If the worktree is dirty, use --force to continue.`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: internal.WorktreeCompletionFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Find git root
		repoRoot, err := internal.FindGitRoot()
		if err != nil {
			return fmt.Errorf("failed to find git repository: %w", err)
		}

		kw := internal.NewKoshoWorktree(repoRoot, name)

		// Check if worktree exists
		if _, err := os.Stat(kw.WorktreePath()); os.IsNotExist(err) {
			return fmt.Errorf("worktree '%s' does not exist", name)
		}

		// Check if worktree is dirty (has uncommitted changes)
		if !forceFlag {
			isDirty, err := kw.IsDirty()
			if err != nil {
				return fmt.Errorf("failed to check worktree status: %w", err)
			}
			if isDirty {
				return fmt.Errorf("worktree '%s' has uncommitted changes, use --force to remove anyway", name)
			}
		}

		// Remove the git worktree
		fmt.Printf("Removing worktree '%s'\n", name)

		err = kw.Remove(forceFlag)
		if err != nil {
			return err
		}

		fmt.Printf("Worktree '%s' removed successfully\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "force removal even if worktree is dirty")
}
