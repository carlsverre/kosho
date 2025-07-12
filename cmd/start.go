package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"kosho/internal/docker"
	"kosho/internal/git"
)

var detachedFlag bool

var startCmd = &cobra.Command{
	Use:   "start [flags] [NAME]",
	Short: "Start a kosho worktree environment",
	Long: `Start an interactive development environment in a kosho worktree.
If NAME is not provided, tries to determine from current directory.
Use -d flag to run in detached mode.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		if len(args) > 0 {
			name = args[0]
		}

		return startWorktree(name, detachedFlag)
	},
}

func startWorktree(name string, detached bool) error {
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

	worktreePath := filepath.Join(repoRoot, ".kosho", name)

	// Check if worktree exists
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		return fmt.Errorf("worktree '%s' does not exist", name)
	}

	if detached {
		fmt.Printf("Starting detached environment for worktree '%s'\n", name)
		return docker.StartDetachedShell(worktreePath)
	} else {
		fmt.Printf("Starting interactive environment for worktree '%s'\n", name)
		return docker.StartInteractiveShell(worktreePath)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolVarP(&detachedFlag, "detached", "d", false, "run in detached mode")
}
