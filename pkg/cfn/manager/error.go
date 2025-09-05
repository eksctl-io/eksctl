package manager

import "fmt"

type StackNotFoundErr struct {
	ClusterName string
}

func (e *StackNotFoundErr) Error() string {
	return fmt.Sprintf("no eksctl-managed CloudFormation stacks found for %q", e.ClusterName)
}

// SecurityGroupTaggingError represents errors related to security group tagging during CloudFormation operations
type SecurityGroupTaggingError struct {
	ClusterName string
	Operation   string
	Cause       error
}

func (e *SecurityGroupTaggingError) Error() string {
	return fmt.Sprintf("failed to %s security group tags for cluster %q: %v", e.Operation, e.ClusterName, e.Cause)
}

func (e *SecurityGroupTaggingError) Unwrap() error {
	return e.Cause
}

// NewSecurityGroupTaggingError creates a new SecurityGroupTaggingError
func NewSecurityGroupTaggingError(clusterName, operation string, cause error) *SecurityGroupTaggingError {
	return &SecurityGroupTaggingError{
		ClusterName: clusterName,
		Operation:   operation,
		Cause:       cause,
	}
}

// CloudFormationTemplateError represents errors during CloudFormation template generation
type CloudFormationTemplateError struct {
	ClusterName string
	Component   string
	Cause       error
}

func (e *CloudFormationTemplateError) Error() string {
	return fmt.Sprintf("failed to generate CloudFormation template for %s in cluster %q: %v", e.Component, e.ClusterName, e.Cause)
}

func (e *CloudFormationTemplateError) Unwrap() error {
	return e.Cause
}

// NewCloudFormationTemplateError creates a new CloudFormationTemplateError
func NewCloudFormationTemplateError(clusterName, component string, cause error) *CloudFormationTemplateError {
	return &CloudFormationTemplateError{
		ClusterName: clusterName,
		Component:   component,
		Cause:       cause,
	}
}

// SecurityGroupConfigurationError represents configuration validation errors for security group tagging
type SecurityGroupConfigurationError struct {
	ClusterName string
	Issue       string
}

func (e *SecurityGroupConfigurationError) Error() string {
	return fmt.Sprintf("invalid security group tagging configuration for cluster %q: %s", e.ClusterName, e.Issue)
}

// NewSecurityGroupConfigurationError creates a new SecurityGroupConfigurationError
func NewSecurityGroupConfigurationError(clusterName, issue string) *SecurityGroupConfigurationError {
	return &SecurityGroupConfigurationError{
		ClusterName: clusterName,
		Issue:       issue,
	}
}
