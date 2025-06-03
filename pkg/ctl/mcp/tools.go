// Package mcp implements the Model Context Protocol (MCP) server functionality for eksctl
package mcp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// registerTools registers all eksctl commands as MCP tools
// This is the entry point for tool registration that processes the entire command tree
func registerTools(s *server.MCPServer, rootCmd *cobra.Command) error {
	return registerToolsRecursive(s, rootCmd)
}

// registerToolsRecursive recursively registers all commands as tools
// It processes the current command and then recursively processes all subcommands
func registerToolsRecursive(s *server.MCPServer, cmd *cobra.Command) error {
	if err := registerTool(s, cmd); err != nil {
		return err
	}

	// Process all subcommands recursively
	for _, subCmd := range cmd.Commands() {
		if err := registerToolsRecursive(s, subCmd); err != nil {
			return err
		}
	}

	return nil
}

// registerTool registers a single command as an MCP tool
// It builds the tool options and creates a handler for the command
func registerTool(s *server.MCPServer, cmd *cobra.Command) error {
	// Skip commands that shouldn't be exposed
	if shouldSkipCommand(cmd) {
		return nil
	}

	// Build tool name and options
	toolName, toolOptions, err := buildToolOptions(cmd)
	if err != nil {
		return fmt.Errorf("error building tool options for %s: %w", cmd.Name(), err)
	}

	if toolName == "" {
		return nil
	}

	// Create and register the tool
	tool := mcp.NewTool(toolName, toolOptions...)
	s.AddTool(tool, createToolHandler(cmd))

	return nil
}

// shouldSkipCommand determines if a command should be skipped for MCP tool registration
func shouldSkipCommand(cmd *cobra.Command) bool {
	return cmd.Hidden || cmd.Name() == "help" || cmd.Name() == "completion"
}

// buildToolOptions builds the MCP tool options from a cobra command
func buildToolOptions(cmd *cobra.Command) (string, []mcp.ToolOption, error) {
	// Build the command path (e.g., "create cluster")
	var cmdPath []string
	current := cmd
	for current != nil && current.Name() != "eksctl" {
		cmdPath = append([]string{current.Name()}, cmdPath...)
		current = current.Parent()
	}

	// Skip if no path (this is the root command)
	if len(cmdPath) == 0 {
		return "", nil, nil
	}

	commandPath := strings.Join(cmdPath, " ")
	toolName := "eksctl_" + strings.Join(cmdPath, "_")

	// Extract description from the command
	description := cmd.Short
	if description == "" {
		description = cmd.Long
	}
	if description == "" {
		description = "Run the eksctl " + commandPath + " command"
	}

	// Get usage information
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	flagGrouping := cmdutils.NewGrouping()

	if err := flagGrouping.Usage(cmd); err != nil {
		return "", nil, fmt.Errorf("error printing usage: %w", err)
	}

	// Create tool options starting with description
	toolOptions := []mcp.ToolOption{mcp.WithDescription(description + "\n\n" + buf.String())}

	// Add parameters based on command flags
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		// Skip deprecated flags
		if flag.Deprecated != "" {
			return
		}

		name := flag.Name
		usage := flag.Usage

		// Check for enum values in usage string
		_ = extractEnumValuesFromUsage(usage)

		switch flag.Value.Type() {
		case "bool":
			toolOptions = append(toolOptions, mcp.WithBoolean(
				name,
				mcp.Description(usage),
			))
		case "stringSlice", "stringArray":
			// Handle string arrays as regular strings with comma-separated values
			toolOptions = append(toolOptions, mcp.WithString(
				name,
				mcp.Description(usage+" (comma-separated values)"),
			))
		case "intSlice", "intArray":
			// Handle int arrays as regular strings with comma-separated values
			toolOptions = append(toolOptions, mcp.WithString(
				name,
				mcp.Description(usage+" (comma-separated values)"),
			))
		case "int", "int32", "int64":
			// Use string for numbers as well for simplicity
			toolOptions = append(toolOptions, mcp.WithString(
				name,
				mcp.Description(usage),
			))
		case "float", "float32", "float64":
			toolOptions = append(toolOptions, mcp.WithString(
				name,
				mcp.Description(usage),
			))
		default:
			// Default to string for all other types
			stringOpts := []mcp.PropertyOption{mcp.Description(usage)}

			// Mark required if the flag is required
			if flag.Annotations != nil {
				if _, required := flag.Annotations["cobra_annotation_required"]; required {
					stringOpts = append(stringOpts, mcp.Required())
				}
			}

			toolOptions = append(toolOptions, mcp.WithString(
				name,
				stringOpts...,
			))
		}
	})

	return toolName, toolOptions, nil
}

