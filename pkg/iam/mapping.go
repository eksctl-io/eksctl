package iam

import (
	"errors"
	"fmt"
	"sort"
)

const (
	// ResourceTypeAccount is the resource type of Accounts
	ResourceTypeAccount = "account"
)

var (
	// ErrNeitherUserNorRole is the error returned when an identity is missing both UserARN
	// and RoleARN.
	ErrNeitherUserNorRole = errors.New("arn is neither user nor role")

	// ErrNoKubernetesIdentity is the error returned when an identity has neither a Kubernetes
	// username nor a list of groups.
	ErrNoKubernetesIdentity = errors.New("neither username nor group are set for iam identity")
)

// Identity represents an IAM identity and its corresponding Kubernetes identity
type Identity interface {
	ARN() string
	Type() string
	Username() string
	Groups() []string
	Account() string
}

// CompareIdentity takes 2 Identity values and checks to see if they are identitcal
func CompareIdentity(a, b Identity) bool {
	sameAccount := a.ARN() == b.ARN() &&
		a.Type() == b.Type() &&
		a.Username() == b.Username() &&
		a.Account() == b.Account()

	if !sameAccount {
		return false
	}

	aGroups := a.Groups()
	bGroups := b.Groups()
	if len(aGroups) != len(bGroups) {
		return false
	}

	sort.Strings(aGroups)
	sort.Strings(bGroups)
	for i := range aGroups {
		if aGroups[i] != bGroups[i] {
			return false
		}
	}

	return true

}

// KubernetesIdentity represents a kubernetes identity to be used in iam mappings
type KubernetesIdentity struct {
	KubernetesUsername string   `json:"username,omitempty"`
	KubernetesGroups   []string `json:"groups,omitempty"`
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

// AccountIdentity represents a mapping from an IAM role to a kubernetes identity
type AccountIdentity struct {
	KubernetesAccount string `json:"account,omitempty"`
	KubernetesIdentity
}

// ARN returns the ARN of the iam mapping
func (a AccountIdentity) ARN() string {
	return ""
}

// Account returns the Account of the iam mapping
func (a AccountIdentity) Account() string {
	return a.KubernetesAccount
}

// Type returns the resource type of the iam mapping
func (a AccountIdentity) Type() string {
	return ResourceTypeAccount
}

// Username returns the Kubernetes username
func (k KubernetesIdentity) Username() string {
	return k.KubernetesUsername
}

// Groups returns the Kubernetes groups
func (k KubernetesIdentity) Groups() []string {
	return k.KubernetesGroups
}

// ARN returns the ARN of the iam mapping
func (u UserIdentity) ARN() string {
	return u.UserARN
}

// Type returns the resource type of the iam mapping
func (u UserIdentity) Type() string {
	return ResourceTypeUser
}

// Account returns the Account of the iam mapping
func (u UserIdentity) Account() string {
	return ""
}

// ARN returns the ARN of the iam mapping
func (r RoleIdentity) ARN() string {
	return r.RoleARN
}

// Account returns the Account of the iam mapping
func (r RoleIdentity) Account() string {
	return ""
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
				KubernetesUsername: username,
				KubernetesGroups:   groups,
			},
		}, nil
	case parsedARN.IsRole():
		return &RoleIdentity{
			RoleARN: arn,
			KubernetesIdentity: KubernetesIdentity{
				KubernetesUsername: username,
				KubernetesGroups:   groups,
			},
		}, nil
	default:
		return nil, ErrNeitherUserNorRole
	}
}
