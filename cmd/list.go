package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"kosho/internal/git"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all kosho worktrees",
	Long:  `List all kosho worktrees, their git status, and container status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get current directory and find git root
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		repoRoot, err := git.FindGitRoot(currentDir)
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
				// Stub: In real implementation, would check git status and container status
				fmt.Printf("%s\t\t[stub]\t\t[stub]\n", name)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
