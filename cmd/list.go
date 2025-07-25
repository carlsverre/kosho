package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
		// Find git root
		repoRoot, err := internal.FindGitRoot()
		if err != nil {
			return fmt.Errorf("failed to find git repository: %w", err)
		}

		koshoDir := filepath.Join(repoRoot, ".kosho")

		// Check if .kosho directory exists
		if _, err := os.Stat(koshoDir); os.IsNotExist(err) {
			fmt.Println("No kosho worktrees found")
			return nil
		}

		// List directories in .kosho
		entries, err := os.ReadDir(koshoDir)
		if err != nil {
			return fmt.Errorf("failed to read .kosho directory: %w", err)
		}

		if len(entries) == 0 {
			fmt.Println("No kosho worktrees found")
			return nil
		}

		// Get current branch of main repo
		currentBranchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		currentBranchCmd.Dir = repoRoot
		currentBranchOutput, err := currentBranchCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}
		currentBranch := strings.TrimSpace(string(currentBranchOutput))

		tbl := table.New("NAME", "STATUS", "REF", "MERGE-READY")

		for _, entry := range entries {
			if entry.IsDir() {
				name := entry.Name()
				kw := internal.NewKoshoWorktree(repoRoot, name)

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

				tbl.AddRow(name, gitStatus, gitRef, mergeReady)
			}
		}

		tbl.Print()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
