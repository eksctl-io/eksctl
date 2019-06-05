package iam

import "errors"

// Identity represents an IAM identity.
type Identity struct {
	Username string   `json:"username"`
	Groups   []string `json:"groups"`
}

// Valid ensures the identity is proper.
func (i Identity) Valid() error {
	if len(i.Groups) == 0 {
		return errors.New("identity mapping needs at least 1 group")
	}
	return nil
}
