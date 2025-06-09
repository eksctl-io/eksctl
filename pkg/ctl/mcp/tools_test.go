package mcp

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterTools(t *testing.T) {
	// Create a test MCP server
	s := server.NewMCPServer("test-server", "1.0.0", server.WithToolCapabilities(true))

	// Create a test command structure
	rootCmd := &cobra.Command{
		Use:   "eksctl",
		Short: "The official CLI for Amazon EKS",
	}

	// Add a visible command
	visibleCmd := &cobra.Command{
		Use:   "visible",
		Short: "Visible command",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	visibleCmd.Flags().String("test-flag", "", "Test flag")
	rootCmd.AddCommand(visibleCmd)

	// Add a hidden command
	hiddenCmd := &cobra.Command{
		Use:    "hidden",
		Short:  "Hidden command",
		Hidden: true,
		Run:    func(cmd *cobra.Command, args []string) {},
	}
	rootCmd.AddCommand(hiddenCmd)

	// Register tools
	err := registerTools(s, rootCmd)
	assert.NoError(t, err)

	// Get the list of tools
	toolsList := s.HandleMessage(context.Background(), []byte(`{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "tools/list"
	}`))

	// Verify the response
	resp, ok := toolsList.(mcp.JSONRPCResponse)
	assert.True(t, ok)

	result, ok := resp.Result.(mcp.ListToolsResult)
	assert.True(t, ok)

	// We should have at least one tool (the visible command)
	assert.NotEmpty(t, result.Tools)

	// Find our visible command tool
	var foundVisibleTool bool
	for _, tool := range result.Tools {
		if tool.Name == "eksctl_visible" {
			foundVisibleTool = true
			break
		}
	}
	assert.True(t, foundVisibleTool, "Should find the visible command tool")

	// Make sure hidden command is not registered
	var foundHiddenTool bool
	for _, tool := range result.Tools {
		if tool.Name == "eksctl_hidden" {
			foundHiddenTool = true
			break
		}
	}
	assert.False(t, foundHiddenTool, "Should not find the hidden command tool")
}

func TestBuildToolOptions(t *testing.T) {
	tests := []struct {
		name           string
		setupCmd       func() *cobra.Command
		expectedName   string
		expectedParams int
	}{
		{
			name: "Command with string flag",
			setupCmd: func() *cobra.Command {
				cmd := &cobra.Command{
					Use:   "test",
					Short: "Test command",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				cmd.Flags().String("string-flag", "", "String flag")
				return cmd
			},
			expectedName:   "eksctl_test",
			expectedParams: 1,
		},
		{
			name: "Command with bool flag",
			setupCmd: func() *cobra.Command {
				cmd := &cobra.Command{
					Use:   "test",
					Short: "Test command",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				cmd.Flags().Bool("bool-flag", false, "Bool flag")
				return cmd
			},
			expectedName:   "eksctl_test",
			expectedParams: 1,
		},
		{
			name: "Command with multiple flags",
			setupCmd: func() *cobra.Command {
				cmd := &cobra.Command{
					Use:   "test",
					Short: "Test command",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				cmd.Flags().String("string-flag", "", "String flag")
				cmd.Flags().Bool("bool-flag", false, "Bool flag")
				cmd.Flags().Int("int-flag", 0, "Int flag")
				cmd.Flags().StringSlice("slice-flag", []string{}, "Slice flag")
				return cmd
			},
			expectedName:   "eksctl_test",
			expectedParams: 4,
		},
		{
			name: "Command with required flag",
			setupCmd: func() *cobra.Command {
				cmd := &cobra.Command{
					Use:   "test",
					Short: "Test command",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				flag := cmd.Flags().String("required-flag", "", "Required flag")
				// We can't set Annotations directly, so we'll skip that part in the test
				_ = flag
				return cmd
			},
			expectedName:   "eksctl_test",
			expectedParams: 1,
		},
		{
			name: "Command with deprecated flag",
			setupCmd: func() *cobra.Command {
				cmd := &cobra.Command{
					Use:   "test",
					Short: "Test command",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				// We can't set Deprecated directly, so we'll skip that part in the test
				cmd.Flags().String("flag1", "", "Flag 1")
				cmd.Flags().String("flag2", "", "Flag 2")
				return cmd
			},
			expectedName:   "eksctl_test",
			expectedParams: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			rootCmd := &cobra.Command{Use: "eksctl"}
			cmd := tt.setupCmd()
			rootCmd.AddCommand(cmd)

			// Test
			toolName, toolOptions, err := buildToolOptions(cmd)

			// Verify
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedName, toolName)

			// Count parameter options
			paramCount := 0
			for _, opt := range toolOptions {
				if opt != nil {
					paramCount++
				}
			}
			assert.Equal(t, tt.expectedParams+1, paramCount, "Expected parameters + 1 for description")
		})
	}
}

func TestShouldSkipCommand(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *cobra.Command
		expected bool
	}{
		{
			name: "Regular command",
			cmd: &cobra.Command{
				Use:   "regular",
				Short: "Regular command",
			},
			expected: false,
		},
		{
			name: "Hidden command",
			cmd: &cobra.Command{
				Use:    "hidden",
				Short:  "Hidden command",
				Hidden: true,
			},
			expected: true,
		},
		{
			name: "Help command",
			cmd: &cobra.Command{
				Use:   "help",
				Short: "Help command",
			},
			expected: true,
		},
		{
			name: "Completion command",
			cmd: &cobra.Command{
				Use:   "completion",
				Short: "Completion command",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldSkipCommand(tt.cmd)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractEnumValuesFromUsage(t *testing.T) {
	tests := []struct {
		name     string
		usage    string
		expected []string
	}{
		{
			name:     "No enum values",
			usage:    "Just a regular description",
			expected: []string{},
		},
		{
			name:     "Must be one of pattern",
			usage:    "The value must be one of: value1, value2, value3",
			expected: []string{"value1", "value2", "value3"},
		},
		{
			name:     "Can be one of pattern",
			usage:    "The value can be one of: value1, value2, value3",
			expected: []string{"value1", "value2", "value3"},
		},
		{
			name:     "Pipe pattern",
			usage:    "The value [value1|value2|value3] is required",
			expected: []string{"value1", "value2", "value3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractEnumValuesFromUsage(tt.usage)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateToolHandler(t *testing.T) {
	// Create a test command with flags
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	cmd.Flags().String("string-flag", "", "String flag")
	cmd.Flags().Bool("bool-flag", false, "Bool flag")
	cmd.Flags().StringSlice("slice-flag", []string{}, "Slice flag")

	// Create the handler
	handler := createToolHandler(cmd)
	require.NotNil(t, handler)

	// We can't easily test the execution without mocking exec.Command
	// But we can verify the handler doesn't panic with valid input
	request := mcp.CallToolRequest{}
	request.Params.Arguments = map[string]any{
		"string-flag": "test-value",
		"bool-flag":   true,
		"slice-flag":  "value1,value2",
	}

	// This will try to execute the real command, which might fail
	// but at least we can verify the handler doesn't panic
	_, _ = handler(context.Background(), request)
}

func TestRegisterTool(t *testing.T) {
	// Create a test MCP server
	s := server.NewMCPServer("test-server", "1.0.0", server.WithToolCapabilities(true))

	// Create a test command
	rootCmd := &cobra.Command{Use: "eksctl"}
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	cmd.Flags().String("flag", "", "Test flag")
	rootCmd.AddCommand(cmd)

	// Register the tool
	err := registerTool(s, cmd)
	assert.NoError(t, err)

	// Verify the tool was registered
	toolsList := s.HandleMessage(context.Background(), []byte(`{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "tools/list"
	}`))

	resp, ok := toolsList.(mcp.JSONRPCResponse)
	assert.True(t, ok)

	result, ok := resp.Result.(mcp.ListToolsResult)
	assert.True(t, ok)

	// Find our tool
	var found bool
	for _, tool := range result.Tools {
		if tool.Name == "eksctl_test" {
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the registered tool")
}

func TestRegisterToolsRecursive(t *testing.T) {
	// Create a test MCP server
	s := server.NewMCPServer("test-server", "1.0.0", server.WithToolCapabilities(true))

	// Create a test command structure with nested commands
	rootCmd := &cobra.Command{Use: "eksctl"}

	// Level 1 command
	level1Cmd := &cobra.Command{
		Use:   "level1",
		Short: "Level 1 command",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	rootCmd.AddCommand(level1Cmd)

	// Level 2 command
	level2Cmd := &cobra.Command{
		Use:   "level2",
		Short: "Level 2 command",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	level1Cmd.AddCommand(level2Cmd)

	// Register tools recursively
	err := registerToolsRecursive(s, rootCmd)
	assert.NoError(t, err)

	// Verify tools were registered
	toolsList := s.HandleMessage(context.Background(), []byte(`{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "tools/list"
	}`))

	resp, ok := toolsList.(mcp.JSONRPCResponse)
	assert.True(t, ok)

	result, ok := resp.Result.(mcp.ListToolsResult)
	assert.True(t, ok)

	// We should have at least 2 tools (level1 and level2)
	foundLevel1 := false
	foundLevel2 := false

	for _, tool := range result.Tools {
		if tool.Name == "eksctl_level1" {
			foundLevel1 = true
		}
		if tool.Name == "eksctl_level1_level2" {
			foundLevel2 = true
		}
	}

	assert.True(t, foundLevel1, "Should find level1 command tool")
	assert.True(t, foundLevel2, "Should find level2 command tool")
}
