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

echo -e "${YELLOW}Running Claude Code in test container...${NC}"
echo -e ""

# Default test instructions if none provided
DEFAULT_INSTRUCTIONS="Test the basic kosho workflow:

1. Run 'kosho list' to see current worktrees (should be none)
2. Create a feature worktree: 'kosho open feature-test -b feature/test'
3. List worktrees again to see the new one
4. Go into the worktree (.kosho/feature-test) and make some changes, commit them
5. Come back to repo root and run 'kosho merge feature-test'
6. Test other commands like 'kosho remove feature-test'"

# Use provided instructions or default
if [ $# -gt 0 ]; then
    TEST_INSTRUCTIONS="$*"
else
    TEST_INSTRUCTIONS="$DEFAULT_INSTRUCTIONS"
fi

# Ensure the image exists
docker build -t kosho-tester kosho-tester

# Start Claude Code with instructions, parse output with jq
docker run --rm -it \
  -v "$(pwd)/kosho:/usr/local/bin/kosho:ro" \
  -v "$(pwd):/kosho:ro" \
  -v "$HOME/.claude/.credentials.json:/home/ubuntu/.claude/.credentials.json:ro" \
  --cap-add=NET_ADMIN \
  --cap-add=NET_RAW \
  --env TEST_INSTRUCTIONS="${TEST_INSTRUCTIONS}" \
  kosho-tester claude.sh
