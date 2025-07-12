package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"kosho/internal/docker"
	"kosho/internal/git"
	"kosho/internal/worktree"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [flags] [NAME]",
	Short: "Start a kosho worktree environment",
	Long: `Start an interactive development environment in a kosho worktree.
If NAME is not provided, tries to determine from current directory.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return startWorktree(args[0])
	},
}

func startWorktree(name string) error {
	// Get current directory and find git root
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	repoRoot, err := git.FindGitRoot(currentDir)
	if err != nil {
		return fmt.Errorf("failed to find git repository: %w", err)
	}

	// If name not provided, try to determine from current directory
	if name == "" {
		// Check if we're in a .kosho subdirectory
		koshoDir := filepath.Join(repoRoot, ".kosho")
		if rel, err := filepath.Rel(koshoDir, currentDir); err == nil {
			parts := filepath.SplitList(rel)
			if len(parts) > 0 && parts[0] != ".." {
				name = parts[0]
			}
		}

		if name == "" {
			return fmt.Errorf("NAME is required when not in a kosho worktree directory")
		}
	}

	kw := worktree.NewKoshoWorktree(repoRoot, name)

	// Check if worktree exists
	if _, err := os.Stat(kw.WorktreePath()); os.IsNotExist(err) {
		return fmt.Errorf("worktree '%s' does not exist", name)
	}

	fmt.Printf("Starting interactive environment for worktree '%s'\n", name)
	return docker.StartInteractiveShell(kw)
}

func init() {
	rootCmd.AddCommand(startCmd)
}
