#!/usr/bin/env bash

# MCP Filesystem Server Launcher
# This script builds and runs the MCP filesystem server

set -e

# Build the server
echo "Building MCP filesystem server..."
go build -o mcp-filesystem-server main.go

# Check if target directory exists
TARGET_DIR="/Users/xu/Documents/mv"
if [ ! -d "$TARGET_DIR" ]; then
    echo "Warning: Target directory $TARGET_DIR does not exist."
    echo "Creating directory..."
    mkdir -p "$TARGET_DIR"
fi

echo "MCP Filesystem Server ready."
echo "Target directory: $TARGET_DIR"
echo "Starting server..."

# Run the server
exec ./mcp-filesystem-server