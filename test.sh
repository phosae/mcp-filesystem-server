#!/usr/bin/env bash

# Test script for MCP filesystem server

echo "Testing MCP Filesystem Server..."

# Start the server in the background
./mcp-filesystem-server &
SERVER_PID=$!

# Give server time to start
sleep 1

# Test initialize
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}' | ./mcp-filesystem-server &

# Test tools/list
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' | ./mcp-filesystem-server &

# Test create directory
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"create_directory","arguments":{"path":"test-directory"}}}' | ./mcp-filesystem-server &

# Test write file
echo '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"write_file","arguments":{"path":"test-directory/test.txt","content":"Hello, MCP World!"}}}' | ./mcp-filesystem-server &

# Test read file
echo '{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"read_file","arguments":{"path":"test-directory/test.txt"}}}' | ./mcp-filesystem-server &

# Test list directory
echo '{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"list_directory","arguments":{"path":"test-directory"}}}' | ./mcp-filesystem-server &

wait

echo "Test completed. Check /Users/xu/Documents/mv for created files."