#!/bin/sh
# Shell-AI installation script

echo "Shell-AI Installation Helper"
echo "============================"
echo ""
echo "Note: The recommended way to install Shell-AI is using Go's built-in install command:"
echo "  go install github.com/jwswj/shell-ai/cmd/shai@latest"
echo ""
echo "This script will help you set up the configuration directory and build the binary locally."
echo ""

# Create config directory
CONFIG_DIR="$HOME/.config/shell-ai"
echo "Creating config directory at $CONFIG_DIR..."
mkdir -p "$CONFIG_DIR"

# Build the binary
echo "Building Shell-AI locally..."
make build

echo ""
echo "Shell-AI has been built locally at $(pwd)/bin/shai"
echo ""
echo "You can run it directly with:"
echo "  $(pwd)/bin/shai find all files modified in the last 24 hours"
echo ""
echo "Or you can install it to your GOPATH with:"
echo "  go install ./cmd/shai"
echo ""
echo "To configure Shell-AI, create a config.json file at $CONFIG_DIR/config.json"
echo "See README.md for configuration options." 