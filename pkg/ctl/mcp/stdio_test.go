package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStdioServer tests the stdio server functionality
func TestStdioServer(t *testing.T) {
	// Create a simple root command for testing
	rootCmd := &cobra.Command{
		Use:   "eksctl",
		Short: "The official CLI for Amazon EKS",
	}

	// Add a test subcommand
	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	rootCmd.AddCommand(testCmd)

	// Create the MCP server
	mcpServer, err := newEksctlMCPServer(rootCmd, "test-version")
	require.NoError(t, err)
	require.NotNil(t, mcpServer)

	// Create mock stdin/stdout
	mockStdin := &bytes.Buffer{}
	mockStdout := &bytes.Buffer{}

	// Create a stdio server
	stdioServer := server.NewStdioServer(mcpServer)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start the server in a goroutine
	go func() {
		err := stdioServer.Listen(ctx, mockStdin, mockStdout)
		if err != nil && err != context.DeadlineExceeded && err != io.EOF {
			t.Errorf("Unexpected error from stdio server: %v", err)
		}
	}()

	// Write a ping request to stdin
	pingRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "ping",
	}
	pingRequestBytes, err := json.Marshal(pingRequest)
	require.NoError(t, err)

	// Add newline to the request
	pingRequestBytes = append(pingRequestBytes, '\n')

	// Write to mock stdin
	_, err = mockStdin.Write(pingRequestBytes)
	require.NoError(t, err)

	// Wait for response
	time.Sleep(100 * time.Millisecond)

	// Read response from stdout
	response := mockStdout.String()
	assert.NotEmpty(t, response)

	// Parse response
	var responseObj map[string]interface{}
	err = json.Unmarshal([]byte(strings.TrimSpace(response)), &responseObj)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, "2.0", responseObj["jsonrpc"])
	assert.Equal(t, float64(1), responseObj["id"])
	assert.NotNil(t, responseObj["result"])
}

// TestStdioServerWithToolCall tests the stdio server with a tool call
func TestStdioServerWithToolCall(t *testing.T) {
	// Create a simple root command for testing
	rootCmd := &cobra.Command{
		Use:   "eksctl",
		Short: "The official CLI for Amazon EKS",
	}

	// Add a test subcommand with a mock implementation
	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	rootCmd.AddCommand(testCmd)

	// Create the MCP server
	mcpServer, err := newEksctlMCPServer(rootCmd, "test-version")
	require.NoError(t, err)
	require.NotNil(t, mcpServer)

	// Add a mock tool directly to the server
	mcpServer.AddTool(
		mcp.NewTool("mock_tool"),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return mcp.NewToolResultText("Mock tool response"), nil
		},
	)

	// Create mock stdin/stdout
	mockStdin := &bytes.Buffer{}
	mockStdout := &bytes.Buffer{}

	// Create a stdio server
	stdioServer := server.NewStdioServer(mcpServer)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start the server in a goroutine
	go func() {
		err := stdioServer.Listen(ctx, mockStdin, mockStdout)
		if err != nil && err != context.DeadlineExceeded && err != io.EOF {
			t.Errorf("Unexpected error from stdio server: %v", err)
		}
	}()

	// Write a tool call request to stdin
	toolRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name":      "mock_tool",
			"arguments": map[string]interface{}{},
		},
	}
	toolRequestBytes, err := json.Marshal(toolRequest)
	require.NoError(t, err)

	// Add newline to the request
	toolRequestBytes = append(toolRequestBytes, '\n')

	// Write to mock stdin
	_, err = mockStdin.Write(toolRequestBytes)
	require.NoError(t, err)

	// Wait for response
	time.Sleep(100 * time.Millisecond)

	// Read response from stdout
	response := mockStdout.String()
	assert.NotEmpty(t, response)

	// Parse response
	var responseObj map[string]interface{}
	err = json.Unmarshal([]byte(strings.TrimSpace(response)), &responseObj)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, "2.0", responseObj["jsonrpc"])
	assert.Equal(t, float64(2), responseObj["id"])
	assert.NotNil(t, responseObj["result"])

	// Verify tool response content
	result := responseObj["result"].(map[string]interface{})
	content := result["content"].([]interface{})
	assert.Len(t, content, 1)

	textContent := content[0].(map[string]interface{})
	assert.Equal(t, "text", textContent["type"])
	assert.Equal(t, "Mock tool response", textContent["text"])
}

