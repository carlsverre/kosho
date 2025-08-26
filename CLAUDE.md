# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Kosho is a CLI tool that manages git worktrees in `.kosho/` directories for isolated development environments, particularly useful for running multiple concurrent AI coding agents.

## Current Development

Core features complete including worktree management and merge workflows.

## Build & Test Commands

- Build: `go build`
- Run: `./kosho ...`
- Lint: `golangci-lint run`
- Format: `golangci-lint fmt`
- Check for dead code: `deadcode .`
- Tidy dependencies: `go mod tidy`
- Manual testing: `./test-kosho.sh`

## Code Style Guidelines

- **Imports**: Standard Go import organization (stdlib, external, internal)
- **Error Handling**: Return errors explicitly
- **Naming**: Use Go conventions (CamelCase for exported, camelCase for unexported)
- **Types**: Use strong typing; prefer interfaces for dependencies
- **Documentation**: Document all exported functions and types

## Testing Guidelines

- **Manual Testing**: Use `./test-kosho.sh` to test kosho functionality in a safe environment
- **Important**: Never manually test kosho commands directly in the kosho repository itself
- The test script creates an isolated environment to safely verify kosho operations

## Pre-Commit Guidelines

- Always run `golangci-lint fmt/run` before committing to ensure code is properly formatted and linted
