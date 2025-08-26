package internal

import (
	"github.com/spf13/cobra"
)

// SplitArgs splits command arguments at the dash separator, returning arguments before and after the dash.
func SplitArgs(cmd *cobra.Command, args []string) ([]string, []string) {
	argsLenAtDash := cmd.ArgsLenAtDash()
	if argsLenAtDash < 0 {
		argsLenAtDash = len(args)
	}
	return args[:argsLenAtDash], args[argsLenAtDash:]
}
