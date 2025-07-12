package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

func CreateWorktree(repoPath, worktreeName, branch, worktreeDir string) error {
	worktreePath := filepath.Join(worktreeDir, worktreeName)
	
	if err := os.MkdirAll(worktreeDir, 0755); err != nil {
		return fmt.Errorf("failed to create worktree directory: %w", err)
	}

	// Use git worktree add command
	// Try to create with new branch first, if it fails try existing branch
	cmd := exec.Command("git", "worktree", "add", "-b", branch, worktreePath, "HEAD")
	cmd.Dir = repoPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If branch already exists, try without -b flag
		cmd = exec.Command("git", "worktree", "add", worktreePath, branch)
		cmd.Dir = repoPath
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create worktree: %w\nOutput: %s", err, string(output))
		}
	}
	
	return nil
}

func FindGitRoot(startPath string) (string, error) {
	currentPath := startPath
	
	for {
		gitPath := filepath.Join(currentPath, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			return currentPath, nil
		}
		
		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			return "", fmt.Errorf("not a git repository (or any of the parent directories)")
		}
		currentPath = parentPath
	}
}

func GetRepoName(repoPath string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	remotes, err := repo.Remotes()
	if err != nil || len(remotes) == 0 {
		return filepath.Base(repoPath), nil
	}

	url := remotes[0].Config().URLs[0]
	name := filepath.Base(url)
	if filepath.Ext(name) == ".git" {
		name = name[:len(name)-4]
	}
	
	return name, nil
}