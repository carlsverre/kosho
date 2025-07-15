package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func (kw *KoshoWorktree) Exists() (bool, error) {
	if _, err := os.Stat(kw.WorktreePath()); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to check worktree path: %w", err)
	}
	return true, nil
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

// Remove removes the worktree using git worktree remove
func (kw *KoshoWorktree) Remove(force bool) error {
	// Build git worktree remove command
	args := []string{"worktree", "remove", kw.WorktreePath()}
	if force {
		args = append(args, "--force")
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = kw.RepoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove worktree: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// IsDirty checks if the worktree has uncommitted changes
func (kw *KoshoWorktree) IsDirty() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = kw.WorktreePath()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to get git status: %w", err)
	}

	// If output is empty, worktree is clean
	return len(output) > 0, nil
}

// RunCommand runs a command in the worktree directory
func (kw *KoshoWorktree) RunCommand(command []string) error {
	if len(command) == 0 {
		return fmt.Errorf("no command provided")
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = kw.WorktreePath()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// GitRef returns the current git reference (branch name or commit hash)
func (kw *KoshoWorktree) GitRef() (string, error) {
	// Try to get current branch name first
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = kw.WorktreePath()

	output, err := cmd.CombinedOutput()
	if err == nil {
		branch := strings.TrimSpace(string(output))
		if branch != "" {
			return branch, nil
		}
	}

	// If no branch name, get short commit hash
	cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Dir = kw.WorktreePath()
	output, err = cmd.CombinedOutput()
	if err != nil {
		return "unknown", nil
	}

	return strings.TrimSpace(string(output)), nil
}

// HasOutstandingCommits checks if the worktree has commits ahead of the main repo's current branch
func (kw *KoshoWorktree) HasOutstandingCommits(mainRepoBranch string) (bool, error) {
	// Get worktree branch
	worktreeBranch, err := kw.GitBranch()
	if err != nil {
		return false, fmt.Errorf("failed to get worktree branch: %w", err)
	}

	// Check if there are commits in worktree branch that are not in main branch
	cmd := exec.Command("git", "rev-list", "--count", mainRepoBranch+".."+worktreeBranch)
	cmd.Dir = kw.WorktreePath()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to check outstanding commits: %w", err)
	}

	count := strings.TrimSpace(string(output))
	return count != "0", nil
}

// GitBranch returns the current branch name of the worktree
func (kw *KoshoWorktree) GitBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = kw.WorktreePath()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree branch: %w", err)
	}

	branch := strings.TrimSpace(string(output))
	if branch == "HEAD" {
		return "", fmt.Errorf("worktree is in detached HEAD state")
	}

	return branch, nil
}
