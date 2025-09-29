# MCP Filesystem Server

A Model Context Protocol (MCP) server that provides secure filesystem access to Claude and other MCP clients.

## Pure Golang Version

Prompt Claude Code

```
Create an MCP server in Golang to give Claude access to files in specific directory. The MCP server can expose file system access through standardized tools.
```

## mark3labs/mcp-go version and refactor

Prompt for `v2: sdk version`

```
1. describe this project
2. add one more command mcp-filesystem-server-mark3labs-mcp-go (the previous is mcp-filesystem-server), implement a new mcp server provides same abilities via sdk                             github.com/mark3labs/mcp-go
3. rename run.sh to build.sh, use build.sh to build two binaries.
```

Prompt for `as standard Go layout`: separate the code of two binaries using Golang's classic cmd dir structure.

Prompt for `works in any dir`: make mcp servers to works in any directories, not just BaseDir. the directory can be as boot argument.

Prompt: update @test.sh to the latest version

## Cost

This learning mcp project costs $4.64.

```
Total cost:            $4.64 (costs may be inaccurate due to usage of unknown models)
Total duration (API):  24m 19.4s
Total duration (wall): 1h 0m 43.1s
Total code changes:    1310 lines added, 64 lines removed
Usage by model:
claude-3-5-haiku-20241022:  22.4k input, 1.3k output, 0 cache read, 0 cache write
cluade-sonnet-4-20250514:  313 input, 27.6k output, 4.5m cache read, 745.9k cache write
```
