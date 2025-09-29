package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const BaseDir = "/Users/xu/mv"

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

func main() {
	// Create a new MCP server using the mark3labs SDK
	s := server.NewMCPServer(
		"filesystem-mcp-server-mark3labs",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	// Add read_file tool
	readFileTool := mcp.NewTool("read_file",
		mcp.WithDescription("Read the complete contents of a file from the filesystem"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to the file to read"),
		),
	)

	s.AddTool(readFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := request.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		validPath, err := validatePath(path)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		content, err := ioutil.ReadFile(validPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error reading file: %v", err)), nil
		}

		return mcp.NewToolResultText(string(content)), nil
	})

	// Add write_file tool
	writeFileTool := mcp.NewTool("write_file",
		mcp.WithDescription("Write content to a file (overwrites existing content)"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to the file to write"),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("Content to write to the file"),
		),
	)

	s.AddTool(writeFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := request.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		content, err := request.RequireString("content")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		validPath, err := validatePath(path)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		err = os.MkdirAll(filepath.Dir(validPath), 0755)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error creating directory: %v", err)), nil
		}

		err = ioutil.WriteFile(validPath, []byte(content), 0644)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error writing file: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully wrote to file: %s", validPath)), nil
	})

	// Add list_directory tool
	listDirectoryTool := mcp.NewTool("list_directory",
		mcp.WithDescription("List the contents of a directory"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to the directory to list"),
		),
	)

	s.AddTool(listDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := request.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		validPath, err := validatePath(path)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		files, err := ioutil.ReadDir(validPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error reading directory: %v", err)), nil
		}

		var fileList []string
		for _, file := range files {
			if file.IsDir() {
				fileList = append(fileList, file.Name()+"/")
			} else {
				fileList = append(fileList, file.Name())
			}
		}

		return mcp.NewToolResultText(fmt.Sprintf("Directory contents:\n%s", strings.Join(fileList, "\n"))), nil
	})

	// Add create_directory tool
	createDirectoryTool := mcp.NewTool("create_directory",
		mcp.WithDescription("Create a new directory"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to the directory to create"),
		),
	)

	s.AddTool(createDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := request.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		validPath, err := validatePath(path)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		err = os.MkdirAll(validPath, 0755)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error creating directory: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully created directory: %s", validPath)), nil
	})

	// Add delete_file tool
	deleteFileTool := mcp.NewTool("delete_file",
		mcp.WithDescription("Delete a file or directory"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to the file or directory to delete"),
		),
	)

	s.AddTool(deleteFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := request.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		validPath, err := validatePath(path)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		err = os.RemoveAll(validPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error deleting file/directory: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully deleted: %s", validPath)), nil
	})

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}