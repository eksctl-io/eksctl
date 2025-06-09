package mcp

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEksctlMCPServer(t *testing.T) {
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

	// Verify server creation
	assert.NoError(t, err)
	assert.NotNil(t, mcpServer)

	// We can't directly access name and version fields, but we can verify the server was created
	require.NotNil(t, mcpServer)
}

func TestCommand(t *testing.T) {
	// Create a root command
	rootCmd := &cobra.Command{
		Use:   "eksctl",
		Short: "The official CLI for Amazon EKS",
	}

	// Create the MCP command
	mcpCmd := Command(rootCmd)

	// Verify command properties
	assert.Equal(t, "mcp", mcpCmd.Use)
	assert.True(t, mcpCmd.Hidden)
	assert.NotEmpty(t, mcpCmd.Short)
	assert.NotEmpty(t, mcpCmd.Long)
}
