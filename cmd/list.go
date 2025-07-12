package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"kosho/internal/docker"
	"kosho/internal/git"
	"kosho/internal/worktree"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all kosho worktrees",
	Long:  `List all kosho worktrees, their git status, and container status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Find git root
		repoRoot, err := git.FindGitRoot()
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
		fmt.Println("NAME\t\tSTATUS\t\tCONTAINER")
		fmt.Println("----\t\t------\t\t---------")

		for _, entry := range entries {
			if entry.IsDir() {
				name := entry.Name()
				kw := worktree.NewKoshoWorktree(repoRoot, name)

				// Check git status
				gitStatus := getWorktreeStatus(kw.WorktreePath())

				// Check container status
				containerStatus := getContainerStatus(kw.ContainerName())

				fmt.Printf("%s\t\t%s\t\t%s\n", name, gitStatus, containerStatus)
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

func getContainerStatus(containerName string) string {
	running, err := docker.IsContainerRunning(containerName)
	if err != nil {
		return "error"
	}
	if running {
		return "running"
	}

	exists, err := docker.ContainerExists(containerName)
	if err != nil {
		return "error"
	}
	if exists {
		return "stopped"
	}

	return "none"
}

func init() {
	rootCmd.AddCommand(listCmd)
}
