package iam

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
)

// ARN implements the pflag.Value interface for aws-sdk-go/aws/arn.ARN
type ARN struct {
	arn.ARN
}

// Parse wraps the aws-sdk-go/aws/arn.Parse function and instead returns a
// authconfigmap.ARN
func Parse(s string) (ARN, error) {
	a, err := arn.Parse(s)
	return ARN{a}, err
}

// Set parses the given string into an arn.ARN and sets the receiver pointer to the
// populated struct
func (a *ARN) Set(s string) error {
	arn, err := arn.Parse(s)
	if err != nil {
		return err
	}
	*a = ARN{arn}
	return nil
}

// Type describes the argument type in the pflag.Value interface
func (a *ARN) Type() string {
	return "aws arn"
}

func (a *ARN) resource() string {
	resource := a.Resource
	if idx := strings.Index(resource, "/"); idx >= 0 {
		resource = resource[:idx] // remove everything following the forward slash
	}

	return resource
}

// User returns whether the arn represents a IAM user or not
func (a *ARN) User() bool {
	return a.resource() == "user"
}

// Role returns whether the arn represents a IAM role or not
func (a *ARN) Role() bool {
	return a.resource() == "role"
}
