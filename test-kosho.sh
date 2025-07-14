#!/bin/bash

# Test script for kosho - creates a temporary git repo with sample history
# and spawns a Claude Code sub-agent to test kosho functionality
#
# Usage: ./test-kosho.sh [additional test instructions...]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building kosho...${NC}"
go build -o kosho

echo -e "${GREEN}Creating temporary test directory...${NC}"
TEMP_DIR=$(mktemp -d)
echo "Test directory: $TEMP_DIR"

# Cleanup function
cleanup() {
    echo -e "${YELLOW}Cleaning up temporary directory: $TEMP_DIR${NC}"
    rm -rf "$TEMP_DIR"
    echo -e "${GREEN}Cleanup complete${NC}"
}

# Set trap to cleanup on exit
trap cleanup EXIT

# Store the original directory (where kosho binary is)
KOSHO_DIR="$PWD"

# Change to temp directory
cd "$TEMP_DIR"

echo -e "${GREEN}Initializing git repository...${NC}"
git init
git config user.name "Test User"
git config user.email "test@example.com"

echo -e "${GREEN}Creating sample git history...${NC}"

# Initial commit
echo "# Test Project" > README.md
echo "print('Hello, World!')" > main.py
git add .
git commit -m "Initial commit: Add README and main.py"

# Second commit - update README
echo -e "\nThis is a test project for kosho." >> README.md
git add README.md
git commit -m "Update README with description"

# Third commit - add new file
echo "def greet(name):" > utils.py
echo "    return f'Hello, {name}!'" >> utils.py
git add utils.py
git commit -m "Add utils.py with greet function"

# Fourth commit - update main.py
echo "from utils import greet" > main.py
echo "print(greet('World'))" >> main.py
git add main.py
git commit -m "Update main.py to use utils.greet"

echo -e "${GREEN}Sample git history created:${NC}"
git log --oneline

echo -e "${YELLOW}Starting Claude Code sub-agent in test directory...${NC}"
echo -e ""

cp "$KOSHO_DIR/kosho" "$TEMP_DIR/"

cd $TEMP_DIR

CLAUDE_ARGS=(
  --print
  --output-format=stream-json
  --verbose
  --allowedTools='Bash(./kosho:*) Bash(git:*) Bash(echo:*) Bash(cat:*) Edit MultiEdit Write'
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

The kosho CLI tool is available at './kosho'.
Kosho manages git worktrees in .kosho/ directories for isolated development environments.

Key kosho commands:
- ./kosho list: Show all worktrees and their status
- ./kosho open <name> [-b branch]: Create/open a worktree
- ./kosho merge <name>: Merge worktree branch into current branch
- ./kosho remove <name>: Remove a worktree
- ./kosho prune: Clean up dangling worktree references

Always report on what works and any issues you find."

# Default test instructions if none provided
DEFAULT_INSTRUCTIONS="Test the basic kosho workflow:

1. Run './kosho list' to see current worktrees (should be none)
2. Create a feature worktree: './kosho open feature-test -b feature/test'
3. List worktrees again to see the new one
4. Go into the worktree (.kosho/feature-test) and make some changes, commit them
5. Come back to repo root and run './kosho merge feature-test'
6. Test other commands like './kosho remove feature-test'"

# Use provided instructions or default
if [ $# -gt 0 ]; then
    TEST_INSTRUCTIONS="$*"
else
    TEST_INSTRUCTIONS="$DEFAULT_INSTRUCTIONS"
fi

FULL_PROMPT="$SYSTEM_PROMPT

$TEST_INSTRUCTIONS"

# Start Claude Code with instructions, parse output with jq
claude "${CLAUDE_ARGS[@]}" "$FULL_PROMPT" | jq -r "$JQ_FILTER"
