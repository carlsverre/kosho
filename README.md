<p align="center">
  <h1 align="center">Kōshō</h1>
  <p align="center"><em>A loyal assistant, trained for quiet efficiency.</em></p>
</p>

---

**Kosho** creates and manages [git worktree]s in `.kosho/` directories, making it easy to work on multiple branches simultaneously without the overhead of cloning repositories or switching contexts.

```bash
go install github.com/carlsverre/kosho@latest
```

## Why Use Kosho?

Kosho is well-suited for **running multiple concurrent AI coding agents** like Claude Code, each working on different features or branches in separate workspaces:

- **🤖 Multiple AI Agents**: Run Claude Code (or other AI tools) in separate worktrees simultaneously
- **⚡ Reduced Context Switching**: Work on multiple branches/features at once
- **🔒 Workspace Isolation**: Each agent operates in its own [git worktree] with independent working state
- **🚀 More Efficient Than Cloning**: Worktrees share the same `.git` directory - no duplicate repositories
- **📋 Easy Coordination**: Use `kosho list` to see what each agent is working on
- **🧹 Easy Cleanup**: Remove completed work environments without affecting your main repository

**Example workflow using Claude Code:**

```bash
# Start Claude Code on a feature branch
kosho open feat/widget claude

# Start another agent on bug fixes in a separate shell
kosho open bug/doodad claude

# Check what each agent is working on
kosho list

# Clean up clean worktrees
kosho prune
```

## Basic Usage

1. **Run a command in a worktree**:

   ```bash
   kosho run my-feature claude
   ```

2. **List all worktrees**:

   ```bash
   kosho list
   ```

3. **Prune clean worktrees**:

   ```bash
   kosho prune
   ```

## Commands

### `kosho run BRANCH [command...]`

Runs the provided command in a worktree checked out at the target `BRANCH`. If the worktree doesn't exist, it will be created.

**Arguments:**

- `BRANCH`: Name of the git branch
- `command...`: Any command you'd like to run in the worktree. I.e., `claude`

**Examples:**

```bash
# run claude in a bugfix worktree
kosho run bugfix claude

# open a shell in a worktree
kosho run playground zsh
```

### `kosho list`

Lists all kosho worktrees with their status and current git reference.

**Output:**

```
NAME        UPSTREAM  REF     STATUS
bugfix      main      bugfix  ahead 1
hotfix      main      bug/1   ahead 2 (dirty)
security    release   sec/1   ahead 1
```

### `kosho prune`

Cleanup clean worktrees and dangling worktree references. This will not delete git branches! If you'd like to clean up merged git branches, I recommend creating a script that looks something like this:

**git-janitor:**

```bash
#!/usr/bin/env bash
kosho prune
git fetch --prune
git branch --merged | grep -vE '^\*|main|\+' | xargs -r git branch -d
```

Then you can run this script (assuming it's on your `$PATH`) via `git janitor`.

## Hooks

Kosho supports hooks that run at specific points during worktree operations. Hooks are executable scripts stored in `.kosho/hooks/` and receive environment variables with context about the operation.

### Available Hooks

- **`create`**: Runs after a new worktree is created, before opening it
- **`run`**: Runs before running a command in a worktree

### Enabling Hooks

Kosho automatically creates sample hook files (`.sample` extension) in `.kosho/hooks/`. To enable a hook:

```bash
# Enable the create hook
mv .kosho/hooks/create.sample .kosho/hooks/create
```

### Environment Variables

Hooks receive these environment variables:

- `$KOSHO_HOOK`: The hook type (`create`, `run`)
- `$KOSHO_WORKTREE`: Name of the worktree being operated on
- `$KOSHO_REPO`: Path to the repository root
- `$PWD` / `$KOSHO_WORKTREE_PATH`: Full path to the worktree directory (the hook is also run within the worktree directory)
- `$KOSHO_CMD`: The name of the command that is being run in the worktree. Only present in the `run` hook

**Example create hook (`.kosho/hooks/create`):**

```bash
#!/bin/sh
echo "Setting up new worktree: $KOSHO_WORKTREE"
echo "Worktree path: $KOSHO_WORKTREE_PATH"
echo "Repository root: $KOSHO_REPO"

# Copy environment file from repo root to worktree
cp "$KOSHO_REPO/.env.example" "$KOSHO_WORKTREE_PATH/.env"

# Install dependencies when creating a new worktree
npm install
```

## How It Works

Kosho manages [git worktree]s in a `.kosho/` directory at your repository root:

```
your-repo/
├── .git/
├── .kosho/               # Kosho root directory
│   ├── .gitignore        # Kosho specific gitignore
│   ├── worktrees/
│   │   ├── feature-a/    # Worktree for feature-a
│   │   ├── bugfix/       # Worktree for bugfix
│   │   └── experiment/   # Worktree for experiment
│   └── hooks/
│       └── create        # Hook script which runs when creating a worktree
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
- Use `kosho prune` periodically to clean up any old worktrees

[git worktree]: https://git-scm.com/docs/git-worktree