// createToolHandler creates a handler function for an MCP tool
func createToolHandler(cmd *cobra.Command) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Add timeout to context
		ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()

		// Build the command path
		var cmdPath []string
		current := cmd
		for current != nil && current.Name() != "eksctl" {
			cmdPath = append([]string{current.Name()}, cmdPath...)
			current = current.Parent()
		}

		// Build arguments for the eksctl command
		args := cmdPath

		// Add flag values from the request
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			name := flag.Name

			// Handle different flag types
			switch flag.Value.Type() {
			case "bool":
				if request.GetBool(name, false) {
					args = append(args, "--"+name)
				}
			case "stringSlice", "stringArray", "intSlice", "intArray":
				// Handle arrays as comma-separated values
				if value := request.GetString(name, ""); value != "" {
					values := strings.Split(value, ",")
					for _, v := range values {
						args = append(args, "--"+name, strings.TrimSpace(v))
					}
				}
			default:
				// Handle string and number types
				if value := request.GetString(name, ""); value != "" {
					args = append(args, "--"+name, value)
				}
			}
		})

		return executeEksctlCommand(ctx, args)
	}
}

func executeEksctlCommand(ctx context.Context, args []string) (*mcp.CallToolResult, error) {
	// Create a context for output collection with a 45-second timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	cmd := exec.Command("eksctl", args...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return mcp.NewToolResultError("Failed to create stdout pipe: " + err.Error()), nil
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return mcp.NewToolResultError("Failed to create stderr pipe: " + err.Error()), nil
	}

	if err := cmd.Start(); err != nil {
		return mcp.NewToolResultError("Failed to start eksctl: " + err.Error()), nil
	}

	// Channels to collect output
	stdoutCh := make(chan string, 1)
	stderrCh := make(chan string, 1)

	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, stdoutPipe)
		stdoutCh <- buf.String()
	}()
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, stderrPipe)
		stderrCh <- buf.String()
	}()

	// Wait for either timeout or command completion
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-timeoutCtx.Done():
		// Timeout: collect whatever output is available
		stdout := ""
		select {
		case stdout = <-stdoutCh:
		default:
		}
		return mcp.NewToolResultText(stdout + "\n\n[Command is still running in the background]"), nil
	case err := <-done:
		// Command finished: collect output
		stdout := <-stdoutCh
		stderr := <-stderrCh
		if err != nil {
			if stderr == "" {
				stderr = err.Error()
			}
			return mcp.NewToolResultError(stderr), nil
		}
		return mcp.NewToolResultText(stdout), nil
	}
}

// Helper functions for parameter enhancement

// extractEnumValuesFromUsage extracts enum values from usage string
func extractEnumValuesFromUsage(usage string) []string {
	// Look for patterns like "must be one of: value1, value2, value3"
	oneOfPattern := regexp.MustCompile(`(?i)(?:must be|one of|can be)(?:\s+one of)?:\s+([^.]+)`)
	matches := oneOfPattern.FindStringSubmatch(usage)
	if len(matches) > 1 {
		values := strings.Split(matches[1], ",")
		for i, v := range values {
			values[i] = strings.Trim(v, " \t\n\r")
		}
		return values
	}

	// Look for patterns like "[value1|value2|value3]"
	pipePattern := regexp.MustCompile(`\[([^\]]+)\]`)
	matches = pipePattern.FindStringSubmatch(usage)
	if len(matches) > 1 && strings.Contains(matches[1], "|") {
		values := strings.Split(matches[1], "|")
		for i, v := range values {
			values[i] = strings.Trim(v, " \t\n\r")
		}
		return values
	}

	return []string{}
}
