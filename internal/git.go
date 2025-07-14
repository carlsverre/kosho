package internal

import (
	"bufio"
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

func EnsureGitIgnored(glob string) error {
	repoRoot, err := FindGitRoot()
	if err != nil {
		return fmt.Errorf("failed to find git repository: %w", err)
	}
	gitignorePath := filepath.Join(repoRoot, ".gitignore")

	// Check if .gitignore exists and if /.kosho is already in it
	if file, err := os.Open(gitignorePath); err == nil {
		defer func() { _ = file.Close() }()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == glob {
				return nil // Already present
			}
		}
	}

	// Update .gitignore
	file, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open .gitignore: %w", err)
	}
	defer func() { _ = file.Close() }()

	_, err = fmt.Fprintf(file, "%s\n", glob)
	if err != nil {
		return fmt.Errorf("failed to write to .gitignore: %w", err)
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
