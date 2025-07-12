package cmd

import (
	"fmt"
	"kosho/internal"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

		fmt.Println("Kosho Worktrees:")
		fmt.Println("NAME\t\tSTATUS\t\tREF")
		fmt.Println("----\t\t------\t\t---")

		for _, entry := range entries {
			if entry.IsDir() {
				name := entry.Name()
				kw := internal.NewKoshoWorktree(repoRoot, name)

				// Check git status
				gitStatus := getWorktreeStatus(kw.WorktreePath())

				// Get current branch/ref
				gitRef := getWorktreeRef(kw.WorktreePath())

				fmt.Printf("%s\t\t%s\t\t%s\n", name, gitStatus, gitRef)
			}
		}

		return nil
	},
}

func getWorktreeStatus(worktreePath string) string {
	// Check if there are uncommitted changes
	gitCmd := exec.Command("git", "status", "--porcelain")
	gitCmd.Dir = worktreePath

	output, err := gitCmd.CombinedOutput()
	if err != nil {
		return "error"
	}

	if len(strings.TrimSpace(string(output))) == 0 {
		return "clean"
	}
	return "dirty"
}

func getWorktreeRef(worktreePath string) string {
	// Get current branch name
	gitCmd := exec.Command("git", "branch", "--show-current")
	gitCmd.Dir = worktreePath

	output, err := gitCmd.CombinedOutput()
	if err != nil {
		// Try to get HEAD ref instead
		gitCmd = exec.Command("git", "rev-parse", "--short", "HEAD")
		gitCmd.Dir = worktreePath
		output, err = gitCmd.CombinedOutput()
		if err != nil {
			return "unknown"
		}
	}

	return strings.TrimSpace(string(output))
}

func init() {
	rootCmd.AddCommand(listCmd)
}
