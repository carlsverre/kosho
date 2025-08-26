package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/carlsverre/kosho/internal"

	"github.com/spf13/cobra"
)

var (
	branchFlag string
)

func checkOpenArgs(cmd *cobra.Command, args []string) error {
	args, _ = internal.SplitArgs(cmd, args)
	if len(args) < 1 {
		return fmt.Errorf("NAME argument is required")
	}
	if len(args) > 2 {
		return fmt.Errorf("too many arguments, expected at most 2 (NAME and commitish)")
	}
	return nil
}

var openCmd = &cobra.Command{
	Use:   "open [flags] [NAME] [commitish] [-- command...]",
	Short: "Open or create a kosho worktree",
	Long: `Open or create a kosho worktree in .kosho/NAME at the repo root.
If the worktree doesn't exist, it will be created.
By default, opens a new shell instance in the worktree.
If a command is provided after --, runs that command instead.`,
	Args:              checkOpenArgs,
	ValidArgsFunction: internal.WorktreeCompletionFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		args, command := internal.SplitArgs(cmd, args)
		name := args[0]

		// reserve worktree names starting with "kosho_"
		if strings.HasPrefix(name, "kosho_") {
			return fmt.Errorf("worktree names cannot start with 'kosho_'")
		}

		var commitish string
		if len(args) > 1 {
			commitish = args[1]
		}

		koshoDir, err := internal.NewKoshoDir()
		if err != nil {
			return fmt.Errorf("failed to load Kosho dir: %w", err)
		}

		kw := internal.NewKoshoWorktree(*koshoDir, name)

		// Check if worktree already exists
		if _, err := os.Stat(kw.WorktreePath()); os.IsNotExist(err) {
			// Create branch specification
			spec := internal.BranchSpec{
				BranchName: branchFlag,
				Commitish:  commitish,
			}

			// Create the worktree
			err := createWorktree(kw, spec)
			if err != nil {
				return err
			}

			// Run the create hook if it exists
			if err := internal.RunKoshoHook(kw, internal.HOOK_CREATE); err != nil {
				if remove_err := kw.Remove(true); remove_err != nil {
					return fmt.Errorf("failed to remove worktree after create hook failure: %w", remove_err)
				}
				return fmt.Errorf("failed to run create hook: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to check worktree path: %w", err)
		}

		// Run the open hook if it exists
		if err := internal.RunKoshoHook(kw, internal.HOOK_OPEN); err != nil {
			return fmt.Errorf("failed to run open hook: %w", err)
		}

		// Change to worktree directory and run shell or command
		return runInWorktree(kw, command)
	},
}

func createWorktree(kw *internal.KoshoWorktree, spec internal.BranchSpec) error {
	fmt.Printf("Creating worktree '%s'\n", kw.Name())

	// Create the worktree
	err := kw.CreateIfNotExists(spec)
	if err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}

	fmt.Printf("Worktree created successfully at %s\n", kw.WorktreePath())
	return nil
}

func runInWorktree(kw *internal.KoshoWorktree, command []string) error {
	if len(command) > 0 {
		// Run the specified command using the worktree method
		return kw.RunCommand(command)
	} else {
		// Open a new shell instance
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/bash"
		}

		// Change to the worktree directory first
		if err := os.Chdir(kw.WorktreePath()); err != nil {
			return fmt.Errorf("failed to change to worktree directory: %w", err)
		}

		// Use syscall.Exec to replace the current process with the shell
		return syscall.Exec(shell, []string{shell}, os.Environ())
	}
}

func init() {
	rootCmd.AddCommand(openCmd)
	openCmd.Flags().StringVarP(&branchFlag, "branch", "b", "", "specify the name of the new branch")
}
