package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"kosho/internal/git"
)

var forceFlag bool

var removeCmd = &cobra.Command{
	Use:   "remove [flags] NAME",
	Short: "Remove a kosho worktree",
	Long: `Remove a kosho worktree and stop any running container.
If the worktree is dirty, use --force to continue.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Get current directory and find git root
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		repoRoot, err := git.FindGitRoot(currentDir)
		if err != nil {
			return fmt.Errorf("failed to find git repository: %w", err)
		}

		worktreePath := filepath.Join(repoRoot, ".kosho", name)

		// Check if worktree exists
		if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
			return fmt.Errorf("worktree '%s' does not exist", name)
		}

		// Check if worktree is dirty (has uncommitted changes)
		if !forceFlag {
			isDirty, err := isWorktreeDirty(worktreePath)
			if err != nil {
				return fmt.Errorf("failed to check worktree status: %w", err)
			}
			if isDirty {
				return fmt.Errorf("worktree '%s' has uncommitted changes, use --force to remove anyway", name)
			}
		}

		// Stub: Stop container if running
		fmt.Printf("Stopping container for worktree '%s' (if running)\n", name)

		// Remove the git worktree
		fmt.Printf("Removing worktree '%s'\n", name)

		gitArgs := []string{"worktree", "remove", worktreePath}
		if forceFlag {
			gitArgs = append(gitArgs, "--force")
		}

		gitCmd := exec.Command("git", gitArgs...)
		gitCmd.Dir = repoRoot

		output, err := gitCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to remove worktree: %w\nOutput: %s", err, string(output))
		}

		fmt.Printf("Worktree '%s' removed successfully\n", name)
		return nil
	},
}

func isWorktreeDirty(worktreePath string) (bool, error) {
	// Check if there are uncommitted changes
	gitCmd := exec.Command("git", "status", "--porcelain")
	gitCmd.Dir = worktreePath

	output, err := gitCmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to get git status: %w", err)
	}

	// If output is empty, worktree is clean
	return len(output) > 0, nil
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "force removal even if worktree is dirty")
}
