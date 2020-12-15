package v1alpha5

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type IdentityProviderType string

const (
	OIDCIdentityProviderType IdentityProviderType = "oidc"
)

// IdentityProviderInterface is a dummy interface
// to give some extra type safety
type IdentityProviderInterface interface {
	isIdentityProvider()
}

// The idea of the `IdentityProvider` struct is to hold an identity provider
// that can be parsed from the following JSON:
// {
// 	"name": "user-pool-1",
// 	"type": "oidc"
// }
// i.e. the type is found adjacent to the other fields of the object
//
// An `IdentityProvider` contains exactly one such identity provider
// which can be accessed with `Inner` and then cast and switched on
// with `.(type)` to get access to the specific type

// IdentityProvider holds an identity provider
// Schema type is one of `OIDCIdentityProvider`
type IdentityProvider struct {
	// Valid variants are:
	// `"oidc"`: OIDC identity provider
	// +required
	Type  string `json:"type"`
	inner IdentityProviderInterface
}

// DeepCopy is needed to generate kubernetes types for IdentityProvider
func (in *IdentityProvider) DeepCopy() *IdentityProvider {
	if in == nil {
		return nil
	}
	out := new(IdentityProvider)
	switch idP := in.inner.(type) {
	case *OIDCIdentityProvider:
		out.inner = idP.DeepCopy()
	default:
		panic("unknown inner identity provider in IdentityProvider")
	}
	return out
}

// Inner returns the contained identity provider
func (ip *IdentityProvider) Inner() IdentityProviderInterface {
	return ip.inner
}

func (ip *IdentityProvider) UnmarshalJSON(data []byte) error {
	var typ struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &typ); err != nil {
		return err
	}
	switch typ.Type {
	case string(OIDCIdentityProviderType):
		oidc := new(OIDCIdentityProvider)
		if err := json.Unmarshal(data, oidc); err != nil {
			return err
		}
		ip.inner = oidc
	default:
		return errors.New("couldn't unmarshal to IdentityProvider, invalid type")
	}
	return nil
}

// OIDCIdentityProvider holds the spec of an OIDC provider
// to use for EKS authzn
type OIDCIdentityProvider struct {
	// +required
	Name           string            `json:"name,omitempty"`
	// +required
	IssuerURL      string            `json:"issuerURL,omitempty"`
	// +required
	ClientID       string            `json:"clientID,omitempty"`
	UsernameClaim  string            `json:"usernameClaim,omitempty"`
	UsernamePrefix string            `json:"usernamePrefix,omitempty"`
	GroupsClaim    string            `json:"groupsClaim,omitempty"`
	GroupsPrefix   string            `json:"groupsPrefix,omitempty"`
	RequiredClaims map[string]string `json:"requiredClaims,omitempty"`
	Tags           map[string]string `json:"tags,omitempty"`
}

func (*OIDCIdentityProvider) isIdentityProvider() {}
