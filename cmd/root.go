package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kosho",
	Short: "A CLI tool for managing git worktrees",
	Long: `Kosho creates and manages git worktrees in .kosho/ directories and helps
launch tools within them for isolated development environments.`,
}

func Execute() error {
	return rootCmd.Execute()
}
