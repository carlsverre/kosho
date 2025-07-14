package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"kosho/internal"

	"github.com/spf13/cobra"
)

func checkMergeArgs(cmd *cobra.Command, args []string) error {
	args, _ = internal.SplitArgs(cmd, args)
	if len(args) < 1 {
		return fmt.Errorf("worktree argument is required")
	}
	if len(args) > 1 {
		return fmt.Errorf("too many arguments, expected at most 1 (worktree name)")
	}
	return nil
}

var mergeCmd = &cobra.Command{
	Use:   "merge [worktree] [-- git-merge-args...]",
	Short: "Merge a worktree branch into the current branch",
	Args:  checkMergeArgs,
	Long: `Merge a worktree branch into the current branch of the main repository.

The worktree must be clean (no uncommitted changes) and the current branch
must be an ancestor of the worktree branch for the merge to proceed.

Any arguments after -- are passed directly to git merge.`,
	// Args: checkMergeArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		args, mergeArgs := internal.SplitArgs(cmd, args)
		worktree := args[0]

		// Find git root
		repoRoot, err := internal.FindGitRoot()
		if err != nil {
			return fmt.Errorf("failed to find git repository: %w", err)
		}

		kw := internal.NewKoshoWorktree(repoRoot, worktree)

		// Check if worktree exists
		exists, err := kw.Exists()
		if err != nil {
			return fmt.Errorf("failed to check worktree path: %w", err)
		}
		if !exists {
			return fmt.Errorf("worktree `%s` not found", worktree)
		}

		// Check if worktree is clean
		isDirty, err := kw.IsDirty()
		if err != nil {
			return fmt.Errorf("failed to check worktree status: %w", err)
		}
		if isDirty {
			return fmt.Errorf("worktree `%s` has uncommitted changes. Commit or stash first", worktree)
		}

		// Get current branch of main repo
		currentBranchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		currentBranchCmd.Dir = repoRoot
		currentBranchOutput, err := currentBranchCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}
		currentBranch := strings.TrimSpace(string(currentBranchOutput))

		// Get worktree branch
		worktreeBranch, err := kw.GitBranch()
		if err != nil {
			return fmt.Errorf("failed to get worktree branch: %w", err)
		}

		// Check if current branch is ancestor of worktree branch
		isAncestor, err := internal.IsAncestor(repoRoot, currentBranch, worktreeBranch)
		if err != nil {
			return fmt.Errorf("failed to check ancestry: %w", err)
		}
		if !isAncestor {
			return fmt.Errorf("current branch '%s' is not an ancestor of worktree branch '%s'. Use 'git log --oneline --graph %s..%s' to see divergent commits",
				currentBranch, worktreeBranch, currentBranch, worktreeBranch)
		}

		// Perform the merge
		mergeCmd := []string{"git", "merge"}
		mergeCmd = append(mergeCmd, mergeArgs...)
		mergeCmd = append(mergeCmd, worktreeBranch)

		fmt.Printf("Merging worktree branch '%s' into '%s'...\n", worktreeBranch, currentBranch)

		gitMerge := exec.Command(mergeCmd[0], mergeCmd[1:]...)
		gitMerge.Dir = repoRoot
		gitMerge.Stdin = os.Stdin
		gitMerge.Stdout = os.Stdout
		gitMerge.Stderr = os.Stderr

		err = gitMerge.Run()
		if err != nil {
			return fmt.Errorf("git merge failed: %w", err)
		}

		fmt.Printf("Successfully merged '%s' into '%s'\n", worktreeBranch, currentBranch)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
}
