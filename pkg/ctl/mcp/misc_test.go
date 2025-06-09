package mcp

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/misc"
)

// TestMiscCommandRegistration tests that misc commands are properly registered
func TestMiscCommandRegistration(t *testing.T) {
	// Create a root command
	rootCmd := &cobra.Command{
		Use:   "eksctl",
		Short: "The official CLI for Amazon EKS",
	}

	// Create a flag grouping for the commands
	flagGrouping := cmdutils.NewGrouping()

	// Register misc commands to the root command
	misc.Command(flagGrouping, rootCmd)

	// Create the MCP server with the root command
	mcpServer, err := newEksctlMCPServer(rootCmd, "test-version")
	require.NoError(t, err)
	require.NotNil(t, mcpServer)

	// Verify that the version command was registered
	versionCmd := findCommandInTest(rootCmd, "version")
	assert.NotNil(t, versionCmd, "version command should be registered")

	// Verify that the info command was registered
	infoCmd := findCommandInTest(rootCmd, "info")
	assert.NotNil(t, infoCmd, "info command should be registered")
}

// findCommandInTest recursively searches for a command by name
func findCommandInTest(cmd *cobra.Command, name string) *cobra.Command {
	for _, subCmd := range cmd.Commands() {
		if subCmd.Name() == name {
			return subCmd
		}
		if found := findCommandInTest(subCmd, name); found != nil {
			return found
		}
	}
	return nil
}
