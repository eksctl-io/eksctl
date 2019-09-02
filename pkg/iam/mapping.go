package iam

import "errors"

var (
	// ErrNeitherUserNorRole is the error returned when an identity is missing both UserARN
	// and RoleARN.
	ErrNeitherUserNorRole = errors.New("arn is neither user nor role")

	// ErrNoKubernetesIdentity is the error returned when an identity has neither a Kubernetes
	// username nor a list of groups.
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
	if i.UserARN == nil && i.RoleARN == nil {
		return ErrNeitherUserNorRole
	}

	if i.Username == nil && len(i.Groups) == 0 {
		return ErrNoKubernetesIdentity
	}
	return nil
}

// ARN returns either the identites UserARN, RoleARN or ErrNeitherUserNorRole
func (i *Identity) ARN() (*ARN, error) {
	if i.UserARN != nil {
		return i.UserARN, nil
	} else if i.RoleARN != nil {
		return i.RoleARN, nil
	} else {
		return nil, ErrNeitherUserNorRole
	}
}

// NewIdentity determines into which field the given arn goes and returns the new identity
// alongside any error resulting for checking its validity.
func NewIdentity(arn ARN, username string, groups []string) (*Identity, error) {
	identity := Identity{}

	if arn.User() {
		identity.UserARN = &arn
	} else if arn.Role() {
		identity.RoleARN = &arn
	} else {
		return nil, ErrNeitherUserNorRole
	}

	if username != "" {
		identity.Username = &username
	}

	if len(groups) > 0 {
		identity.Groups = groups
	}

	return &identity, identity.Valid()
}
