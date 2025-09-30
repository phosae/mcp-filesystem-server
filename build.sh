#!/usr/bin/env bash

# MCP Filesystem Server Builder
# This script builds both MCP filesystem server implementations

set -e

echo "Building MCP filesystem servers..."

# Build the raw implementation
echo "Building raw implementation (mcp-filesystem-server)..."
go build -o mcp-filesystem-server ./cmd/mcp-filesystem-server

# Build the SDK implementation
echo "Building SDK implementation (mcp-filesystem-server-mark3labs-mcp-go)..."
go build -o mcp-filesystem-server-mark3labs-mcp-go ./cmd/mcp-filesystem-server-mark3labs-mcp-go

# Check if target directory exists
TARGET_DIR="/Users/xu/Documents/mv"
if [ ! -d "$TARGET_DIR" ]; then
    echo "Warning: Target directory $TARGET_DIR does not exist."
    echo "Creating directory..."
    mkdir -p "$TARGET_DIR"
fi

echo "Both MCP Filesystem Servers built successfully:"
echo "  - mcp-filesystem-server (raw implementation)"
echo "  - mcp-filesystem-server-mark3labs-mcp-go (SDK implementation)"
echo "Target directory: $TARGET_DIR"
echo ""
echo "To run servers:"
echo "  Raw: ./mcp-filesystem-server"
echo "  SDK: ./mcp-filesystem-server-mark3labs-mcp-go"