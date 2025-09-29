# MCP Filesystem Server

Simple Golang server for learning Model Context Protocol (MCP).

This server provides secure filesystem access to Claude and other MCP clients.

Integrate it with Claude Code

```bash
{
./build.sh
DIR=$(pwd)
claude mcp add-json access-fs '{"type":"stdio","command":"'$DIR'/mcp-filesystem-server","args":["-dir", "'$DIR'"]}' --scope user
claude mcp add-json access-fs-sdk '{"type":"stdio","command":"'$DIR'/mcp-filesystem-server-mark3labs-mcp-go","args":["-dir", "'$DIR'"]}' --scope user
}
```

let Claude Code to use the MCP server to list the file,

```bash
claude "use access-fs to list file in ."
```

The output will similar to

```
access-fs - list_directory (MCP)(path: ".")
  ⎿  Directory contents:
     .git/
     .gitignore
     … +9 lines (ctrl+r to expand)

⏺ - .git/
  - .gitignore
  - README.md
  - build.sh
  - cmd/
  - go.mod
  - go.sum
  - internal/
  - mcp-filesystem-server
  - mcp-filesystem-server-mark3labs-mcp-go
  - test.sh
```

## Prompts and Cost for this MCP server

### Pure Golang Version

Prompt Claude Code

```
Create an MCP server in Golang to give Claude access to files in specific directory. The MCP server can expose file system access through standardized tools.
```

### mark3labs/mcp-go version and refactor

Prompt for `v2: sdk version`

```
1. describe this project
2. add one more command mcp-filesystem-server-mark3labs-mcp-go (the previous is mcp-filesystem-server), implement a new mcp server provides same abilities via sdk                             github.com/mark3labs/mcp-go
3. rename run.sh to build.sh, use build.sh to build two binaries.
```

Prompt for `as standard Go layout`: separate the code of two binaries using Golang's classic cmd dir structure.

Prompt for `works in any dir`: make mcp servers to works in any directories, not just BaseDir. the directory can be as boot argument.

Prompt: update @test.sh to the latest version

### Cost

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