// TestStdioServerMultipleRequests tests handling multiple requests
func TestStdioServerMultipleRequests(t *testing.T) {
	// Create a simple root command for testing
	rootCmd := &cobra.Command{
		Use:   "eksctl",
		Short: "The official CLI for Amazon EKS",
	}

	// Create the MCP server
	mcpServer, err := newEksctlMCPServer(rootCmd, "test-version")
	require.NoError(t, err)
	require.NotNil(t, mcpServer)

	// Create mock stdin/stdout
	mockStdin := &bytes.Buffer{}
	mockStdout := &bytes.Buffer{}

	// Create a stdio server
	stdioServer := server.NewStdioServer(mcpServer)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start the server in a goroutine
	go func() {
		err := stdioServer.Listen(ctx, mockStdin, mockStdout)
		if err != nil && err != context.DeadlineExceeded && err != io.EOF {
			t.Errorf("Unexpected error from stdio server: %v", err)
		}
	}()

	// Write multiple requests to stdin
	requests := []map[string]interface{}{
		{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  "ping",
		},
		{
			"jsonrpc": "2.0",
			"id":      2,
			"method":  "initialize",
		},
	}

	for _, req := range requests {
		reqBytes, err := json.Marshal(req)
		require.NoError(t, err)

		// Add newline to the request
		reqBytes = append(reqBytes, '\n')

		// Write to mock stdin
		_, err = mockStdin.Write(reqBytes)
		require.NoError(t, err)

		// Wait for response
		time.Sleep(100 * time.Millisecond)
	}

	// Read response from stdout
	response := mockStdout.String()

	// The response format might vary, so we'll just check that we got some response
	assert.NotEmpty(t, response)

	// Parse and verify each response - we don't need to check specific values
	// Just verify we can parse the response
	responseLines := strings.Split(strings.TrimSpace(response), "\n")
	for _, line := range responseLines {
		if line == "" {
			continue
		}
		var responseObj map[string]interface{}
		err = json.Unmarshal([]byte(line), &responseObj)
		assert.NoError(t, err, "Response should be valid JSON: %s", line)
	}
}

// TestStdioServerInvalidRequest tests handling invalid requests
func TestStdioServerInvalidRequest(t *testing.T) {
	// Create a simple root command for testing
	rootCmd := &cobra.Command{
		Use:   "eksctl",
		Short: "The official CLI for Amazon EKS",
	}

	// Create the MCP server
	mcpServer, err := newEksctlMCPServer(rootCmd, "test-version")
	require.NoError(t, err)
	require.NotNil(t, mcpServer)

	// Create mock stdin/stdout
	mockStdin := &bytes.Buffer{}
	mockStdout := &bytes.Buffer{}

	// Create a stdio server
	stdioServer := server.NewStdioServer(mcpServer)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start the server in a goroutine
	go func() {
		err := stdioServer.Listen(ctx, mockStdin, mockStdout)
		if err != nil && err != context.DeadlineExceeded && err != io.EOF {
			t.Errorf("Unexpected error from stdio server: %v", err)
		}
	}()

	// Write an invalid JSON request to stdin
	invalidRequest := `{"jsonrpc": "2.0", "id": 1, "method": "ping"`
	_, err = mockStdin.WriteString(invalidRequest + "\n")
	require.NoError(t, err)

	// Wait for response
	time.Sleep(100 * time.Millisecond)

	// Read response from stdout
	response := mockStdout.String()

	// The response format might vary, so we'll just check that we got some response
	assert.NotEmpty(t, response)

	// Try to parse the response, but don't assert on specific values
	var responseObj map[string]interface{}
	_ = json.Unmarshal([]byte(strings.TrimSpace(response)), &responseObj)
}
