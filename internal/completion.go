package internal

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// RunCompletion provides autocompletion for `kosho run`
func RunCompletion(cmd *cobra.Command, args []string, prefix string) ([]string, cobra.ShellCompDirective) {
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
	} else if len(args) == 1 {
		// completing the command name - return executables from PATH
		commands := getExecutablesFromPath(prefix)
		return commands, cobra.ShellCompDirectiveNoFileComp
	}

	// completing arguments to the command - use default completion
	return nil, cobra.ShellCompDirectiveDefault
}

// getExecutablesFromPath returns a list of executables from PATH that match the given prefix
func getExecutablesFromPath(prefix string) []string {
	if len(prefix) == 0 {
		return nil
	}

	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return nil
	}

	seen := make(map[string]bool)
	var results []string

	for _, dir := range filepath.SplitList(pathEnv) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			name := entry.Name()

			// Skip if already seen or doesn't match prefix
			if seen[name] || (prefix != "" && !strings.HasPrefix(name, prefix)) {
				continue
			}

			// Check if executable
			fullPath := filepath.Join(dir, name)
			if isExecutable(fullPath) {
				results = append(results, name)
				seen[name] = true
			}
		}
	}

	return results
}

// isExecutable checks if a file is executable
func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	// Check if it's a regular file and has execute permission
	return info.Mode().IsRegular() && info.Mode().Perm()&0111 != 0
}
