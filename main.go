package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const BaseDir = "/Users/xu/mv"

type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      interface{}   `json:"id"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type InitializeParams struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ClientCapabilities `json:"capabilities"`
	ClientInfo      ClientInfo         `json:"clientInfo"`
}

type ClientCapabilities struct {
	Roots    *RootsCapability    `json:"roots,omitempty"`
	Sampling *SamplingCapability `json:"sampling,omitempty"`
}

type RootsCapability struct {
	ListChanged *bool `json:"listChanged,omitempty"`
}

type SamplingCapability struct{}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	Tools *ToolsCapability `json:"tools,omitempty"`
}

type ToolsCapability struct {
	ListChanged *bool `json:"listChanged,omitempty"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Required   []string               `json:"required"`
}

type ToolsListResult struct {
	Tools []Tool `json:"tools"`
}

type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type CallToolResult struct {
	Content []ToolContent `json:"content"`
	IsError *bool         `json:"isError,omitempty"`
}

type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func main() {
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		var request JSONRPCRequest
		if err := decoder.Decode(&request); err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error decoding request: %v", err)
			continue
		}

		response := handleRequest(request)
		if response != nil {
			if err := encoder.Encode(response); err != nil {
				log.Printf("Error encoding response: %v", err)
			}
		}
	}
}

func handleRequest(request JSONRPCRequest) *JSONRPCResponse {
	switch request.Method {
	case "initialize":
		return handleInitialize(request)
	case "tools/list":
		return handleToolsList(request)
	case "tools/call":
		return handleToolCall(request)
	case "notifications/initialized":
		return nil
	default:
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}
}

func handleInitialize(request JSONRPCRequest) *JSONRPCResponse {
	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: ServerCapabilities{
			Tools: &ToolsCapability{},
		},
		ServerInfo: ServerInfo{
			Name:    "filesystem-mcp-server",
			Version: "1.0.0",
		},
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

func handleToolsList(request JSONRPCRequest) *JSONRPCResponse {
	tools := []Tool{
		{
			Name:        "read_file",
			Description: "Read the complete contents of a file from the filesystem",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Path to the file to read",
					},
				},
				Required: []string{"path"},
			},
		},
		{
			Name:        "write_file",
			Description: "Write content to a file (overwrites existing content)",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Path to the file to write",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "Content to write to the file",
					},
				},
				Required: []string{"path", "content"},
			},
		},
		{
			Name:        "list_directory",
			Description: "List the contents of a directory",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Path to the directory to list",
					},
				},
				Required: []string{"path"},
			},
		},
		{
			Name:        "create_directory",
			Description: "Create a new directory",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Path to the directory to create",
					},
				},
				Required: []string{"path"},
			},
		},
		{
			Name:        "delete_file",
			Description: "Delete a file or directory",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Path to the file or directory to delete",
					},
				},
				Required: []string{"path"},
			},
		},
	}

	result := ToolsListResult{Tools: tools}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

func handleToolCall(request JSONRPCRequest) *JSONRPCResponse {
	var params CallToolParams
	paramBytes, _ := json.Marshal(request.Params)
	json.Unmarshal(paramBytes, &params)

	switch params.Name {
	case "read_file":
		return handleReadFile(request, params)
	case "write_file":
		return handleWriteFile(request, params)
	case "list_directory":
		return handleListDirectory(request, params)
	case "create_directory":
		return handleCreateDirectory(request, params)
	case "delete_file":
		return handleDeleteFile(request, params)
	default:
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Unknown tool",
			},
		}
	}
}

func validatePath(path string) (string, error) {
	cleanPath := filepath.Clean(path)

	if filepath.IsAbs(cleanPath) {
		if !strings.HasPrefix(cleanPath, BaseDir) {
			return "", fmt.Errorf("access denied: path outside allowed directory")
		}
		return cleanPath, nil
	}

	fullPath := filepath.Join(BaseDir, cleanPath)
	cleanFullPath := filepath.Clean(fullPath)

	if !strings.HasPrefix(cleanFullPath, BaseDir) {
		return "", fmt.Errorf("access denied: path outside allowed directory")
	}

	return cleanFullPath, nil
}

func handleReadFile(request JSONRPCRequest, params CallToolParams) *JSONRPCResponse {
	path, exists := params.Arguments["path"].(string)
	if !exists {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Missing required parameter: path",
			},
		}
	}

	validPath, err := validatePath(path)
	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32603,
				Message: err.Error(),
			},
		}
	}

	content, err := ioutil.ReadFile(validPath)
	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Error reading file: %v", err),
			},
		}
	}

	result := CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: string(content),
			},
		},
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

func handleWriteFile(request JSONRPCRequest, params CallToolParams) *JSONRPCResponse {
	path, pathExists := params.Arguments["path"].(string)
	content, contentExists := params.Arguments["content"].(string)

	if !pathExists || !contentExists {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Missing required parameters: path and content",
			},
		}
	}

	validPath, err := validatePath(path)
	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32603,
				Message: err.Error(),
			},
		}
	}

	err = os.MkdirAll(filepath.Dir(validPath), 0755)
	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Error creating directory: %v", err),
			},
		}
	}

	err = ioutil.WriteFile(validPath, []byte(content), 0644)
	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Error writing file: %v", err),
			},
		}
	}

	result := CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: fmt.Sprintf("Successfully wrote to file: %s", validPath),
			},
		},
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

func handleListDirectory(request JSONRPCRequest, params CallToolParams) *JSONRPCResponse {
	path, exists := params.Arguments["path"].(string)
	if !exists {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Missing required parameter: path",
			},
		}
	}

	validPath, err := validatePath(path)
	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32603,
				Message: err.Error(),
			},
		}
	}

	files, err := ioutil.ReadDir(validPath)
	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Error reading directory: %v", err),
			},
		}
	}

	var fileList []string
	for _, file := range files {
		if file.IsDir() {
			fileList = append(fileList, file.Name()+"/")
		} else {
			fileList = append(fileList, file.Name())
		}
	}

	result := CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: fmt.Sprintf("Directory contents:\n%s", strings.Join(fileList, "\n")),
			},
		},
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

func handleCreateDirectory(request JSONRPCRequest, params CallToolParams) *JSONRPCResponse {
	path, exists := params.Arguments["path"].(string)
	if !exists {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Missing required parameter: path",
			},
		}
	}

	validPath, err := validatePath(path)
	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32603,
				Message: err.Error(),
			},
		}
	}

	err = os.MkdirAll(validPath, 0755)
	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Error creating directory: %v", err),
			},
		}
	}

	result := CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: fmt.Sprintf("Successfully created directory: %s", validPath),
			},
		},
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

func handleDeleteFile(request JSONRPCRequest, params CallToolParams) *JSONRPCResponse {
	path, exists := params.Arguments["path"].(string)
	if !exists {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Missing required parameter: path",
			},
		}
	}

	validPath, err := validatePath(path)
	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32603,
				Message: err.Error(),
			},
		}
	}

	err = os.RemoveAll(validPath)
	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Error deleting file/directory: %v", err),
			},
		}
	}

	result := CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: fmt.Sprintf("Successfully deleted: %s", validPath),
			},
		},
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}
