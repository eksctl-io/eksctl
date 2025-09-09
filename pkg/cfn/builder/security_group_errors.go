package builder

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/smithy-go"
	"github.com/kris-nova/logger"
)

// SecurityGroupErrorHandler provides utilities for handling security group related errors
type SecurityGroupErrorHandler struct {
	clusterName string
}

// NewSecurityGroupErrorHandler creates a new SecurityGroupErrorHandler
func NewSecurityGroupErrorHandler(clusterName string) *SecurityGroupErrorHandler {
	return &SecurityGroupErrorHandler{
		clusterName: clusterName,
	}
}

// WrapTemplateGenerationError wraps template generation errors with actionable messages
func (h *SecurityGroupErrorHandler) WrapTemplateGenerationError(operation string, err error) error {
	if err == nil {
		return nil
	}

	logger.Debug("security group template generation error for cluster %q during %s: %v", h.clusterName, operation, err)

	// Check for specific AWS errors and provide actionable messages
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.ErrorCode() {
		case "InvalidParameterValue":
			return fmt.Errorf("failed to %s for cluster %q: invalid parameter value - %s. "+
				"Please check your cluster configuration and ensure all required fields are properly set",
				operation, h.clusterName, apiErr.ErrorMessage())
		case "InvalidVpcId.NotFound":
			return fmt.Errorf("failed to %s for cluster %q: VPC not found - %s. "+
				"Please ensure the VPC exists and you have the necessary permissions to access it",
				operation, h.clusterName, apiErr.ErrorMessage())
		case "UnauthorizedOperation":
			return fmt.Errorf("failed to %s for cluster %q: insufficient permissions - %s. "+
				"Please ensure your IAM role has the following permissions: ec2:CreateSecurityGroup, ec2:CreateTags, ec2:DescribeSecurityGroups",
				operation, h.clusterName, apiErr.ErrorMessage())
		}
	}

	// Check for common configuration issues
	if strings.Contains(err.Error(), "metadata") {
		return fmt.Errorf("failed to %s for cluster %q: %w. "+
			"This error typically occurs when cluster metadata is missing or invalid. "+
			"Please ensure your cluster configuration includes valid metadata with the required tags",
			operation, h.clusterName, err)
	}

	if strings.Contains(err.Error(), "karpenter.sh/discovery") {
		return fmt.Errorf("failed to %s for cluster %q: %w. "+
			"To enable automatic security group tagging, ensure both conditions are met: "+
			"1) Karpenter is enabled (karpenter.version is specified), "+
			"2) karpenter.sh/discovery tag is present in metadata.tags",
			operation, h.clusterName, err)
	}

	// Generic error with troubleshooting guidance
	return fmt.Errorf("failed to %s for cluster %q: %w. "+
		"Please check your cluster configuration and ensure you have the necessary IAM permissions: "+
		"ec2:CreateSecurityGroup, ec2:CreateTags, ec2:DescribeSecurityGroups",
		operation, h.clusterName, err)
}

// WrapConfigurationError wraps configuration validation errors with helpful guidance
func (h *SecurityGroupErrorHandler) WrapConfigurationError(err error) error {
	if err == nil {
		return nil
	}

	logger.Debug("security group configuration error for cluster %q: %v", h.clusterName, err)

	if strings.Contains(err.Error(), "nil") {
		return fmt.Errorf("invalid security group configuration for cluster %q: %w. "+
			"Please ensure your cluster configuration includes valid metadata and tags",
			h.clusterName, err)
	}

	if strings.Contains(err.Error(), "empty") {
		return fmt.Errorf("invalid security group configuration for cluster %q: %w. "+
			"The karpenter.sh/discovery tag value cannot be empty. "+
			"Please provide a valid tag value, typically the cluster name",
			h.clusterName, err)
	}

	if strings.Contains(err.Error(), "length") {
		return fmt.Errorf("invalid security group configuration for cluster %q: %w. "+
			"AWS tag values must be 256 characters or less. "+
			"Please use a shorter value for the karpenter.sh/discovery tag",
			h.clusterName, err)
	}

	return fmt.Errorf("invalid security group configuration for cluster %q: %w", h.clusterName, err)
}

// LogSuccessfulTagging logs successful security group tagging operations
func (h *SecurityGroupErrorHandler) LogSuccessfulTagging(tagKey, tagValue string) {
	logger.Info("successfully configured security group tagging for cluster %q", h.clusterName)
	logger.Debug("added tag %s=%s to cluster shared node security group for cluster %q", tagKey, tagValue, h.clusterName)
}

// LogSkippedTagging logs when security group tagging is skipped
func (h *SecurityGroupErrorHandler) LogSkippedTagging(reason string) {
	logger.Debug("skipping security group tagging for cluster %q: %s", h.clusterName, reason)
}

// ValidateTaggingPrerequisites validates that all prerequisites for security group tagging are met
func (h *SecurityGroupErrorHandler) ValidateTaggingPrerequisites(hasKarpenter, hasDiscoveryTag bool, discoveryValue string) error {
	if !hasKarpenter && !hasDiscoveryTag {
		h.LogSkippedTagging("neither Karpenter configuration nor karpenter.sh/discovery tag found")
		return nil // This is not an error, just skip tagging
	}

	if !hasKarpenter {
		h.LogSkippedTagging("Karpenter is not enabled (karpenter.version not specified)")
		return nil // This is not an error, just skip tagging
	}

	if !hasDiscoveryTag {
		h.LogSkippedTagging("karpenter.sh/discovery tag not found in metadata.tags")
		return nil // This is not an error, just skip tagging
	}

	if discoveryValue == "" {
		return fmt.Errorf("karpenter.sh/discovery tag value cannot be empty")
	}

	if len(discoveryValue) > 256 {
		return fmt.Errorf("karpenter.sh/discovery tag value exceeds maximum length of 256 characters")
	}

	return nil
}
