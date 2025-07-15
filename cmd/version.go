package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print the version information including the git SHA used to build kosho.`,
	Run: func(cmd *cobra.Command, args []string) {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			fmt.Println("kosho dev (commit unknown)")
			return
		}

		version := "dev"
		commit := "unknown"
		fullSHA, _ := cmd.Flags().GetBool("full")

		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				commit = setting.Value
				if !fullSHA && len(commit) > 7 {
					commit = commit[:7]
				}
			case "vcs.modified":
				if setting.Value == "true" {
					commit += "-dirty"
				}
			}
		}

		fmt.Printf("kosho %s (commit %s)\n", version, commit)
	},
}

func init() {
	versionCmd.Flags().BoolP("full", "f", false, "show full git SHA instead of short version")
	rootCmd.AddCommand(versionCmd)
}
