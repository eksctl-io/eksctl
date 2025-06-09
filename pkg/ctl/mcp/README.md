# eksctl MCP (Model Context Protocol) Package

## Overview

The `pkg/ctl/mcp` package implements a [Model Context Protocol (MCP)](https://github.com/mark3labs/mcp) server for eksctl. This server enables AI assistants like Amazon Q to interact with eksctl functionality directly through a standardized protocol, allowing for seamless integration of eksctl commands into AI-powered workflows.

## What is MCP?

Model Context Protocol (MCP) is an open protocol that standardizes how applications provide context and tools to Large Language Models (LLMs). It enables AI assistants to:

1. Discover available tools and their capabilities
2. Execute commands and receive structured responses
3. Provide contextual information to enhance AI interactions

## How the eksctl MCP Server Works

The eksctl MCP server exposes all eksctl commands as tools that can be invoked by AI assistants. The implementation consists of several key components:

### Core Components

1. **`mcp.go`**: Entry point that defines the `mcp` command and starts the MCP server
2. **`server.go`**: Creates and configures the MCP server with all eksctl commands
3. **`tools.go`**: Handles the registration of eksctl commands as MCP tools
4. **`types.go`**: Defines types and interfaces used by the MCP implementation

### Implementation Details

- The server uses the `mcp-go` library to implement the Model Context Protocol
- It dynamically registers all eksctl commands (create, get, delete, etc.) as MCP tools
- Each command's help text, flags, and parameters are exposed through the protocol
- The server runs as a stdio server, communicating through standard input/output

### Command Registration Process

1. The server recursively traverses the eksctl command tree
2. For each command, it extracts:
   - Command description and usage information
   - Available flags and their descriptions
   - Required and optional parameters

## Using the eksctl MCP Server with Amazon Q Chat

To use the eksctl MCP server with Amazon Q Chat, you need to register it in the Amazon Q configuration file.

Create or edit the file at `$HOME/.aws/amazonq/mcp.json` with the following content:

```json
{
  "mcpServers": {
    "eks-tools": {
      "command": "eksctl",
      "args": [
        "mcp"
      ]
    }
  }
}
```

This configuration tells Amazon Q Chat to:
1. Register a server named "eks-tools"
2. Use the `eksctl mcp` command to start the MCP server
3. Make all eksctl commands available as tools with the prefix `eks___`

## Benefits

- **Seamless Integration**: AI assistants can execute eksctl commands directly
- **Structured Responses**: Commands return structured data that can be parsed by AI models
- **Discoverability**: AI assistants can discover available commands and their parameters
- **Context-Aware**: Provides rich context about EKS clusters and resources

## Development

When extending the MCP server functionality:

1. Add new command registrations in `server.go` if new eksctl commands are created
2. Extend type definitions in `types.go` as needed
3. Modify `tools.go` if changes to tool registration logic are required

## Example Usage in Amazon Q Chat

Once configured, users can interact with eksctl through Amazon Q Chat:

```
User: What EKS clusters do I have?
Amazon Q: Let me check your EKS clusters.
[Amazon Q executes the eksctl get cluster command and displays the results]

User: Create a new EKS cluster
Amazon Q: [Guides the user through cluster creation using eksctl commands]
```

The MCP server enables these interactions by providing Amazon Q with the ability to execute eksctl commands and interpret their results.

## Quick test for MCP server
To quickly test the MCP server, you can run the following command in your terminal:

```bash
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | eksctl mcp | jq
```

## Recommendations
- Monitor background operations for long-running commands. Commands such as cluster creations have a 45-second timeout for command responses but the processes continue in background.
