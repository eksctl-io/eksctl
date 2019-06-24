package authconfigmap

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/iam"
)

// MapIdentity represents an IAM identity with an ARN.
type MapIdentity struct {
	iam.Identity `json:",inline"`
	ARN          ARN `json:"-"` // This field is (un)marshaled manually
}

// UnmarshalJSON for MapIdentity makes sure there is either a "rolearn" or "userarn" key
// and places the data of that key into the "ARN" field. All other fields are handled by
// the default implementation of unmarshal.
func (m *MapIdentity) UnmarshalJSON(data []byte) error {
	// Handle everything as usual
	type alias MapIdentity
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	// We want to unmarshal "(rolearn|userarn)" into the "ARN" field and then unmarshal
	// the rest as usual
	var outerKeys map[string]*json.RawMessage
	if err := json.Unmarshal(data, &outerKeys); err != nil {
		return err
	}

	arnData, ok := outerKeys[roleKey]
	if !ok {
		arnData, ok = outerKeys[userKey]
		if !ok {
			return errors.New("missing arn")
		}
	}

	var canonicalArn string
	if err := json.Unmarshal(*arnData, &canonicalArn); err != nil {
		return err
	}

	arn, err := arn.Parse(canonicalArn)
	if err != nil {
		return err
	}
	a.ARN = ARN(arn)

	*m = MapIdentity(a)
	return nil
}

// MarshalJSON for MapIdentity marshals the ARN field into either "rolearn" or "userarn"
// depending on what is appropriate and returns and error if it cannot determine that.
// All other fields are marshaled with the default marshaler
func (m *MapIdentity) MarshalJSON() ([]byte, error) {
	// Marshal everything as one would normally do
	type alias MapIdentity
	a := (*alias)(m)
	data, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	// Partially unmarshal everything into a map in order to more easily append the (role|user)arn
	partial := map[string]*json.RawMessage{}
	if err := json.Unmarshal(data, &partial); err != nil {
		return nil, err
	}

	arn, err := json.Marshal(m.ARN.String())
	if err != nil {
		return nil, err
	}

	switch {
	case m.role():
		partial[roleKey] = (*json.RawMessage)(&arn)
	case m.user():
		partial[userKey] = (*json.RawMessage)(&arn)
	default:
		return nil, errors.Errorf("cannot determine if %q refers to a user or role", arn)
	}

	return json.Marshal(partial)
}

func (m MapIdentity) resource() string {
	resource := m.ARN.Resource
	if idx := strings.Index(resource, "/"); idx >= 0 {
		resource = resource[:idx] // remove everything following the forward slash
	}

	return resource
}

func (m MapIdentity) role() bool {
	return m.resource() == "role"
}

func (m MapIdentity) user() bool {
	return m.resource() == "user"
}

// MapIdentities is a list of IAM identities with a role or user ARN.
type MapIdentities []MapIdentity

// Get returns all matching role mappings. Note that at this moment
// aws-iam-authenticator only considers the last one!
func (rs MapIdentities) Get(arn ARN) MapIdentities {
	var m MapIdentities
	for _, r := range rs {
		if r.ARN.String() == arn.String() {
			m = append(m, r)
		}
	}
	return m
}
