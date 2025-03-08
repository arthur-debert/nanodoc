#!/bin/bash
# Wrapper script for package-managers/debian/update-apt-package.sh

# Get the directory of the current script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Execute the actual script with all arguments passed through
"$PROJECT_ROOT/package-managers/debian/update-apt-package.sh" "$@"
