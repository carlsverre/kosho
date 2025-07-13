package cmd

import (
	"fmt"
	"os"
	"syscall"

	"kosho/internal"

	"github.com/spf13/cobra"
)

var (
	branchFlag      string
	resetBranchFlag string
)

func splitArgs(cmd *cobra.Command, args []string) ([]string, []string) {
	argsLenAtDash := cmd.ArgsLenAtDash()
	if argsLenAtDash < 0 {
		argsLenAtDash = len(args)
	}
	return args[:argsLenAtDash], args[argsLenAtDash:]
}

func checkArgs(cmd *cobra.Command, args []string) error {
	args, _ = splitArgs(cmd, args)
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
	Args: checkArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		args, command := splitArgs(cmd, args)
		name := args[0]

		var commitish string
		if len(args) > 1 {
			commitish = args[1]
		}

		// Find git root
		repoRoot, err := internal.FindGitRoot()
		if err != nil {
			return fmt.Errorf("failed to find git repository: %w", err)
		}

		kw := internal.NewKoshoWorktree(repoRoot, name)

		// Check if worktree already exists
		if _, err := os.Stat(kw.WorktreePath()); os.IsNotExist(err) {
			// Create branch specification
			spec := internal.BranchSpec{
				BranchName: branchFlag,
				Commitish:  commitish,
				Reset:      false,
			}

			// If resetBranchFlag is set, use it instead and set Reset to true
			if resetBranchFlag != "" {
				spec.BranchName = resetBranchFlag
				spec.Reset = true
			}

			// Create the worktree
			err := createWorktree(name, kw, spec)
			if err != nil {
				return err
			}
		} else if err != nil {
			return fmt.Errorf("failed to check worktree path: %w", err)
		}

		// Change to worktree directory and run shell or command
		return runInWorktree(kw, command)
	},
}

func createWorktree(name string, kw *internal.KoshoWorktree, spec internal.BranchSpec) error {
	// Ensure .kosho is in .gitignore
	err := internal.EnsureGitIgnored("/.kosho**")
	if err != nil {
		return fmt.Errorf("failed to update .gitignore: %w", err)
	}

	fmt.Printf("Creating worktree '%s' in %s\n", name, kw.WorktreePath())

	// Create the worktree
	err = kw.CreateIfNotExists(spec)
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
	openCmd.Flags().StringVarP(&branchFlag, "branch", "b", "", "create a new branch")
	openCmd.Flags().StringVarP(&resetBranchFlag, "reset-branch", "B", "", "create or reset a branch")
}
