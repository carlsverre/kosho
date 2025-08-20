package internal

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lithammer/dedent"
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

	KOSHO_DIR        = ".kosho"
	KOSHO_CONFIG_DIR = KOSHO_DIR + "/kosho_config"
	KOSHO_HOOKS_DIR  = KOSHO_CONFIG_DIR + "/hooks"
)

var (
	//go:embed sample-hooks
	KoshoHooks embed.FS

	KoshoHooksGitIgnore = []byte(strings.TrimLeft(dedent.Dedent(`
		*.sample
	`), "\n"))

	KoshoGitIgnore = []byte(strings.TrimLeft(dedent.Dedent(`
		*
		!/.gitignore
		!/kosho_config/
		!/kosho_config/**
	`), "\n"))
)

func SetupKoshoDir() error {
	repoRoot, err := FindGitRoot()
	if err != nil {
		return fmt.Errorf("failed to find git repository: %w", err)
	}

	// if the root .gitignore contains .kosho, remove it
	// this is an upgrade step from an earlier Kosho version
	rootGitIgnorePath := filepath.Join(repoRoot, ".gitignore")
	if file, err := os.Open(rootGitIgnorePath); err == nil {
		defer func() { _ = file.Close() }()
		scanner := bufio.NewScanner(file)
		var lines []string
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.Contains(line, ".kosho") {
				lines = append(lines, line)
			}
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading .gitignore: %w", err)
		}
		// Rewrite .gitignore without the .kosho line
		if err := os.WriteFile(rootGitIgnorePath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
			return fmt.Errorf("failed to rewrite .gitignore: %w", err)
		}
	}

	// Recursively create .kosho/kosho_config/hooks directory structure
	hooksDir := filepath.Join(repoRoot, KOSHO_HOOKS_DIR)
	err = os.MkdirAll(hooksDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create hooks directory %s: %w", hooksDir, err)
	}

	// Initialize the hooks directory with samples
	samples, err := KoshoHooks.ReadDir("sample-hooks")
	if err != nil {
		return fmt.Errorf("failed to read sample hooks directory: %w", err)
	}
	for _, sample := range samples {
		srcPath := filepath.Join("sample-hooks", sample.Name())
		destPath := filepath.Join(hooksDir, sample.Name())
		// Read the sample hook file
		data, err := KoshoHooks.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("failed to read sample hook file %s: %w", srcPath, err)
		}
		// Write the sample hook file to the hooks directory
		if err := os.WriteFile(destPath, data, 0755); err != nil {
			return fmt.Errorf("failed to write sample hook file %s: %w", destPath, err)
		}
	}

	// Create .kosho/kosho_config/hooks/.gitignore if it doesn't exist
	koshoHooksGitIgnorePath := filepath.Join(hooksDir, ".gitignore")
	if _, err := os.Stat(koshoHooksGitIgnorePath); os.IsNotExist(err) {
		if err := os.WriteFile(koshoHooksGitIgnorePath, KoshoHooksGitIgnore, 0644); err != nil {
			return fmt.Errorf("failed to create %s: %w", koshoHooksGitIgnorePath, err)
		}
	}

	// Create .kosho/.gitignore if it doesn't exist
	koshoGitIgnorePath := filepath.Join(repoRoot, KOSHO_DIR, ".gitignore")
	if _, err := os.Stat(koshoGitIgnorePath); os.IsNotExist(err) {
		if err := os.WriteFile(koshoGitIgnorePath, KoshoGitIgnore, 0644); err != nil {
			return fmt.Errorf("failed to create %s: %w", koshoGitIgnorePath, err)
		}
	}

	return nil
}

func RunKoshoHook(hook KoshoHook, worktree *KoshoWorktree) error {
	repoRoot, err := FindGitRoot()
	if err != nil {
		return fmt.Errorf("failed to find git repository: %w", err)
	}

	hookFile := filepath.Join(repoRoot, KOSHO_HOOKS_DIR, string(hook))

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
	cmd.Env = append(cmd.Env, "KOSHO_HOOK="+string(hook))
	cmd.Env = append(cmd.Env, "KOSHO_WORKTREE="+worktree.WorktreeName)
	cmd.Env = append(cmd.Env, "KOSHO_REPO="+repoRoot)

	if err = cmd.Run(); err != nil {
		return fmt.Errorf("failed to run hook %s: %w", hook, err)
	}
	return nil
}
