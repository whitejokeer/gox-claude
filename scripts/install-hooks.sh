#!/bin/bash
set -e

echo "Installing pre-commit hooks..."

# Check if pre-commit is installed
if ! command -v pre-commit &> /dev/null; then
    echo "pre-commit not found. Installing..."
    if command -v pip &> /dev/null; then
        pip install pre-commit
    elif command -v brew &> /dev/null; then
        brew install pre-commit
    else
        echo "Error: Neither pip nor brew found. Please install pre-commit manually."
        echo "Visit: https://pre-commit.com/#install"
        exit 1
    fi
fi

# Install the git hook scripts
pre-commit install
pre-commit install --hook-type commit-msg

echo "Pre-commit hooks installed successfully!"
echo "Running pre-commit on all files..."
pre-commit run --all-files || true

echo "Done! Pre-commit hooks are now active."