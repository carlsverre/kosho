package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FindGitRoot finds the root directory of the git repository, handling both regular repos and worktrees.
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

// RemoveLinesFromGitIgnore removes lines containing the specified substring from a .gitignore file,
// preserving the original formatting including trailing newlines. Only writes if changes are needed.
func RemoveLinesFromGitIgnore(gitIgnorePath, substring string) error {
	// Read the entire file to preserve original format
	originalContent, err := os.ReadFile(gitIgnorePath)
	if err != nil {
		return err
	}

	// Process lines directly from the content
	contentLines := strings.Split(string(originalContent), "\n")

	var lines []string
	var foundMatch bool
	for _, line := range contentLines {
		if strings.Contains(line, substring) {
			foundMatch = true
		} else {
			lines = append(lines, line)
		}
	}

	// Only write if we found lines to remove
	if !foundMatch {
		return nil
	}

	// Rewrite .gitignore
	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(gitIgnorePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to rewrite .gitignore: %w", err)
	}

	return nil
}
