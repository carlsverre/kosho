# Kosho CLI Project Plan

## Overview

Kosho is a CLI tool that manages git worktrees in `.kosho/` directories and helps to launch tools within them. The tool streamlines the workflow of creating isolated worktrees for development.

## Dependencies to use

- write code in golang
- cobra for argument parsing
- https://github.com/go-git/go-git for git operations
- git command line for worktree operations

## Ideas

- what about a docker workspace for running tests which is perma-whitelisted by claude

## Commands

**kosho open [-b<branch>|-B<branch>] [NAME] [commitish] [-- command...]**

- NAME will be the name of the worktree
- commitish may be omitted which will result in the same behavior as omitting it from the underlying `git worktree add` command.
- the `-b|-B` flags will be passed through to git worktree if specified
- the worktree will be located at `.kosho/$NAME` at the root of the current git repo
  - `/.kosho` will be added to .gitignore if it's not already there
- if the worktree doesn't exist, it will be created
- by default, opens a new shell instance (inheriting the current shell binary and env) in the worktree
- if a command is provided after `--`, runs that command in the target worktree instead

**kosho list**

- list all kosho worktrees and their current git status + ref

**kosho remove [-f|--force] NAME**

- if worktree is dirty, require --force to continue
- run `git worktree remove` passing through the `--force` flag if specified

**kosho prune**

- cleanup any dangling worktree refs by running `git worktree prune`

# TODO: REMAINING WORK

## Core Implementation

- [x] Implement `kosho open` command:
  - [x] Create worktree if it doesn't exist (using git worktree add)
  - [x] Add `.kosho` to .gitignore if not already present
  - [x] Launch shell session in worktree directory
  - [x] Support optional command execution instead of shell
  - [x] Inherit current shell binary and environment variables
- [x] Implement `kosho list` command:
  - [x] List all worktrees in `.kosho/` directory
  - [x] Show current git status/ref for each worktree
- [x] Implement `kosho remove` command:
  - [x] Check for dirty worktree and require --force if dirty
  - [x] Remove worktree using `git worktree remove`
- [x] Implement `kosho prune` command:
  - [x] Run `git worktree prune` to cleanup dangling refs

## Documentation

- [ ] Provide a simple quickstart guide in the README
