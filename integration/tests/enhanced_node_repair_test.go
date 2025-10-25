//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEnhancedNodeRepairCLIFlags tests that the enhanced node repair CLI flags are properly parsed
func TestEnhancedNodeRepairCLIFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []string // Expected strings in the dry-run output
	}{
		{
			name: "basic node repair",
			args: []string{
				"create", "cluster",
				"--name", "test-basic-repair",
				"--enable-node-repair",
				"--dry-run",
			},
			expected: []string{
				"NodeRepairConfig",
				"Enabled\": true",
			},
		},
		{
			name: "node repair with percentage thresholds",
			args: []string{
				"create", "cluster",
				"--name", "test-percentage-repair",
				"--enable-node-repair",
				"--node-repair-max-unhealthy-percentage=25",
				"--node-repair-max-parallel-percentage=20",
				"--dry-run",
			},
			expected: []string{
				"NodeRepairConfig",
				"Enabled\": true",
				"MaxUnhealthyNodeThresholdPercentage\": 25",
				"MaxParallelNodesRepairedPercentage\": 20",
			},
		},
		{
			name: "node repair with count thresholds",
			args: []string{
				"create", "cluster",
				"--name", "test-count-repair",
				"--enable-node-repair",
				"--node-repair-max-unhealthy-count=5",
				"--node-repair-max-parallel-count=2",
				"--dry-run",
			},
			expected: []string{
				"NodeRepairConfig",
				"Enabled\": true",
				"MaxUnhealthyNodeThresholdCount\": 5",
				"MaxParallelNodesRepairedCount\": 2",
			},
		},
		{
			name: "nodegroup with node repair flags",
			args: []string{
				"create", "nodegroup",
				"--cluster", "existing-cluster",
				"--name", "test-ng-repair",
				"--enable-node-repair",
				"--node-repair-max-unhealthy-percentage=30",
				"--node-repair-max-parallel-count=1",
				"--dry-run",
			},
			expected: []string{
				"NodeRepairConfig",
				"Enabled\": true",
				"MaxUnhealthyNodeThresholdPercentage\": 30",
				"MaxParallelNodesRepairedCount\": 1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run eksctl with the test arguments
			cmd := exec.Command("./eksctl", tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				// For dry-run, we expect it to fail due to missing AWS credentials or cluster
				// but we should still get the configuration output
				t.Logf("Command failed as expected (dry-run): %v", err)
				t.Logf("Stderr: %s", stderr.String())
			}

			output := stdout.String() + stderr.String()
			t.Logf("Full output: %s", output)

			// Check that all expected strings are present in the output
			for _, expected := range tt.expected {
				assert.Contains(t, output, expected, "Expected string not found in output: %s", expected)
			}
		})
	}
}

// TestEnhancedNodeRepairConfigFile tests that enhanced node repair configuration files are properly parsed
func TestEnhancedNodeRepairConfigFile(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
		expected   []string
	}{
		{
			name:       "basic config file",
			configFile: "examples/44-node-repair.yaml",
			expected: []string{
				"NodeRepairConfig",
				"Enabled\": true",
			},
		},
		{
			name:       "enhanced config file",
			configFile: "examples/44-enhanced-node-repair.yaml",
			expected: []string{
				"NodeRepairConfig",
				"maxUnhealthyNodeThresholdPercentage",
				"maxParallelNodesRepairedPercentage",
				"nodeRepairConfigOverrides",
				"AcceleratedInstanceNotReady",
				"NvidiaXID13Error",
				"NetworkNotReady",
				"InterfaceNotUp",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run eksctl with the config file
			cmd := exec.Command("./eksctl", "create", "cluster", "--config-file", tt.configFile, "--dry-run")
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				// For dry-run, we expect it to fail due to missing AWS credentials
				// but we should still get the configuration output
				t.Logf("Command failed as expected (dry-run): %v", err)
			}

			output := stdout.String() + stderr.String()
			t.Logf("Full output: %s", output)

			// Check that all expected strings are present in the output
			for _, expected := range tt.expected {
				assert.Contains(t, output, expected, "Expected string not found in output: %s", expected)
			}
		})
	}
}

