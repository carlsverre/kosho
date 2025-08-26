package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/carlsverre/kosho/internal"

	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all kosho worktrees",
	Long:  `List all kosho worktrees and their current git status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		koshoDir, err := internal.NewKoshoDir()
		if err != nil {
			return fmt.Errorf("failed to load Kosho dir: %w", err)
		}

		worktrees, err := koshoDir.ListWorktrees()
		if err != nil {
			return fmt.Errorf("failed to list worktrees: %w", err)
		}

		if len(worktrees) == 0 {
			fmt.Println("No kosho worktrees found")
			return nil
		}

		// Get current branch of main repo
		currentBranchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		currentBranchCmd.Dir = koshoDir.RepoPath()
		currentBranchOutput, err := currentBranchCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}
		currentBranch := strings.TrimSpace(string(currentBranchOutput))

		tbl := table.New("NAME", "STATUS", "REF", "MERGE-READY")
		for _, kw := range worktrees {
			// Check git status
			isDirty, err := kw.IsDirty()
			var gitStatus string
			if err != nil {
				gitStatus = "error"
			} else if isDirty {
				gitStatus = "dirty"
			} else {
				gitStatus = "clean"
			}

			// Get current branch/ref
			gitRef, err := kw.GitRef()
			if err != nil {
				gitRef = "unknown"
			}

			// Check if worktree has outstanding commits
			var mergeReady string
			hasCommits, err := kw.HasOutstandingCommits(currentBranch)
			if err != nil {
				mergeReady = "error"
			} else if hasCommits {
				mergeReady = "yes"
			} else {
				mergeReady = "no"
			}

			tbl.AddRow(kw.Name(), gitStatus, gitRef, mergeReady)
		}

		tbl.Print()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
