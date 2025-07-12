# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Kosho is a CLI tool that manages git worktrees in `.kosho/` directories and launches Docker container development environments attached to those worktrees.

## Build & Test Commands

- Build: `go build`
- Run: `kosho ...`
- Lint: `golangci-lint run`
- Format: `golangci-lint fmt`
- Check for dead code: `deadcode .`
- Tidy dependencies: `go mod tidy`

## Code Style Guidelines

- **Imports**: Standard Go import organization (stdlib, external, internal)
- **Error Handling**: Return errors explicitly
- **Naming**: Use Go conventions (CamelCase for exported, camelCase for unexported)
- **Types**: Use strong typing; prefer interfaces for dependencies
- **Documentation**: Document all exported functions and types

## PLAN.md

- Make sure to mark completed tasks in PLAN.md before completing your work.
