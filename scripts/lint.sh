#!/bin/bash

# Script to run golangci-lint with proper configuration

set -e

echo "Running golangci-lint..."

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null
then
    echo "golangci-lint could not be found, installing..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

# Run golangci-lint
if [ -f ".golangci.local.yml" ]; then
    echo "Using local configuration"
    golangci-lint run -c .golangci.local.yml ./...
else
    echo "Using default configuration"
    golangci-lint run ./...
fi

echo "Linting completed successfully!"