package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// BranchSpec encapsulates branch creation parameters for git worktree
type BranchSpec struct {
	// BranchName is the name of the branch to create or checkout
	BranchName string
	// Commitish is the commit-ish to base the branch on (optional)
	Commitish string
	// Reset indicates whether to reset an existing branch to the target commitish
	Reset bool
}

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

// KoshoDir returns the path to the .kosho directory
func (kw *KoshoWorktree) KoshoDir() string {
	return filepath.Join(kw.RepoPath, ".kosho")
}

func (kw *KoshoWorktree) CreateIfNotExists(spec BranchSpec) error {
	worktreePath := kw.WorktreePath()

	if err := os.MkdirAll(worktreePath, 0755); err != nil {
		return fmt.Errorf("failed to create worktree directory: %w", err)
	}

	// Build git worktree add command
	args := []string{"worktree", "add"}

	// Add branch flags if branch name is specified
	if spec.BranchName != "" {
		if spec.Reset {
			args = append(args, "-B", spec.BranchName)
		} else {
			args = append(args, "-b", spec.BranchName)
		}
	}

	args = append(args, worktreePath)

	// Add commitish if specified
	if spec.Commitish != "" {
		args = append(args, spec.Commitish)
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = kw.RepoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create worktree: %w\nOutput: %s", err, string(output))
	}

	return nil
}
