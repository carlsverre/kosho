# Kosho CLI Project Plan

## Overview

Kosho is a CLI tool that creates git worktrees in pre-configured locations and launches interactive Docker development environments. The tool aims to streamline the workflow of creating isolated development environments with proper volume mounts and container configuration.

## Dependencies to use

- write code in golang
- cobra for argument parsing
- github.com/docker/docker for managing docker containers
- https://github.com/go-git/go-git for git

## Core Requirements

### 1. Worktree Management

- **Location Strategy**: All worktrees for repo "foo" should be stored in `$XDG_DATA_HOME/foo/...` or `$HOME/.local/share` if XDG_DATA_HOME is not specified
- **Dynamic Naming**: Generate Docker-style names (adjective_noun format) for worktrees
- **Name Override**: Allow users to specify custom worktree names

### 2. Docker Integration

- **Run containers with an equivalent command**:
  ```
  docker run -it \
    -v$repo-worktree-config:/home/ubuntu/.claude \
    -v$repo-worktree-history:/commandhistory \
    -v$repo-worktree-path:/workspace \
    --cap-add=NET_ADMIN \
    --cap-add=NET_RAW \
    kosho-img
  ```
- **Volume Management**: Create and manage the three required volumes
  - $repo-worktree-config and $repo-worktree-history should be named docker volumes
  - $repo-worktree-path should be the path to the git worktree
- **Image Management**: Assume the kosho-img is already built

## Implementation Plan

- [ ] Set up Cobra CLI framework with root command
- [ ] Add worktree command with name generation
- [ ] Implement XDG data directory resolution
- [ ] Create git worktree using go-git
- [ ] Add Docker container management
- [ ] Create named volumes for config and history
- [ ] Mount worktree path as workspace volume
- [ ] Add container lifecycle management (start/stop)
- [ ] Add cleanup commands for worktrees and volumes
- [ ] Handle error cases and validation
