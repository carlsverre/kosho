#!/bin/bash
set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

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
