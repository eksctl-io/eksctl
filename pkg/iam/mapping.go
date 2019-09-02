package iam

import "errors"

var (
	ErrNeitherUserNorRole   = errors.New("arn is neither user nor role")
	ErrNoKubernetesIdentity = errors.New("neither username nor groups is set for iam identity")
)

// Identity represents an IAM identity and its corresponding Kubernetes identity
type Identity struct {
	UserARN  *ARN     `json:"userarn,omitempty"`
	RoleARN  *ARN     `json:"rolearn,omitempty"`
	Username *string  `json:"username,omitempty"`
	Groups   []string `json:"groups,omitempty"`
}

// Valid ensures the identity is proper.
func (i Identity) Valid() error {
	if len(i.Groups) == 0 {
		return errors.New("identity mapping needs at least 1 group")
	}
	return nil
}

func NewIdentity(arn ARN, username string, groups []string) (Identity, error) {
	identity := Identity{}

	if arn.User() {
		identity.UserARN = &arn
	} else if arn.Role() {
		identity.RoleARN = &arn
	} else {
		return Identity{}, ErrNeitherUserNorRole
	}

	if username != "" {
		identity.Username = &username
	}

	if len(groups) > 0 {
		identity.Groups = groups
	}

	if identity.Username == nil && len(identity.Groups) == 0 {
		return Identity{}, ErrNoKubernetesIdentity
	}

	return identity, nil
}
