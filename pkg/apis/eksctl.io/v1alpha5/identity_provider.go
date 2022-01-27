package v1alpha5

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type IdentityProviderType string

const (
	OIDCIdentityProviderType IdentityProviderType = "oidc"
)

// IdentityProviderInterface is a dummy interface
// to give some extra type safety
type IdentityProviderInterface interface {
	DeepCopyIdentityProviderInterface() IdentityProviderInterface
	Type() IdentityProviderType
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

// IdentityProvider holds an identity provider configuration.
// See [the example eksctl config](https://github.com/weaveworks/eksctl/blob/main/examples/27-oidc-provider.yaml).
// Schema type is one of `OIDCIdentityProvider`
type IdentityProvider struct {
	// Valid variants are:
	// `"oidc"`: OIDC identity provider
	// +required
	type_ string `json:"type"` //nolint
	Inner IdentityProviderInterface
}

func FromIdentityProvider(idp IdentityProviderInterface) IdentityProvider {
	return IdentityProvider{
		type_: string(idp.Type()),
		Inner: idp,
	}
}

func (ip *IdentityProvider) MarshalJSON() ([]byte, error) {
	switch ip.Inner.Type() {
	case OIDCIdentityProviderType:
		oidc, ok := ip.Inner.(*OIDCIdentityProvider)
		if !ok {
			return nil, fmt.Errorf("failed to cast oidc of type %s into OIDCIdentityProvider", OIDCIdentityProviderType)
		}
		return json.Marshal(struct {
			*OIDCIdentityProvider
			Type string `json:"type"`
		}{
			Type:                 string(OIDCIdentityProviderType),
			OIDCIdentityProvider: oidc,
		})

	default:
		return nil, errors.New("couldn't marshal to IdentityProvider, invalid type")
	}
}

func (ip *IdentityProvider) UnmarshalJSON(data []byte) error {
	var typ struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &typ); err != nil {
		return err
	}
	var inner IdentityProviderInterface
	switch typ.Type {
	case string(OIDCIdentityProviderType):
		oidc := new(OIDCIdentityProvider)
		if err := json.Unmarshal(data, oidc); err != nil {
			return err
		}
		inner = oidc
	default:
		return errors.New("couldn't unmarshal to IdentityProvider, invalid type")
	}
	ip.Inner = inner
	return nil
}

// OIDCIdentityProvider holds the spec of an OIDC provider
// to use for EKS authzn
type OIDCIdentityProvider struct {
	// +required
	Name string `json:"name,omitempty"`
	// +required
	IssuerURL string `json:"issuerURL,omitempty"`
	// +required
	ClientID       string            `json:"clientID,omitempty"`
	UsernameClaim  string            `json:"usernameClaim,omitempty"`
	UsernamePrefix string            `json:"usernamePrefix,omitempty"`
	GroupsClaim    string            `json:"groupsClaim,omitempty"`
	GroupsPrefix   string            `json:"groupsPrefix,omitempty"`
	RequiredClaims map[string]string `json:"requiredClaims,omitempty"`
	Tags           map[string]string `json:"tags,omitempty"`
}

func (p *OIDCIdentityProvider) DeepCopyIdentityProviderInterface() IdentityProviderInterface {
	return p.DeepCopy()
}

func (p *OIDCIdentityProvider) Type() IdentityProviderType {
	return OIDCIdentityProviderType
}
