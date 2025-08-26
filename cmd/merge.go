package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/carlsverre/kosho/internal"

	"github.com/spf13/cobra"
)

var keepFlag bool
var commitMessage string

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
	Use:               "merge [worktree] [-- git-merge-args...]",
	Short:             "Merge a worktree branch into the current branch",
	Args:              checkMergeArgs,
	ValidArgsFunction: internal.WorktreeCompletionFunc,
	Long: `Merge a worktree branch into the current branch of the main repository.

The worktree must be clean (no uncommitted changes) and the current branch
must be an ancestor of the worktree branch for the merge to proceed.

Any arguments after -- are passed directly to git merge.

Worktrees are automatically removed after successful merge unless --keep is specified.

Use -m/--message to commit changes in the worktree before merging.`,
	// Args: checkMergeArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		args, mergeArgs := internal.SplitArgs(cmd, args)
		worktree := args[0]

		koshoDir, err := internal.NewKoshoDir()
		if err != nil {
			return fmt.Errorf("failed to load Kosho dir: %w", err)
		}

		kw := internal.NewKoshoWorktree(*koshoDir, worktree)

		// Check if worktree exists
		exists, err := kw.Exists()
		if err != nil {
			return fmt.Errorf("failed to check worktree path: %w", err)
		}
		if !exists {
			return fmt.Errorf("worktree `%s` not found", worktree)
		}

		// Run the merge hook if it exists
		if err := internal.RunKoshoHook(kw, internal.HOOK_MERGE); err != nil {
			return fmt.Errorf("failed to run merge hook: %w", err)
		}

		// Handle pre-merge commit if message is provided
		if commitMessage != "" {
			// Check for untracked files
			cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
			cmd.Dir = kw.WorktreePath()
			untrackedOutput, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("failed to check for untracked files: %w", err)
			}
			if len(strings.TrimSpace(string(untrackedOutput))) > 0 {
				return fmt.Errorf("worktree `%s` has untracked files. Add them first or remove them before using --message", worktree)
			}

			// Commit with the provided message
			commitCmd := exec.Command("git", "commit", "-am", commitMessage)
			commitCmd.Dir = kw.WorktreePath()
			commitOutput, err := commitCmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("failed to commit changes: %w\nOutput: %s", err, string(commitOutput))
			}
			fmt.Printf("Committed changes in worktree '%s': %s\n", worktree, commitMessage)
		} else {
			// Check if worktree is clean
			isDirty, err := kw.IsDirty()
			if err != nil {
				return fmt.Errorf("failed to check worktree status: %w", err)
			}
			if isDirty {
				return fmt.Errorf("worktree `%s` has uncommitted changes. Commit or stash first, or use -m/--message to commit them", worktree)
			}
		}

		// Get current branch of main repo
		currentBranchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		currentBranchCmd.Dir = koshoDir.RepoPath()
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
		isAncestor, err := internal.IsAncestor(koshoDir.RepoPath(), currentBranch, worktreeBranch)
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
		gitMerge.Dir = koshoDir.RepoPath()
		gitMerge.Stdin = os.Stdin
		gitMerge.Stdout = os.Stdout
		gitMerge.Stderr = os.Stderr

		err = gitMerge.Run()
		if err != nil {
			return fmt.Errorf("git merge failed: %w", err)
		}

		fmt.Printf("Successfully merged '%s' into '%s'\n", worktreeBranch, currentBranch)

		// Remove worktree by default unless --keep flag is set
		if !keepFlag {
			fmt.Printf("Removing worktree '%s'...\n", worktree)
			err = kw.Remove(false)
			if err != nil {
				return fmt.Errorf("merge succeeded but failed to remove worktree: %w", err)
			}
			fmt.Printf("Worktree '%s' removed successfully\n", worktree)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
	mergeCmd.Flags().BoolVar(&keepFlag, "keep", false, "keep worktree after successful merge")
	mergeCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "commit message for changes in worktree before merging")
}
