# Kosho Merge Command Implementation

## Command: `kosho merge <worktree> [-- git-merge-args...]`

Merges worktree branch into current main repo branch.

## Implementation Steps

### 1. Add to internal package:

- `GetWorktreeBranch(name) (string, error)` - Get branch name from worktree
- `IsAncestor(ancestor, descendant) (bool, error)` - Check git ancestry

### 2. Create cmd/merge.go:

- Parse worktree name and optional git merge args (split on `--`)
- Validate worktree exists and is clean
- Check ancestry: current branch is ancestor of worktree branch
- Execute: `git merge [args...] <worktree-branch>`

## Key Implementation Details

### Git Commands Used

```bash
# Get worktree branch
git -C .kosho/<name> rev-parse --abbrev-ref HEAD

# Check ancestry (exit code 0 = is ancestor)
git merge-base --is-ancestor <current-branch> <worktree-branch>

# Execute merge in main repo
git merge [args...] <worktree-branch>
```

### Error Cases

- Worktree doesn't exist → fail early
- Worktree is dirty → abort with helpful message
- Current branch not ancestor → abort with git log suggestion
- Git merge failure → pass through git's error handling

### Argument Parsing

Use Cobra's `ArgsLenAtDash()` like the open command:

- Before `--`: worktree name (required)
- After `--`: arguments passed to git merge
