package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	verbose bool
	config  string
	force   bool
	timeout int
)

var rootCmd = &cobra.Command{
	Use:   "kosho",
	Short: "Docker image and container management tool",
}

var dummyCmd = &cobra.Command{
	Use:     "dummy [image]",
	Aliases: []string{"d"},
	Short:   "Dummy subcommand placeholder with image argument",
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Global flags - Verbose: %v, Config: %s\n", verbose, config)
		fmt.Printf("Command flags - Force: %v, Timeout: %d\n", force, timeout)
		
		if len(args) > 0 {
			fmt.Printf("Arguments: %v\n", args)
		} else {
			fmt.Println("No arguments provided")
		}
		
		fmt.Println("This is a dummy subcommand. Implementation coming soon!")
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "~/.kosho.yaml", "config file path")
	
	dummyCmd.Flags().BoolVarP(&force, "force", "f", false, "force operation")
	dummyCmd.Flags().IntVarP(&timeout, "timeout", "t", 30, "timeout in seconds")
	
	rootCmd.AddCommand(dummyCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
