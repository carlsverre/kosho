package internal

import (
	"github.com/spf13/cobra"
)

// RunCompletion provides autocompletion for `kosho run`
func RunCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	repoRoot, err := FindGitRoot()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError

	}
	if len(args) < 1 {
		// completing the branch name
		branches, err := ListBranches(repoRoot)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return branches, cobra.ShellCompDirectiveNoFileComp
	} else {
		// completing the shell command to run
	}

	return nil, cobra.ShellCompDirectiveDefault
}
