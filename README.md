# MCP Filesystem Server

A Model Context Protocol (MCP) server that provides secure filesystem access to Claude and other MCP clients.

## Pure Golang Version

Prompt Claude Code

```
Create an MCP server in Golang to give Claude access to files in specific directory. The MCP server can expose file system access through standardized tools.
```

## github.com/mark3labs/mcp-go version

Prompt for `v2: sdk version`

```
1. describe this project
2. add one more command mcp-filesystem-server-mark3labs-mcp-go (the previous is mcp-filesystem-server), implement a new mcp server provides same abilities via sdk                             github.com/mark3labs/mcp-go
3. rename run.sh to build.sh, use build.sh to build two binaries.
```

Prompt for `as standard Go layout`: separate the code of two binaries using Golang's classic cmd dir structure.

Prompt for `works in any dir`: make mcp servers to works in any directories, not just BaseDir. the directory can be as boot argument.
