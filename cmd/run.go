package cmd

import (
	"fmt"
	"strings"

	"github.com/carlsverre/kosho/internal"

	"github.com/spf13/cobra"
)

func hasHelp(args []string) bool {
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			break
		}
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}

func checkRunArgs(cmd *cobra.Command, args []string) error {
	if hasHelp(args) {
		return nil
	}

	if len(args) == 0 {
		return fmt.Errorf("BRANCH argument is required")
	}
	if len(args) == 1 {
		return fmt.Errorf("command is required")
	}
	return nil
}

var runCmd = &cobra.Command{
	Use:   "run BRANCH COMMAND [args...]",
	Short: "Runs COMMAND in a Git worktree checked out to BRANCH",
	Long: `Runs COMMAND in a Git worktree located at .kosho/BRANCH.
If the worktree or branch doesn't exist, it will be created. Any
additional arguments and flags will be passed through as-is to the
command.`,
	Example:            "kosho run bugfix pnpm build",
	Args:               checkRunArgs,
	ValidArgsFunction:  internal.RunCompletion,
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle --help explicitly when flag parsing is disabled
		if hasHelp(args) {
			return cmd.Help()
		}

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
