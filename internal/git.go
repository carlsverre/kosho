package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func FindGitRoot() (string, error) {
	// First check if we're in a worktree by examining the git directory
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("not a git repository (or any of the parent directories): %w", err)
	}

	gitDir := strings.TrimSpace(string(output))

	// Check if we're in a worktree (git-dir contains .git/worktrees/)
	if strings.Contains(gitDir, ".git/worktrees/") {
		// Extract main repo path: /path/to/repo/.git/worktrees/name -> /path/to/repo/.git -> /path/to/repo
		mainGitDir := strings.Split(gitDir, "/worktrees/")[0]
		mainRepoRoot := filepath.Dir(mainGitDir)
		return mainRepoRoot, nil
	}

	// Not in a worktree, use standard method
	cmd = exec.Command("git", "rev-parse", "--show-toplevel")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("not a git repository (or any of the parent directories): %w", err)
	}

	gitRoot := strings.TrimSpace(string(output))
	return gitRoot, nil
}

// EnsureKoshoGitIgnore creates a .gitignore file in the .kosho directory 
// to ignore worktrees while allowing hooks to be committed
func EnsureKoshoGitIgnore() error {
	repoRoot, err := FindGitRoot()
	if err != nil {
		return fmt.Errorf("failed to find git repository: %w", err)
	}
	
	koshoDir := filepath.Join(repoRoot, ".kosho")
	gitignorePath := filepath.Join(koshoDir, ".gitignore")

	// Create .kosho directory if it doesn't exist
	if err := os.MkdirAll(koshoDir, 0755); err != nil {
		return fmt.Errorf("failed to create .kosho directory: %w", err)
	}

	// Content for .kosho/.gitignore
	gitignoreContent := `# Ignore worktree directories
*
# But keep hooks and this gitignore file
!_hooks/
!_hooks/**
!.gitignore
`

	// Write the .gitignore file
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to write .kosho/.gitignore: %w", err)
	}

	return nil
}

// IsAncestor checks if ancestorRef is an ancestor of descendantRef using git merge-base
func IsAncestor(repoPath, ancestorRef, descendantRef string) (bool, error) {
	cmd := exec.Command("git", "merge-base", "--is-ancestor", ancestorRef, descendantRef)
	cmd.Dir = repoPath

	err := cmd.Run()
	if err != nil {
		// Exit code 1 means not an ancestor, other exit codes are actual errors
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return false, nil
		}
		return false, fmt.Errorf("failed to check ancestry: %w", err)
	}

	return true, nil
}
