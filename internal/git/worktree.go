package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CreateKoshoWorktree(repoPath, worktreeName, worktreeDir, branchFlag, newBranchFlag, commitish string) error {
	worktreePath := filepath.Join(worktreeDir, worktreeName)

	if err := os.MkdirAll(worktreeDir, 0755); err != nil {
		return fmt.Errorf("failed to create worktree directory: %w", err)
	}

	// Build git worktree add command
	args := []string{"worktree", "add"}

	// Add branch flags if specified
	if newBranchFlag != "" {
		args = append(args, "-B", newBranchFlag)
	} else if branchFlag != "" {
		args = append(args, "-b", branchFlag)
	}

	args = append(args, worktreePath)

	// Add commitish if specified
	if commitish != "" {
		args = append(args, commitish)
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create worktree: %w\nOutput: %s", err, string(output))
	}

	return nil
}

func FindGitRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = currentDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("not a git repository (or any of the parent directories): %w", err)
	}

	gitRoot := strings.TrimSpace(string(output))
	return gitRoot, nil
}
