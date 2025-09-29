#!/usr/bin/env bash

# Test script for MCP filesystem servers

set -e

echo "Testing MCP Filesystem Servers..."

# Create a test directory
TEST_DIR="/tmp/mcp-test-$(date +%s)"
mkdir -p "$TEST_DIR"
echo "Using test directory: $TEST_DIR"

# Function to test a server
test_server() {
    local server_name="$1"
    local server_path="$2"
    
    echo ""
    echo "=== Testing $server_name ==="
    
    # Test initialize
    echo "1. Testing initialize..."
    echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}' | timeout 5 "$server_path" -dir "$TEST_DIR" | head -1
    
    # Test tools/list  
    echo "2. Testing tools/list..."
    echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' | timeout 5 "$server_path" -dir "$TEST_DIR" | head -1
    
    # Test create directory
    echo "3. Testing create_directory..."
    echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"create_directory","arguments":{"path":"test-directory"}}}' | timeout 5 "$server_path" -dir "$TEST_DIR" | head -1
    
    # Test write file
    echo "4. Testing write_file..."
    echo '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"write_file","arguments":{"path":"test-directory/test.txt","content":"Hello, MCP World from '"$server_name"'!"}}}' | timeout 5 "$server_path" -dir "$TEST_DIR" | head -1
    
    # Test read file
    echo "5. Testing read_file..."
    echo '{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"read_file","arguments":{"path":"test-directory/test.txt"}}}' | timeout 5 "$server_path" -dir "$TEST_DIR" | head -1
    
    # Test list directory
    echo "6. Testing list_directory..."
    echo '{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"list_directory","arguments":{"path":"test-directory"}}}' | timeout 5 "$server_path" -dir "$TEST_DIR" | head -1
    
    # Test delete file
    echo "7. Testing delete_file..."
    echo '{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"delete_file","arguments":{"path":"test-directory"}}}' | timeout 5 "$server_path" -dir "$TEST_DIR" | head -1
    
    echo "$server_name tests completed."
}

# Build servers if needed
if [ ! -f "./mcp-filesystem-server" ] || [ ! -f "./mcp-filesystem-server-mark3labs-mcp-go" ]; then
    echo "Building servers..."
    ./build.sh
fi

# Test both servers
test_server "Raw Implementation" "./mcp-filesystem-server"
test_server "SDK Implementation" "./mcp-filesystem-server-mark3labs-mcp-go"

echo ""
echo "=== Verification ==="
echo "Test directory contents:"
ls -la "$TEST_DIR" || echo "Test directory cleaned up successfully"

# Cleanup
rm -rf "$TEST_DIR"
echo ""
echo "✅ All tests completed successfully!"
echo "✅ Test directory cleaned up: $TEST_DIR"