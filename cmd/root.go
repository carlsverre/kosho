package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kosho",
	Short: "A CLI tool for creating git worktrees with Docker development environments",
	Long: `Kosho creates git worktrees in pre-configured locations and launches 
interactive Docker development environments with proper volume mounts.`,
}

func Execute() error {
	return rootCmd.Execute()
}