package internal

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type KoshoHook string

const (
	// Runs after a worktree is created
	HOOK_CREATE KoshoHook = "create"

	// Runs before opening a worktree
	HOOK_OPEN KoshoHook = "open"

	// Runs while merging a worktree, before validation
	HOOK_MERGE KoshoHook = "merge"

	// Runs before removing a worktree
	HOOK_REMOVE KoshoHook = "remove"
)

var (
	//go:embed sample-hooks
	KoshoHooks embed.FS
)

func writeKoshoHookSamples(hookDir string) error {
	samples, err := KoshoHooks.ReadDir("sample-hooks")
	if err != nil {
		return fmt.Errorf("failed to read sample hooks directory: %w", err)
	}
	for _, sample := range samples {
		srcPath := filepath.Join("sample-hooks", sample.Name())
		destPath := filepath.Join(hookDir, sample.Name())
		data, err := KoshoHooks.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("failed to read sample hook file %s: %w", srcPath, err)
		}
		if err := writeFileIfNotExists(destPath, data, 0755); err != nil {
			return fmt.Errorf("failed to write sample hook file %s: %w", destPath, err)
		}
	}

	return nil
}

func RunKoshoHook(worktree *KoshoWorktree, hook KoshoHook) error {
	hookFile := worktree.KoshoDir.HookPath(hook)

	// abort if hook does not exist
	if _, err := os.Stat(hookFile); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to stat hook file %s: %w", hookFile, err)
	}

	// try to run the hook in the worktree directory
	cmd := exec.Command(hookFile)
	cmd.Dir = worktree.WorktreePath()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"KOSHO_HOOK="+string(hook),
		"KOSHO_WORKTREE="+worktree.WorktreeName,
		"KOSHO_REPO="+worktree.KoshoDir.RepoPath(),
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run hook %s: %w", hook, err)
	}
	return nil
}
