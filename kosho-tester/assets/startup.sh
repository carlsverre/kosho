#!/bin/bash
set -euo pipefail
sudo init-firewall.sh
exec "$@"