// TestEnhancedNodeRepairCLIHelp tests that the CLI help includes the new flags
func TestEnhancedNodeRepairCLIHelp(t *testing.T) {
	cmd := exec.Command("../../eksctl", "create", "cluster", "--help")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	require.NoError(t, err, "Help command should not fail")

	output := stdout.String()

	// Check that all new flags are documented in help
	expectedFlags := []string{
		"--enable-node-repair",
		"--node-repair-max-unhealthy-percentage",
		"--node-repair-max-unhealthy-count",
		"--node-repair-max-parallel-percentage",
		"--node-repair-max-parallel-count",
	}

	for _, flag := range expectedFlags {
		assert.Contains(t, output, flag, "Flag not found in help output: %s", flag)
	}

	// Check that flags have proper descriptions
	assert.Contains(t, output, "managed nodegroups only", "Flags should indicate they're for managed nodegroups only")
}

// TestEnhancedNodeRepairBackwardCompatibility tests that existing configurations still work
func TestEnhancedNodeRepairBackwardCompatibility(t *testing.T) {
	// Test that the original example still works
	cmd := exec.Command("../../eksctl", "create", "cluster", "--config-file", "../../examples/44-node-repair.yaml", "--dry-run")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Expected to fail due to missing AWS credentials, but should parse config
		t.Logf("Command failed as expected (dry-run): %v", err)
	}

	output := stdout.String() + stderr.String()

	// Should not contain any parsing errors
	assert.NotContains(t, strings.ToLower(output), "unknown field", "Should not have unknown field errors")
	assert.NotContains(t, strings.ToLower(output), "invalid", "Should not have invalid field errors")

	// Should contain the basic node repair config
	assert.Contains(t, output, "nodeRepairConfig", "Should contain nodeRepairConfig")
}

// TestEnhancedNodeRepairSchemaValidation tests that the schema includes new fields
func TestEnhancedNodeRepairSchemaValidation(t *testing.T) {
	cmd := exec.Command("../../eksctl", "utils", "schema")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	require.NoError(t, err, "Schema command should not fail")

	output := stdout.String()

	// Check that the schema includes all new fields
	expectedFields := []string{
		"NodeGroupNodeRepairConfig",
		"maxUnhealthyNodeThresholdPercentage",
		"maxUnhealthyNodeThresholdCount",
		"maxParallelNodesRepairedPercentage",
		"maxParallelNodesRepairedCount",
		"nodeRepairConfigOverrides",
	}

	for _, field := range expectedFields {
		assert.Contains(t, output, field, "Schema should include field: %s", field)
	}
}



// TestEnhancedNodeRepairErrorHandling tests error handling for invalid configurations
func TestEnhancedNodeRepairErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedError string
	}{
		{
			name: "conflicting percentage and count thresholds",
			args: []string{
				"create", "cluster",
				"--name", "test-conflict",
				"--enable-node-repair",
				"--node-repair-max-unhealthy-percentage=25",
				"--node-repair-max-unhealthy-count=5",
				"--dry-run",
			},
			expectedError: "cannot specify both percentage and count",
		},
		{
			name: "conflicting parallel percentage and count",
			args: []string{
				"create", "cluster",
				"--name", "test-parallel-conflict",
				"--enable-node-repair",
				"--node-repair-max-parallel-percentage=20",
				"--node-repair-max-parallel-count=2",
				"--dry-run",
			},
			expectedError: "cannot specify both percentage and count",
		},
		{
			name: "invalid percentage value too high",
			args: []string{
				"create", "cluster",
				"--name", "test-invalid-percentage",
				"--enable-node-repair",
				"--node-repair-max-unhealthy-percentage=150",
				"--dry-run",
			},
			expectedError: "percentage must be between 1 and 100",
		},
		{
			name: "invalid percentage value zero",
			args: []string{
				"create", "cluster",
				"--name", "test-zero-percentage",
				"--enable-node-repair",
				"--node-repair-max-unhealthy-percentage=0",
				"--dry-run",
			},
			expectedError: "percentage must be between 1 and 100",
		},
		{
			name: "invalid count value zero",
			args: []string{
				"create", "cluster",
				"--name", "test-zero-count",
				"--enable-node-repair",
				"--node-repair-max-unhealthy-count=0",
				"--dry-run",
			},
			expectedError: "count must be greater than 0",
		},
		{
			name: "node repair flags without enable flag",
			args: []string{
				"create", "cluster",
				"--name", "test-no-enable",
				"--node-repair-max-unhealthy-percentage=25",
				"--dry-run",
			},
			expectedError: "node repair flags require --enable-node-repair",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run eksctl with the test arguments
			cmd := exec.Command("../../eksctl", tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			require.Error(t, err, "Command should fail with invalid configuration")

			output := stdout.String() + stderr.String()
			t.Logf("Full output: %s", output)

			// Check that the expected error message is present
			assert.Contains(t, strings.ToLower(output), strings.ToLower(tt.expectedError), 
				"Expected error message not found: %s", tt.expectedError)
		})
	}
}

