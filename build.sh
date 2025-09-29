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

echo "Both MCP Filesystem Servers built successfully:"
echo "  - mcp-filesystem-server (raw implementation)"
echo "  - mcp-filesystem-server-mark3labs-mcp-go (SDK implementation)"
echo ""
echo "Usage:"
echo "  Both servers accept a -dir argument to specify the base directory"
echo ""
echo "Examples:"
echo "  # Use current directory as base"
echo "  ./mcp-filesystem-server"
echo "  ./mcp-filesystem-server-mark3labs-mcp-go"
echo ""
echo "  # Use specific directory as base"
echo "  ./mcp-filesystem-server -dir /path/to/directory"
echo "  ./mcp-filesystem-server-mark3labs-mcp-go -dir /path/to/directory"
echo ""
echo "  # Use relative directory"
echo "  ./mcp-filesystem-server -dir ./my-files"
echo "  ./mcp-filesystem-server-mark3labs-mcp-go -dir ./my-files"