# Kosho Codebase Conventions

This document outlines the coding conventions and patterns used throughout the kosho project to maintain consistency and readability.

## Go Coding Conventions

### Code Formatting and Style
- Uses standard `gofmt` formatting throughout
- Import grouping follows Go standards:
  - Standard library imports first
  - Blank line separator
  - Third-party imports
  - Local imports with proper ordering
- Opening braces on same line as function/struct declarations
- Uses both short `:=` and explicit `var` declarations appropriately

### Error Handling Patterns
- Consistent error wrapping with `fmt.Errorf` and `%w` verb
- Error prefixing with context: `"kosho: invalid .kosho/settings.json –"`
- Early returns on errors rather than deep nesting
- Main function uses simple error handling with stderr output and `os.Exit(1)`

### Function and Method Organization
- Pointer receivers for structs: `(kw *KoshoWorktree)`
- Constructor pattern: `NewKoshoWorktree` functions
- Clear, descriptive method names: `WorktreePath()`, `CreateIfNotExists()`, `IsDirty()`
- Related functions grouped together (git operations, worktree operations)

### Package Structure
- `main` package only in `main.go`
- `cmd` package for CLI commands
- `internal` package for core functionality
- Uses `github.com/carlsverre/kosho` module prefix

### Documentation Style
- Clear comments above struct declarations
- Function comments explaining purpose and behavior
- Struct fields include inline comments
- Package-level documentation with package name and description

### Variable Declaration Patterns
- Short `:=` declarations for new variables within functions
- Leverages Go's zero values appropriately
- camelCase for unexported variables, PascalCase for exported
- `var` declarations for package-level variables

## Project Structure

### Directory Organization
- **Standard Go Layout**:
  - `cmd/` - Command-line interface code
  - `internal/` - Private application and library code
  - `main.go` - Application entry point
- **Configuration Directories**:
  - `.kosho/` - Project-specific configuration
  - `.github/` - GitHub workflows and settings

### File Naming Patterns
- **Go Source Files**: Descriptive action names (`list.go`, `merge.go`, `open.go`)
- **Internal Files**: Functional names (`config.go`, `git.go`, `worktree.go`)
- **Test Files**: Standard `*_test.go` suffix
- **Documentation**: Uppercase for important docs (`README.md`, `LICENSE`)

## Naming Conventions

### Functions and Methods
- **Exported**: PascalCase (`LoadSettings`, `FindGitRoot`, `NewKoshoWorktree`)
- **Unexported**: camelCase (`checkMergeArgs`, `createWorktree`, `validateWorktreeInit`)
- **Test Functions**: `Test` + PascalCase with descriptive scenarios

### Variables
- **Local Variables**: camelCase (`settingsPath`, `worktreePath`, `currentBranch`)
- **Package-level**: camelCase (`keepFlag`, `commitMessage`, `rootCmd`)
- **Short Names**: Common patterns use abbreviated names (`cmd`, `err`, `kw`)

### Types and Structs
- **Exported Types**: PascalCase (`Settings`, `BranchSpec`, `KoshoWorktree`)
- **Struct Fields**: PascalCase with appropriate JSON tags
- **JSON Tags**: snake_case (`json:"worktree_init"`)

### Packages
- **Package Names**: lowercase, single word (`main`, `cmd`, `internal`)
- **Import Paths**: Follow module structure (`github.com/carlsverre/kosho/cmd`)

## Testing Conventions

### Test Organization
- Tests co-located with source code in same package
- Standard `*_test.go` naming convention
- No separate test directories

### Test Function Naming
- Pattern: `TestFunctionName_Scenario`
- Examples: `TestLoadSettings_ValidJSON`, `TestRunInitHooks_EmptyCommands`
- Table-driven tests for comprehensive scenario coverage

### Testing Framework
- Uses standard Go `testing` package only
- Zero external testing dependencies
- Pure Go standard library approach

### Test Data and Fixtures
- Temporary directories for filesystem tests: `os.MkdirTemp("", "kosho-config-test")`
- Helper functions: `setupTempWorktree(t *testing.T)`
- Inline JSON strings for configuration testing
- Dynamic test fixtures rather than static files

### Integration Testing
- Manual integration testing via `test-kosho.sh`
- Docker containerization for safe testing environment
- Ubuntu 24.04 container with git, node, and development tools
- Automated test scenarios with environment isolation

## Configuration Patterns

### Configuration Files
- `.kosho/settings.json` - Project-specific settings
- JSON format for configuration data
- snake_case naming for JSON keys (`worktree_init`)

### Settings Structure
- Clear struct definitions with JSON tags
- Validation functions for configuration integrity
- Error handling for missing or invalid configuration

## Build and Development

### Commands
- **Build**: `go build`
- **Lint**: `golangci-lint run`
- **Format**: `golangci-lint fmt`
- **Manual Testing**: `./test-kosho.sh`

### Safety Practices
- Containerized testing to prevent repository damage
- Comprehensive error testing with edge cases
- Platform-specific considerations (Windows compatibility)
- Clean setup/teardown patterns in tests

## Additional Patterns

### Struct Composition
- Uses composition with embedded types
- Clear field purposes and relationships
- Proper JSON marshaling/unmarshaling

### File Operations
- Consistent error handling for file operations
- Proper cleanup using `defer`
- Leverages Go's standard library for file manipulation

### Command Execution
- Uses `exec.Command` for external processes
- Proper error handling and output capture
- Consistent patterns for git operations

This document reflects the mature Go development practices demonstrated throughout the kosho codebase, emphasizing clarity, consistency, and maintainability.
