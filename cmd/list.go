package cmd

import (
	"fmt"

	"github.com/carlsverre/kosho/internal"

	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all kosho worktrees",
	Long:  `List all kosho worktrees and their current git status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		koshoDir, err := internal.NewKoshoDir()
		if err != nil {
			return fmt.Errorf("failed to load Kosho dir: %w", err)
		}

		worktrees, err := koshoDir.ListWorktrees()
		if err != nil {
			return fmt.Errorf("failed to list worktrees: %w", err)
		}

		if len(worktrees) == 0 {
			fmt.Println("No kosho worktrees found")
			return nil
		}

		tbl := table.New("NAME", "UPSTREAM", "REF", "STATUS")
		for _, kw := range worktrees {
			upstream, err := kw.GetUpstream()
			if err != nil {
				upstream = "unknown"
			}

			gitRef, err := kw.GitBranch()
			if err != nil {
				gitRef = "detached"
			}

			status, err := kw.Status()
			if err != nil {
				status = "error"
			}

			tbl.AddRow(kw.Name(), upstream, gitRef, status)
		}

		tbl.Print()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
