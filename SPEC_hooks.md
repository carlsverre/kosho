Specification: Worktree-Initialization Hooks for Kosho
====================================================

Overview
--------
Add first-class support for running arbitrary shell commands automatically after a new worktree is created.  
Configuration lives in `.kosho/settings.json` inside the repository and currently exposes a single key:

```json
{
  "worktree_init": ["pnpm i"]
}
```

When `kosho open …` creates a new worktree it shall execute every command in `worktree_init` sequentially inside the freshly-created worktree directory.

--------------------------------------------------------------------
1. Configuration file
--------------------------------------------------------------------

1.1 File location  
• Path: `<repo-root>/.kosho/settings.json`  
• Must live inside the repository so it is shareable and version-controllable.

1.2 JSON schema (v1)  

```
{
  "worktree_init": string[]   // optional
}
```

Notes  
• `worktree_init` is optional; omit or empty array ⇒ no hooks executed.  
• Each string is interpreted exactly like a user would type it in the shell, therefore it is split by the OS shell, not by Kosho (i.e. passed verbatim to `/bin/sh -c <cmd>` on Unix, `cmd.exe /C` on Windows).  
• Future-proofing: other top-level keys must be ignored with a warning.

1.3 Validation rules  
• File must contain valid UTF-8 JSON.  
• If `worktree_init` exists:  
 – It must be a JSON array.  
 – Every element must be a non-empty string.  
• Validation errors are surfaced to the user and abort the `open` command.

1.4 Default behaviour  
• If the file does not exist → act as if `{}`.  
• Corrupted JSON: abort with `kosho: invalid .kosho/settings.json – <reason>`.

--------------------------------------------------------------------
2. Implementation approach
--------------------------------------------------------------------

2.1 New package / file  
Create `internal/config.go`

```
type Settings struct {
    WorktreeInit []string `json:"worktree_init"`
}

func LoadSettings(repoRoot string) (Settings, error)
```

• Locate file, read, decode, validate (rules §1.3).  
• Return zero-value Settings on "file not found".

2.2 Command execution helper  
Add to `internal/worktree.go`:

```
func (kw *KoshoWorktree) RunInitHooks(cmds []string) error
```

Pseudo-code:

```
for _, c := range cmds {
    fmt.Printf("Running init hook: %q …\n", c)
    execCmd := exec.Command(shellCmd, shellArg, c) // /bin/sh -c on *nix
    execCmd.Dir = kw.WorktreePath()
    execCmd.Stdout, execCmd.Stderr, execCmd.Stdin = os.Stdout, os.Stderr, os.Stdin
    if err := execCmd.Run(); err != nil {
        return fmt.Errorf("init hook %q failed: %w", c, err)
    }
}
return nil
```

`shellCmd`/`shellArg` cross-platform helpers live in new `internal/shell.go`.

2.3 Wire-up in open workflow  

open.go:

```
settings, err := internal.LoadSettings(repoRoot)
…
err := createWorktree(name, kw, spec, settings)
```

Update `createWorktree` signature:

```
func createWorktree(name string, kw *internal.KoshoWorktree,
                    spec internal.BranchSpec, settings internal.Settings) error
```

Inside:

```
if err := kw.CreateIfNotExists(spec); err != nil { … }

if err := kw.RunInitHooks(settings.WorktreeInit); err != nil {
    // Remove partially-initialised worktree to keep repo clean
    _ = kw.Remove(true)
    return err
}
```

2.4 Backwards compatibility  
• Existing consumers without `.kosho/settings.json` remain unaffected.  
• Existing public API exported from `internal` remains unchanged except for the new optional helpers; CLI surface unchanged.

--------------------------------------------------------------------
3. Error handling policy
--------------------------------------------------------------------

• Any validation error ⇒ abort before mutating repo.  
• Any hook command returning non-zero exit status ⇒  
 – Print combined stdout/stderr.  
 – Delete the just-created worktree (best-effort).  
 – Exit `kosho open` with the same exit code.  
• Partial failures are *not* ignored; users who want "best-effort" can wrap their command: `"pnpm i || true"`.

--------------------------------------------------------------------
4. Testing strategy
--------------------------------------------------------------------

4.1 Unit tests  
• `config_test.go`  
 – valid file parsed correctly  
 – missing file returns defaults  
 – invalid JSON fails validation.  

• `worktree_hooks_test.go`  
 – `RunInitHooks` executes commands and stops on failure; use a temporary directory and stub commands (`echo`, `false`).  
 – Cross-platform shell path detection.

4.2 Integration tests (using `os/exec` and a temp git repo)  
• Scenario:  
 1. Create temp repo, commit, write `.kosho/settings.json` with `["touch initialized"]`.  
 2. Run `kosho open tempwt`.  
 3. Assert `initialized` file exists inside `.kosho/tempwt`.  

• Failure scenario: same but hook is `["false"]`; assert command exits non-zero and worktree directory no longer exists.

4.3 CI  
• Add above tests to `go test ./...` matrix on Linux/macOS/Windows.

--------------------------------------------------------------------
5. Documentation updates
--------------------------------------------------------------------

5.1 README  
• Add "Initialization hooks" subsection with purpose and example config.

```
echo '{
  "worktree_init": ["pnpm i"]
}' > .kosho/settings.json
```

• Explain execution order, failure policy, and quoting rules.

5.2 `kosho open --help`  
• Append note: "If .kosho/settings.json contains worktree_init hooks, they run automatically after the worktree is created."

5.3 CHANGELOG  
• `Added: worktree initialization hooks via .kosho/settings.json (#xyz).`

--------------------------------------------------------------------
6. Migration / versioning
--------------------------------------------------------------------

Version number bumps from `v0.x.y` → `v0.(y+1).0`. No breaking changes.

--------------------------------------------------------------------
Implementation is now fully specified and ready for development.
