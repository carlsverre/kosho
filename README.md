<p align="center">
  <h1 align="center">Kōshō</h1>
  <p align="center"><em>A loyal assistant, trained for quiet efficiency.</em></p>
</p>

---

**Kosho** creates and manages [git worktree]s in `.kosho/` directories, making it easy to work on multiple branches simultaneously without the overhead of cloning repositories or switching contexts.

```bash
go install github.com/carlsverre/kosho
```

## Why Use Kosho?

Kosho is well-suited for **running multiple concurrent AI coding agents** like Claude Code, each working on different features or branches in separate workspaces:

- **🤖 Multiple AI Agents**: Run Claude Code (or other AI tools) in separate worktrees simultaneously
- **⚡ Reduced Context Switching**: Work on multiple branches/features at once without git checkout delays
- **🔒 Workspace Isolation**: Each agent operates in its own [git worktree] with independent working state
- **🚀 More Efficient Than Cloning**: Worktrees share the same `.git` directory - no duplicate repositories
- **📋 Easy Coordination**: Use `kosho list` to see what each agent is working on
- **🧹 Easy Cleanup**: Remove completed work environments without affecting your main repository

**Example workflow using Claude Code:**

```bash
# Start Claude Code on feature branch
kosho open feature-auth -b feature/user-auth -- claude

# Start another agent on bug fixes in a separate shell
kosho open bugfix-session -- claude

# Check what each agent is working on
kosho list

# Merge completed work back to main branch
kosho merge feature-auth -- --squash
kosho merge bugfix-session --remove

# Clean up merged worktrees
kosho remove feature-auth
```

## Basic Usage

1. **Create and open a worktree**:

   ```bash
   kosho open my-feature
   ```

   This creates a new worktree at `.kosho/my-feature`, checks out or creates the branch called `my-feature`, and opens a shell session.

2. **Run a command in a worktree instead of opening a shell**:

   ```bash
   kosho open my-feature -- claude
   ```

   Start a Claude Code agent in the new worktree.

3. **List all worktrees**:

   ```bash
   kosho list
   ```

4. **Remove a worktree**:

   ```bash
   kosho remove my-feature
   ```

## Commands

### `kosho open [flags] NAME [commitish] [-- command...]`

Creates or opens a worktree. If the worktree doesn't exist, it will be created.

**Arguments:**

- `NAME` (required): Name of the worktree
- `commitish` (optional): Git commit-ish to base the worktree on

**Flags:**

- `-b, --branch <name>`: Create a new branch
- `-B, --reset-branch <name>`: Create or reset a branch to the target commitish

**Examples:**

```bash
# Create worktree from current HEAD
kosho open bugfix

# Create worktree with new branch
kosho open feature-work -b feature/awesome-feature

# Create worktree from specific commit
kosho open hotfix v1.2.3

# Run command instead of opening shell
kosho open testing -- npm test

# Reset existing branch to specific commit
kosho open release -B release/v2.0 v2.0.0
```

### `kosho list`

Lists all kosho worktrees with their status and current git reference.

**Output:**

```
NAME        STATUS  REF
bugfix      clean   main
hotfix      dirty   hotfix
my-feature  dirty   feature/my-feature
```

- **STATUS**: `clean` (no uncommitted changes) or `dirty` (has uncommitted changes)
- **REF**: Current branch name or commit hash

### `kosho remove [flags] NAME`

Removes a worktree and cleans up git references.

**Flags:**

- `-f, --force`: Force removal even if worktree has uncommitted changes

**Examples:**

```bash
# Remove clean worktree
kosho remove my-feature

# Force remove dirty worktree
kosho remove my-feature --force
```

### `kosho merge [worktree] [-- git-merge-args...]`

Merges a worktree branch into the current branch of the main repository.

**Requirements:**

- Worktree must be clean (no uncommitted changes)
- Current branch must be an ancestor of the worktree branch

**Examples:**

```bash
# Standard merge
kosho merge feature-auth

# Squash merge
kosho merge feature-auth -- --squash

# No-fast-forward merge with message
kosho merge feature-auth -- --no-ff -m "Add authentication feature"
```

### `kosho prune`

Cleans up any dangling worktree references using `git worktree prune`.

## How It Works

Kosho manages [git worktree]s in a `.kosho/` directory at your repository root:

```
your-repo/
├── .git/
├── .kosho/           # Kosho worktrees directory (auto-added to .gitignore)
│   ├── feature-a/    # Worktree for feature-a
│   ├── bugfix/       # Worktree for bugfix
│   └── experiment/   # Worktree for experiment
├── src/
└── README.md
```

Each worktree is a complete working directory that shares the same git history but can have different branches checked out and different working states.

## Shell Completion

Kosho supports tab completion for bash, zsh, fish, and PowerShell. To set up completion for your shell:

```bash
kosho completion <shell> --help
```

This command outputs detailed setup instructions specific to your shell.

## Tips

- The `.kosho` directory is automatically added to your `.gitignore`
- Each worktree maintains its own working state and uncommitted changes
- Use `kosho list` to see the status of all your worktrees at a glance
- Use `kosho prune` periodically to clean up any orphaned references

[git worktree]: https://git-scm.com/docs/git-worktree
