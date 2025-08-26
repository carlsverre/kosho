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

		fmt.Println(info.Main.Version)

		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				fmt.Println(setting.Value)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
