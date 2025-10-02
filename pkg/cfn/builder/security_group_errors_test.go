package builder

import (
	"errors"
	"strings"
	"testing"

	"github.com/aws/smithy-go"
	"github.com/stretchr/testify/assert"
)

func TestSecurityGroupErrorHandler_WrapTemplateGenerationError(t *testing.T) {
	handler := NewSecurityGroupErrorHandler("test-cluster")

	tests := []struct {
		name          string
		operation     string
		inputError    error
		expectedError string
	}{
		{
			name:          "nil error returns nil",
			operation:     "generate tags",
			inputError:    nil,
			expectedError: "",
		},
		{
			name:          "InvalidParameterValue AWS error",
			operation:     "generate tags",
			inputError:    &smithy.GenericAPIError{Code: "InvalidParameterValue", Message: "Invalid parameter"},
			expectedError: "failed to generate tags for cluster \"test-cluster\": invalid parameter value",
		},
		{
			name:          "UnauthorizedOperation AWS error",
			operation:     "create security group",
			inputError:    &smithy.GenericAPIError{Code: "UnauthorizedOperation", Message: "Access denied"},
			expectedError: "failed to create security group for cluster \"test-cluster\": insufficient permissions",
		},
		{
			name:          "metadata related error",
			operation:     "generate tags",
			inputError:    errors.New("cluster metadata is nil"),
			expectedError: "failed to generate tags for cluster \"test-cluster\": cluster metadata is nil",
		},
		{
			name:          "karpenter discovery related error",
			operation:     "configure tags",
			inputError:    errors.New("karpenter.sh/discovery tag not found"),
			expectedError: "failed to configure tags for cluster \"test-cluster\": karpenter.sh/discovery tag not found",
		},
		{
			name:          "generic error",
			operation:     "process template",
			inputError:    errors.New("unknown error"),
			expectedError: "failed to process template for cluster \"test-cluster\": unknown error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.WrapTemplateGenerationError(tt.operation, tt.inputError)

			if tt.expectedError == "" {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Contains(t, result.Error(), tt.expectedError)
			}
		})
	}
}

func TestSecurityGroupErrorHandler_WrapConfigurationError(t *testing.T) {
	handler := NewSecurityGroupErrorHandler("test-cluster")

	tests := []struct {
		name          string
		inputError    error
		expectedError string
	}{
		{
			name:          "nil error returns nil",
			inputError:    nil,
			expectedError: "",
		},
		{
			name:          "nil related error",
			inputError:    errors.New("cluster metadata is nil"),
			expectedError: "invalid security group configuration for cluster \"test-cluster\"",
		},
		{
			name:          "empty value error",
			inputError:    errors.New("tag value cannot be empty"),
			expectedError: "invalid security group configuration for cluster \"test-cluster\"",
		},
		{
			name:          "length related error",
			inputError:    errors.New("exceeds maximum length"),
			expectedError: "invalid security group configuration for cluster \"test-cluster\"",
		},
		{
			name:          "generic configuration error",
			inputError:    errors.New("invalid configuration"),
			expectedError: "invalid security group configuration for cluster \"test-cluster\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.WrapConfigurationError(tt.inputError)

			if tt.expectedError == "" {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Contains(t, result.Error(), tt.expectedError)
			}
		})
	}
}

func TestSecurityGroupErrorHandler_ValidateTaggingPrerequisites(t *testing.T) {
	handler := NewSecurityGroupErrorHandler("test-cluster")

	tests := []struct {
		name            string
		hasKarpenter    bool
		hasDiscoveryTag bool
		discoveryValue  string
		expectError     bool
		expectedError   string
	}{
		{
			name:            "neither condition met - no error",
			hasKarpenter:    false,
			hasDiscoveryTag: false,
			discoveryValue:  "",
			expectError:     false,
		},
		{
			name:            "only Karpenter enabled - no error",
			hasKarpenter:    true,
			hasDiscoveryTag: false,
			discoveryValue:  "",
			expectError:     false,
		},
		{
			name:            "only discovery tag present - no error",
			hasKarpenter:    false,
			hasDiscoveryTag: true,
			discoveryValue:  "test-cluster",
			expectError:     false,
		},
		{
			name:            "both conditions met with valid value - no error",
			hasKarpenter:    true,
			hasDiscoveryTag: true,
			discoveryValue:  "test-cluster",
			expectError:     false,
		},
		{
			name:            "both conditions met but empty value - error",
			hasKarpenter:    true,
			hasDiscoveryTag: true,
			discoveryValue:  "",
			expectError:     true,
			expectedError:   "tag value cannot be empty",
		},
		{
			name:            "both conditions met but value too long - error",
			hasKarpenter:    true,
			hasDiscoveryTag: true,
			discoveryValue:  strings.Repeat("a", 257), // 257 characters
			expectError:     true,
			expectedError:   "exceeds maximum length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.ValidateTaggingPrerequisites(tt.hasKarpenter, tt.hasDiscoveryTag, tt.discoveryValue)

			if tt.expectError {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestNewSecurityGroupErrorHandler(t *testing.T) {
	clusterName := "test-cluster"
	handler := NewSecurityGroupErrorHandler(clusterName)

	assert.NotNil(t, handler)
	assert.Equal(t, clusterName, handler.clusterName)
}