// TestEnhancedNodeRepairConfigFileErrorHandling tests error handling for invalid config files
func TestEnhancedNodeRepairConfigFileErrorHandling(t *testing.T) {
	// Create a temporary invalid config file
	invalidConfig := `
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: test-invalid-config
  region: us-west-2
managedNodeGroups:
- name: ng-1
  nodeRepairConfig:
    enabled: true
    maxUnhealthyNodeThresholdPercentage: 25
    maxUnhealthyNodeThresholdCount: 5  # This conflicts with percentage
    maxParallelNodesRepairedPercentage: 150  # Invalid percentage > 100
`

	tmpFile, err := os.CreateTemp("", "invalid-config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(invalidConfig)
	require.NoError(t, err)
	tmpFile.Close()

	// Run eksctl with the invalid config file
	cmd := exec.Command("../../eksctl", "create", "cluster", "--config-file", tmpFile.Name(), "--dry-run")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	require.Error(t, err, "Command should fail with invalid configuration")

	output := stdout.String() + stderr.String()
	t.Logf("Full output: %s", output)

	// Check for validation errors
	expectedErrors := []string{
		"cannot specify both percentage and count",
		"percentage must be between 1 and 100",
	}

	for _, expectedError := range expectedErrors {
		assert.Contains(t, strings.ToLower(output), strings.ToLower(expectedError), 
			"Expected error message not found: %s", expectedError)
	}
}

// TestEnhancedNodeRepairUnmanagedNodegroupError tests that node repair flags are rejected for unmanaged nodegroups
func TestEnhancedNodeRepairUnmanagedNodegroupError(t *testing.T) {
	// Create a config with unmanaged nodegroup and node repair config
	invalidConfig := `
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: test-unmanaged-repair
  region: us-west-2
nodeGroups:  # Unmanaged nodegroup
- name: ng-1
  nodeRepairConfig:  # This should be invalid for unmanaged nodegroups
    enabled: true
`

	tmpFile, err := os.CreateTemp("", "unmanaged-repair-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(invalidConfig)
	require.NoError(t, err)
	tmpFile.Close()

	// Run eksctl with the invalid config file
	cmd := exec.Command("../../eksctl", "create", "cluster", "--config-file", tmpFile.Name(), "--dry-run")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	require.Error(t, err, "Command should fail with node repair config on unmanaged nodegroup")

	output := stdout.String() + stderr.String()
	t.Logf("Full output: %s", output)

	// Check that the error mentions managed nodegroups only
	assert.Contains(t, strings.ToLower(output), "managed nodegroups only", 
		"Should indicate that node repair is for managed nodegroups only")
}

// TestEnhancedNodeRepairValidationRecovery tests that validation errors don't leave resources in inconsistent states
func TestEnhancedNodeRepairValidationRecovery(t *testing.T) {
	// This test ensures that when validation fails, no partial resources are created
	// Since we're using --dry-run, we're mainly testing that the validation happens early
	// and doesn't proceed to resource creation
	
	cmd := exec.Command("../../eksctl", 
		"create", "cluster",
		"--name", "test-validation-recovery",
		"--enable-node-repair",
		"--node-repair-max-unhealthy-percentage=200", // Invalid percentage
		"--dry-run",
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	require.Error(t, err, "Command should fail early with validation error")

	output := stdout.String() + stderr.String()
	t.Logf("Full output: %s", output)

	// Ensure that validation happens before any resource creation attempts
	assert.Contains(t, strings.ToLower(output), "percentage must be between 1 and 100", 
		"Should show validation error")
	
	// Should not contain CloudFormation template generation messages
	assert.NotContains(t, strings.ToLower(output), "cloudformation template", 
		"Should not proceed to CloudFormation template generation")
	assert.NotContains(t, strings.ToLower(output), "creating stack", 
		"Should not proceed to stack creation")
}