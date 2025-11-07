#!/bin/bash

# Script to install Git hooks

set -e

echo "Installing Git hooks..."

# Create symbolic links for hooks
HOOKS_DIR=".git/hooks"
SCRIPTS_DIR="scripts"

# Remove old hooks if they exist
rm -f "${HOOKS_DIR}/pre-commit-original"

# Pre-commit hook - copy the script directly to hooks directory
if [ -f "${SCRIPTS_DIR}/pre-commit" ]; then
    cp "${SCRIPTS_DIR}/pre-commit" "${HOOKS_DIR}/pre-commit"
    chmod +x "${HOOKS_DIR}/pre-commit"
    echo "Pre-commit hook installed"
else
    echo "Pre-commit hook script not found: ${SCRIPTS_DIR}/pre-commit"
fi

# Additional pre-commit hook for go modules - no longer needed as it's integrated
rm -f "${HOOKS_DIR}/pre-commit-go-mod" 2>/dev/null || true
rm -f "${HOOKS_DIR}/pre-commit-original" 2>/dev/null || true

echo "Git hooks installation completed!"