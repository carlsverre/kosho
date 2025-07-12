package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"kosho/internal/git"
	"kosho/internal/worktree"
)

var (
	branchFlag    string
	newBranchFlag string
)

var newCmd = &cobra.Command{
	Use:   "new [flags] NAME [commitish]",
	Short: "Create a new kosho worktree",
	Long: `Create a new kosho worktree in .kosho/NAME at the repo root.
The NAME parameter is required. An optional commitish can be provided.
Use -b to create a new branch or -B to create/reset a branch.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		var commitish string
		if len(args) > 1 {
			commitish = args[1]
		}

		// Get current directory and find git root
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		repoRoot, err := git.FindGitRoot(currentDir)
		if err != nil {
			return fmt.Errorf("failed to find git repository: %w", err)
		}

		kw := worktree.NewKoshoWorktree(repoRoot, name)

		// Check if worktree already exists
		if _, err := os.Stat(kw.WorktreePath()); err == nil {
			return fmt.Errorf("worktree '%s' already exists", name)
		}

		// Ensure .kosho is in .gitignore
		err = git.EnsureKoshoInGitignore(repoRoot)
		if err != nil {
			return fmt.Errorf("failed to update .gitignore: %w", err)
		}

		fmt.Printf("Creating worktree '%s' in %s\n", name, kw.WorktreePath())

		// Create the worktree
		err = git.CreateKoshoWorktree(repoRoot, name, kw.KoshoDir(), branchFlag, newBranchFlag, commitish)
		if err != nil {
			return fmt.Errorf("failed to create worktree: %w", err)
		}

		fmt.Printf("Worktree created successfully at %s\n", kw.WorktreePath())

		// Fall through to start command
		return startWorktree(name)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringVarP(&branchFlag, "branch", "b", "", "create a new branch")
	newCmd.Flags().StringVarP(&newBranchFlag, "new-branch", "B", "", "create/reset a branch")
}
