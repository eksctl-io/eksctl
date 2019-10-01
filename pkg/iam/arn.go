package iam

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
)

const (
	// ResourceTypeRole is the resource type of the role ARN
	ResourceTypeRole = "role"
	// ResourceTypeUser is the resource type of the user ARN
	ResourceTypeUser = "user"
)

// ARN implements the pflag.Value interface for aws-sdk-go/aws/arn.ARN
type ARN struct {
	arn.ARN
}

// Parse wraps the aws-sdk-go/aws/arn.Parse function and instead returns a
// iam.ARN
func Parse(s string) (ARN, error) {
	a, err := arn.Parse(s)
	return ARN{a}, err
}

// ResourceType returns the type of the resource specified in the ARN.
// Typically, in the case of IAM, it is a role or a user
func (a *ARN) ResourceType() string {
	t := a.Resource
	if idx := strings.Index(t, "/"); idx >= 0 {
		t = t[:idx] // remove everything following the forward slash
	}

	return t
}

// IsUser returns whether the arn represents a IAM user or not
func (a *ARN) IsUser() bool {
	return a.ResourceType() == ResourceTypeUser
}

// IsRole returns whether the arn represents a IAM role or not
func (a *ARN) IsRole() bool {
	return a.ResourceType() == ResourceTypeRole
}
