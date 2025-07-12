# Kosho CLI Project Plan

## Overview

Kosho is a CLI tool that creates git worktrees in a `.kosho` folder at the repo root and launches interactive development environments. The tool aims to streamline the workflow of creating isolated development environments.

## Dependencies to use

- write code in golang
- cobra for argument parsing
- github.com/docker/docker/pkg/namesgenerator for naming
- https://github.com/go-git/go-git for git operations
- git command line for worktree operations

## Commands

**kosho new [-b<branch>|-B<branch>] [NAME] [commitish]**

- NAME will be the name of the worktree
- commitish may be omitted which will result in the same behavior as omitting it from the underlying `git worktree add` command.
- the `-b|-B` flags will be passed through to git worktree if specified
- the worktree will be located at `.kosho/$NAME` at the root of the current git repo
  - `/.kosho` will be added to .gitignore if it's not already there
- this command will fail if the worktree already exists
- after creating the worktree this command will fall through to `kosho start`

**kosho start [-d] [NAME]**

- start the Kosho docker container in the worktree
- if -d is specified, run the container in the background using some kind of sleep-forever command
- if -d is not specified, run the container interactively by fully passing through stdin/out/err and all signals and so on.

**kosho list**

- list all kosho worktrees, their current git status/ref, and their running state (along with the container name if running)

**kosho stop**

- stop a kosho container if running

**kosho remove [-f|--force] NAME**

- if worktree is dirty, require --force to continue
- stop the container if running
- remove the container
- run `git worktree remove` passing through the `--force` flag if specified

# TODO: REMAINING WORK

## Docker Container Integration

- [ ] Replace bash sessions with actual Docker container management
- [ ] Implement container naming convention (e.g., `kosho-{repo-name}-{worktree-name}`)
- [ ] Create and manage Docker volumes for persistent data:
  - [ ] Config volume: `{repo-name}-{worktree-name}-config` mounted to `/home/ubuntu/.claude`
  - [ ] History volume: `{repo-name}-{worktree-name}-history` mounted to `/commandhistory`
  - [ ] Workspace volume: bind mount worktree path to `/workspace`
- [ ] Implement Docker container lifecycle:
  - [ ] `kosho start` interactive mode: full stdin/stdout/stderr passthrough with signal handling
  - [ ] `kosho start -d` detached mode: run container with sleep-forever command
  - [ ] `kosho stop` functionality: stop running containers
  - [ ] Container and volume cleanup in `kosho remove`
- [ ] Enhanced `kosho list` command:
  - [ ] Show actual git status/ref for each worktree
  - [ ] Show container running status and container name
  - [ ] Show last activity/created dates

## Container Configuration

- [ ] Add container configuration options:
  - [ ] `--cap-add=NET_ADMIN --cap-add=NET_RAW` capabilities
  - [ ] Environment variable passthrough
  - [ ] Port mapping options
  - [ ] Custom image selection
  - [ ] Pass through the `TZ` environment variable into the container
- [ ] Container image management:
  - [ ] If `.kosho/Dockerfile` exists then build and use that docker image
        Otherwise use the image `kosho-runtime` which is expected to exist.

## Enhanced Features

- [ ] Add `kosho attach NAME` - attach an interactive zsh shell to a running kosho container
- [ ] Add `kosho prune` - cleanup any dangling volumes and run `git worktree prune` to cleanup any dangling worktree refs

## Documentation

- [ ] Provide a simple quickstart guide in the README
