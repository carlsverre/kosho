#!/bin/bash
set -euo pipefail

CLAUDE_ARGS=(
  --print
  --output-format=stream-json
  --verbose
  --dangerously-skip-permissions
)

# jq filter for pretty-printing Claude Code streaming JSON
JQ_FILTER='
if .type == "system" then
  "\u001b[35m[SYSTEM]\u001b[0m \(.subtype // "unknown") | Session: \(.session_id[0:8])... | CWD: \(.cwd)"
elif .type == "user" then
  if .message.content[0].type == "tool_result" then
    "\u001b[32m[TOOL_RESULT]\u001b[0m \(.message.content[0].tool_use_id[0:12])...\n\(.message.content[0].content | if type == "string" and (. | length) > 200 then .[0:200] + "..." else . end)"
  else
    "\u001b[32m[USER]\u001b[0m \(.message.content[0].text // "No text content")"
  end
elif .type == "assistant" then
  if .message.content[0].type == "text" then
    "\u001b[34m[ASSISTANT]\u001b[0m \(.message.content[0].text)"
  elif .message.content[0].type == "tool_use" then
    "\u001b[36m[TOOL_USE]\u001b[0m \(.message.content[0].name) | ID: \(.message.content[0].id[0:12])...\n  Input: \(.message.content[0].input | tostring | if length > 100 then .[0:100] + "..." else . end)"
  else
    "\u001b[34m[ASSISTANT]\u001b[0m Unknown content type: \(.message.content[0].type // "none")"
  end
elif .type == "result" then
  "\u001b[33m[RESULT]\u001b[0m \(.subtype) | Duration: \(.duration_ms)ms | Turns: \(.num_turns) | Cost: $\(.total_cost_usd)\n\(.result | if type == "string" and (. | length) > 300 then .[0:300] + "..." else . end)"
else
  "\u001b[31m[UNKNOWN]\u001b[0m Type: \(.type) | \(.)"
end'

# Build the test instruction
SYSTEM_PROMPT="You are now in a test repository for kosho. This repo has sample git history with main.py, utils.py, and README.md files.

The kosho CLI tool is available on the path as 'kosho'.
The kosho codebase and docs is available at '/kosho'.

Kosho manages git worktrees in .kosho/ directories for isolated development environments.

Key kosho commands:
- kosho list: Show all worktrees and their status
- kosho open <name> [-b branch]: Create/open a worktree
- kosho merge <name>: Merge worktree branch into current branch
- kosho remove <name>: Remove a worktree
- kosho prune: Clean up dangling worktree references

Always report on what works and any issues you find."

FULL_PROMPT="${SYSTEM_PROMPT}

${TEST_INSTRUCTIONS}"

# Run Claude Code with the constructed prompt
exec claude "${CLAUDE_ARGS[@]}" "$FULL_PROMPT" | jq -r "$JQ_FILTER"
