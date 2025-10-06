package internal

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// KoshoWorktree represents a git worktree managed by Kosho
type KoshoWorktree struct {
	KoshoDir     KoshoDir
	BranchName   string
	WorktreeName string
}

// NewKoshoWorktree creates a new KoshoWorktree instance
func NewKoshoWorktree(root KoshoDir, branchName string) *KoshoWorktree {
	return &KoshoWorktree{
		KoshoDir:     root,
		BranchName:   branchName,
		WorktreeName: sluggify(branchName),
	}
}

func (kw *KoshoWorktree) Name() string {
	return kw.WorktreeName
}

// WorktreePath returns the full path to the worktree directory
func (kw *KoshoWorktree) WorktreePath() string {
	return kw.KoshoDir.WorktreePath(kw.WorktreeName)
}

func (kw *KoshoWorktree) Exists() (bool, error) {
	if _, err := os.Stat(kw.WorktreePath()); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to check worktree path: %w", err)
	}
	return true, nil
}

func (kw *KoshoWorktree) CreateWorktree() error {
	worktreePath := kw.WorktreePath()

	args := []string{"worktree", "add"}
	if !BranchExists(kw.KoshoDir.repoPath, kw.BranchName) {
		args = append(args, "-b", kw.BranchName, worktreePath)
	} else if kw.BranchName != kw.WorktreeName {
		args = append(args, worktreePath, kw.BranchName)
	} else {
		args = append(args, worktreePath)
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = kw.KoshoDir.RepoPath()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create worktree: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Remove the worktree if it's clean, but leaves the branch as is.
// force will cause the worktree to be removed even if it's dirty
func (kw *KoshoWorktree) Remove(force bool) error {
	// Build git worktree remove command
	args := []string{"worktree", "remove", kw.WorktreePath()}
	if force {
		args = append(args, "--force")
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = kw.KoshoDir.RepoPath()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove worktree: %w\nOutput: %s", err, string(output))
	}

	return nil
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
	cmd.Env = os.Environ()

	return cmd.Run()
}

// AheadBehind returns the number of commits ahead and behind the upstream
func (kw *KoshoWorktree) AheadBehind() (int, int, error) {
	cmd := exec.Command("git", "rev-list", "--left-right", "--count", "@{u}...HEAD")
	cmd.Dir = kw.WorktreePath()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to run git rev-list: %w", err)
	}

	counts := strings.Fields(strings.TrimSpace(string(output)))
	if len(counts) != 2 {
		return 0, 0, fmt.Errorf("unexpected output from git rev-list: %s", string(output))
	}

	behind, err := strconv.Atoi(counts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse rev-list output: %w", err)
	}
	ahead, err := strconv.Atoi(counts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse rev-list output: %w", err)
	}
	return ahead, behind, nil
}

// Status returns a string describing the worktree's current status relative to its base
func (kw *KoshoWorktree) Status() (string, error) {
	var statusParts []string

	upstream, err := kw.GetUpstream()
	if err != nil {
		return "", err
	}

	if upstream != "" {
		ahead, behind, err := kw.AheadBehind()
		if err != nil {
			return "", err
		}
		if ahead != 0 {
			statusParts = append(statusParts, fmt.Sprintf("ahead %d", ahead))
		}
		if behind != 0 {
			statusParts = append(statusParts, fmt.Sprintf("behind %d", behind))
		}
	}

	isDirty, err := kw.IsDirty()
	if err != nil {
		return "", fmt.Errorf("failed to check if worktree is dirty: %w", err)
	}
	if isDirty {
		statusParts = append(statusParts, "(dirty)")
	} else {
		statusParts = append(statusParts, "(clean)")
	}

	return strings.Join(statusParts, " "), nil
}

// IsClean checks if the worktree is not dirty
func (kw *KoshoWorktree) IsClean() (bool, error) {
	isDirty, err := kw.IsDirty()
	if err != nil {
		return false, fmt.Errorf("failed to check if worktree is dirty: %w", err)
	}
	return !isDirty, nil
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

// GetUpstream returns the upstream branch name if one exists
func (kw *KoshoWorktree) GetUpstream() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "@{u}")
	cmd.Dir = kw.WorktreePath()

	output, err := cmd.CombinedOutput()
	if err != nil {
		// No upstream configured
		return "", nil
	}

	upstream := strings.TrimSpace(string(output))
	return upstream, nil
}

var slugRegex = regexp.MustCompile(`[^\p{L}\p{N}]+`)

// returns a version of the string with the following changes:
// - lowercased
// - replace all punctuation, symbols, whitespace with `-`
func sluggify(s string) string {
	s = strings.ToLower(s)
	s = slugRegex.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
