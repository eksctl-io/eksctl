package iam

import (
	"errors"
	"fmt"
)

var (
	// ErrNeitherUserNorRole is the error returned when an identity is missing both UserARN
	// and RoleARN.
	ErrNeitherUserNorRole = errors.New("arn is neither user nor role")

	// ErrNoKubernetesIdentity is the error returned when an identity has neither a Kubernetes
	// username nor a list of groups.
	ErrNoKubernetesIdentity = errors.New("neither username nor groups is set for iam identity")
)

// Identity represents an IAM identity and its corresponding Kubernetes identity
type Identity interface {
	GetARN() string
	Type() string
	GetUsername() string
	GetGroups() []string
}

// KubernetesIdentity represents a kubernetes identity to be used in iam mappings
type KubernetesIdentity struct {
	Username string   `json:"username,omitempty"`
	Groups   []string `json:"groups,omitempty"`
}

// UserIdentity represents a mapping from an IAM user to a kubernetes identity
type UserIdentity struct {
	UserARN string `json:"userarn,omitempty"`
	KubernetesIdentity
}

// RoleIdentity represents a mapping from an IAM role to a kubernetes identity
type RoleIdentity struct {
	RoleARN string `json:"rolearn,omitempty"`
	KubernetesIdentity
}

// GetUsername returns the Kubernetes username
func (k KubernetesIdentity) GetUsername() string {
	return k.Username
}

// GetGroups returns the Kubernetes groups
func (k KubernetesIdentity) GetGroups() []string {
	return k.Groups
}

// GetARN returns the ARN of the iam mapping
func (u UserIdentity) GetARN() string {
	return u.UserARN
}

// Type returns the resource type of the iam mapping
func (u UserIdentity) Type() string {
	return ResourceTypeUser
}

// GetARN returns the ARN of the iam mapping
func (r RoleIdentity) GetARN() string {
	return r.RoleARN
}

// Type returns the resource type of the iam mapping
func (r RoleIdentity) Type() string {
	return ResourceTypeRole
}

// NewIdentity determines into which field the given arn goes and returns the new identity
// alongside any error resulting for checking its validity.
func NewIdentity(arn string, username string, groups []string) (Identity, error) {
	if arn == "" {
		return nil, fmt.Errorf("expected a valid arn but got empty string")
	}
	if username == "" && len(groups) == 0 {
		return nil, ErrNoKubernetesIdentity
	}

	parsedARN, err := Parse(arn)
	if err != nil {
		return nil, err
	}

	switch {
	case parsedARN.IsUser():
		return &UserIdentity{
			UserARN: arn,
			KubernetesIdentity: KubernetesIdentity{
				Username: username,
				Groups:   groups,
			},
		}, nil
	case parsedARN.IsRole():
		return &RoleIdentity{
			RoleARN: arn,
			KubernetesIdentity: KubernetesIdentity{
				Username: username,
				Groups:   groups,
			},
		}, nil
	default:
		return nil, ErrNeitherUserNorRole
	}
}
