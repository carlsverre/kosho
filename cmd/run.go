package cmd

import (
	"fmt"

	"github.com/carlsverre/kosho/internal"

	"github.com/spf13/cobra"
)

func checkRunArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("BRANCH argument is required")
	}
	if len(args) == 1 {
		return fmt.Errorf("command is required")
	}
	return nil
}

var runCmd = &cobra.Command{
	Use:   "run [BRANCH] [command...]",
	Short: "Runs a command in a Git worktree on a specific branch",
	Long: `Runs a command in a Git worktree located at .kosho/BRANCH.
	       If the worktree or branch doesn't exist, it will be created.`,
	Args:              checkRunArgs,
	ValidArgsFunction: internal.RunCompletion,
	RunE: func(cmd *cobra.Command, args []string) error {
		branch, rest := args[0], args[1:]

		koshoDir, err := internal.LoadKoshoDir()
		if err != nil {
			return fmt.Errorf("failed to load Kosho dir: %w", err)
		}

		kw := internal.NewKoshoWorktree(*koshoDir, branch)

		createdWorktree := false

		// Check if worktree already exists
		if exists, err := kw.Exists(); !exists {
			if err := createWorktree(kw); err != nil {
				return err
			}
			if err := runHook(kw, internal.HOOK_CREATE, true); err != nil {
				return err
			}
			createdWorktree = true
		} else if err != nil {
			return fmt.Errorf("failed to check worktree path: %w", err)
		}

		// Run the run hook if it exists
		if err := runHook(kw, internal.HOOK_RUN, createdWorktree, fmt.Sprintf("KOSHO_CMD=%q", rest[0])); err != nil {
			return err
		}

		return kw.RunCommand(rest)
	},
}

func runHook(kw *internal.KoshoWorktree, hook internal.KoshoHook, deleteWorktreeOnFailure bool, extraEnv ...string) error {
	if err := internal.RunKoshoHook(kw, hook, extraEnv...); err != nil {
		fmt.Printf("Failed to run hook `%s`\n", hook)
		if deleteWorktreeOnFailure {
			fmt.Printf("cleaning up worktree '%s'... ", kw.Name())
			if remove_err := kw.Remove(true); remove_err != nil {
				fmt.Printf("ERROR\n")
				return fmt.Errorf("failed to remove worktree after '%s' hook failure: %w", hook, remove_err)
			}
			fmt.Printf("DONE\n")
		}
		return err
	}
	return nil
}

func createWorktree(kw *internal.KoshoWorktree) error {
	fmt.Printf("Creating worktree '%s'... ", kw.Name())

	// Create the worktree
	err := kw.CreateWorktree()
	if err != nil {
		fmt.Printf("ERROR\n")
		return fmt.Errorf("failed to create worktree: %w", err)
	}

	fmt.Printf("DONE\n")
	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)
}
