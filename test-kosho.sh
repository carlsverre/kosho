#!/bin/bash

# Test script for kosho - creates a temporary git repo with sample history
# and opens an interactive shell with kosho in PATH for testing

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
git add README.md main.py
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

echo -e "${YELLOW}Starting interactive shell with kosho in PATH...${NC}"
echo -e "${YELLOW}Try commands like:${NC}"
echo -e "  ${GREEN}kosho open feature-test${NC}"
echo -e "  ${GREEN}kosho list${NC}"
echo -e "  ${GREEN}kosho remove feature-test${NC}"
echo -e ""
echo -e "${YELLOW}Type 'exit' to close the shell and cleanup${NC}"
echo -e ""

# Add kosho directory to PATH and start shell
export PATH="$KOSHO_DIR:$PATH"

$SHELL
