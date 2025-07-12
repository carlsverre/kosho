package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/spf13/cobra"
	"kosho/internal/docker"
	"kosho/internal/git"
)

func generateWorktreeName() string {
	return namesgenerator.GetRandomName(0)
}

var worktreeCmd = &cobra.Command{
	Use:   "worktree [branch] [name]",
	Short: "Create a new git worktree",
	Long: `Create a new git worktree in the configured location. 
If no name is provided, a Docker-style name will be generated.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		branch := args[0]
		
		var name string
		if len(args) > 1 {
			name = args[1]
		} else {
			name = generateWorktreeName()
		}
		
		// Get current directory and find git root
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		
		repoRoot, err := git.FindGitRoot(currentDir)
		if err != nil {
			return fmt.Errorf("failed to find git repository: %w", err)
		}
		
		// Worktree directory is .worktrees in repo root
		worktreeDir := filepath.Join(repoRoot, ".worktrees")
		worktreePath := filepath.Join(worktreeDir, name)
		
		// Ensure .worktrees is in .gitignore
		err = git.EnsureWorktreesInGitignore(repoRoot)
		if err != nil {
			return fmt.Errorf("failed to update .gitignore: %w", err)
		}
		
		fmt.Printf("Creating worktree '%s' for branch '%s' in %s\n", name, branch, worktreePath)
		
		err = git.CreateWorktree(repoRoot, name, branch, worktreeDir)
		if err != nil {
			return fmt.Errorf("failed to create worktree: %w", err)
		}
		
		fmt.Printf("Worktree created successfully at %s\n", worktreePath)
		
		// Start interactive bash shell in the worktree
		fmt.Println()
		err = docker.StartInteractiveShell(worktreePath)
		if err != nil {
			return fmt.Errorf("failed to start shell: %w", err)
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(worktreeCmd)
}