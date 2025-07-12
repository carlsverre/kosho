#!/bin/bash
set -euo pipefail

# Run local startup script if it exists
if [ -x /workspace/.kosho/startup.sh ]; then
    echo "Running local startup script..."
    source /workspace/.kosho/startup.sh
fi

# Ensure the firewall is setup
sudo init-firewall.sh

# Run zsh with any arguments passed to this script as the ubuntu user
exec /usr/bin/zsh "$@"
