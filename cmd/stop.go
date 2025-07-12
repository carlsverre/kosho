package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [NAME]",
	Short: "Stop a kosho container",
	Long:  `Stop a running kosho container. If NAME is not provided, tries to determine from current directory.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		if len(args) > 0 {
			name = args[0]
		}

		if name == "" {
			// In real implementation, would try to determine from current directory
			return fmt.Errorf("NAME is required when not in a kosho worktree directory")
		}

		// Stub implementation
		fmt.Printf("Stopping container for worktree '%s'\n", name)
		fmt.Println("Note: Stop command is stubbed - would stop Docker container")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
