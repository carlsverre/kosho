package worktree

import (
	"fmt"
	"path/filepath"
)

// KoshoWorktree represents a git worktree managed by Kosho
type KoshoWorktree struct {
	// RepoPath is the absolute path to the git repository root
	RepoPath string
	// WorktreeName is the name of the worktree directory
	WorktreeName string
}

// NewKoshoWorktree creates a new KoshoWorktree instance
func NewKoshoWorktree(repoPath, worktreeName string) *KoshoWorktree {
	return &KoshoWorktree{
		RepoPath:     repoPath,
		WorktreeName: worktreeName,
	}
}

// WorktreePath returns the full path to the worktree directory
func (kw *KoshoWorktree) WorktreePath() string {
	return filepath.Join(kw.RepoPath, ".kosho", kw.WorktreeName)
}

// ContainerName returns the Docker container name for this worktree
func (kw *KoshoWorktree) ContainerName() string {
	repoName := filepath.Base(kw.RepoPath)
	return fmt.Sprintf("%s-%s", repoName, kw.WorktreeName)
}

// ConfigVolumeName returns the Docker volume name for config storage
func (kw *KoshoWorktree) ConfigVolumeName() string {
	repoName := filepath.Base(kw.RepoPath)
	return fmt.Sprintf("%s-%s-config", repoName, kw.WorktreeName)
}

// HistoryVolumeName returns the Docker volume name for history storage
func (kw *KoshoWorktree) HistoryVolumeName() string {
	repoName := filepath.Base(kw.RepoPath)
	return fmt.Sprintf("%s-%s-history", repoName, kw.WorktreeName)
}

// RepoName returns the repository name (basename of repo path)
func (kw *KoshoWorktree) RepoName() string {
	return filepath.Base(kw.RepoPath)
}

// KoshoDir returns the path to the .kosho directory
func (kw *KoshoWorktree) KoshoDir() string {
	return filepath.Join(kw.RepoPath, ".kosho")
}