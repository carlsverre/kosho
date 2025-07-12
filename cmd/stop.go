package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"kosho/internal/docker"
	"kosho/internal/git"
	"kosho/internal/worktree"
)

var stopCmd = &cobra.Command{
	Use:   "stop [NAME]",
	Short: "Stop a kosho container",
	Long:  `Stop a running kosho container. If NAME is not provided, tries to determine from current directory.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		if len(args) > 0 {
			name = args[0]
		}

		return stopWorktree(name)
	},
}

func stopWorktree(name string) error {
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

	// Stop the container
	fmt.Printf("Stopping container for worktree '%s'\n", name)
	return docker.StopContainer(kw.ContainerName())
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
