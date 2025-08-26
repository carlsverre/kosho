package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lithammer/dedent"
)

const (
	KOSHO_DIR          = ".kosho"
	KOSHO_HOOKS_DIR    = "hooks"
	KOSHO_WORKTREE_DIR = "worktrees"
)

var (
	KoshoGitIgnore = []byte(strings.TrimLeft(dedent.Dedent(`
		/worktrees/
		/worktrees/**
		/hooks/*.sample
	`), "\n"))
)

type KoshoDir struct {
	repoPath string
}

func NewKoshoDir() (*KoshoDir, error) {
	repoPath, err := FindGitRoot()
	if err != nil {
		return nil, err
	}
	err = setupKoshoRepo(repoPath)
	if err != nil {
		return nil, err
	}
	return &KoshoDir{repoPath: repoPath}, nil
}

func setupKoshoRepo(repoDir string) error {
	// if the root .gitignore contains .kosho, remove it
	// this is an upgrade step from an earlier Kosho version
	rootGitIgnorePath := filepath.Join(repoDir, ".gitignore")
	if err := RemoveLinesFromGitIgnore(rootGitIgnorePath, ".kosho"); err != nil {
		// Only return error if it's not a "file doesn't exist" error
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove .kosho from .gitignore: %w", err)
		}
	}

	// Create .kosho directory structure
	koshoDir := filepath.Join(repoDir, KOSHO_DIR)
	dirs := []string{KOSHO_HOOKS_DIR, KOSHO_WORKTREE_DIR}
	for _, dir := range dirs {
		dirPath := filepath.Join(koshoDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create %s directory: %w", dirPath, err)
		}
	}

	// Initialize the hooks directory with samples
	if err := writeKoshoHookSamples(filepath.Join(koshoDir, KOSHO_HOOKS_DIR)); err != nil {
		return fmt.Errorf("failed to initialize hooks directory: %w", err)
	}

	// Create .kosho/.gitignore if it doesn't exist
	koshoGitIgnorePath := filepath.Join(koshoDir, ".gitignore")
	if err := writeFileIfNotExists(koshoGitIgnorePath, KoshoGitIgnore, 0644); err != nil {
		return fmt.Errorf("failed to create %s: %w", koshoGitIgnorePath, err)
	}

	return nil
}

func writeFileIfNotExists(path string, contents []byte, perm os.FileMode) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.WriteFile(path, contents, perm); err != nil {
			return fmt.Errorf("failed to create %s: %w", path, err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check if %s exists: %w", path, err)
	}
	return nil
}

func (kr *KoshoDir) RepoPath() string {
	return kr.repoPath
}

func (kr *KoshoDir) WorktreePath(worktreeName string) string {
	return filepath.Join(kr.repoPath, KOSHO_DIR, KOSHO_WORKTREE_DIR, worktreeName)
}

func (kr *KoshoDir) HookPath(hook KoshoHook) string {
	return filepath.Join(kr.repoPath, KOSHO_DIR, KOSHO_HOOKS_DIR, string(hook))
}

func (kr *KoshoDir) ListWorktrees() ([]KoshoWorktree, error) {
	entries, err := os.ReadDir(filepath.Join(kr.repoPath, KOSHO_DIR, KOSHO_WORKTREE_DIR))
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	worktrees := make([]KoshoWorktree, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		worktrees = append(worktrees, *NewKoshoWorktree(*kr, entry.Name()))
	}
	return worktrees, nil
}
