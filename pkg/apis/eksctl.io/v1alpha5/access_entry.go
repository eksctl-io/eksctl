package v1alpha5

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
)

// AccessEntry represents an access entry for managing access to a cluster.
type AccessEntry struct {
	PrincipalARN ARN `json:"principalARN"`
	// +optional
	KubernetesGroups []string `json:"kubernetesGroups,omitempty"`
	// +optional
	KubernetesUsername string `json:"kubernetesUsername,omitempty"`
	// +optional
	AccessPolicies []AccessPolicy `json:"accessPolicies,omitempty"`
}

// An AccessPolicy represents a policy to associate with an access entry.
type AccessPolicy struct {
	PolicyARN   ARN         `json:"policyARN"`
	AccessScope AccessScope `json:"accessScope"`
}

// AccessScope defines the scope of an access policy.
type AccessScope struct {
	Type string `json:"type"`
	// +optional
	Namespaces []string `json:"namespaces,omitempty"`
}

// ARN provides custom unmarshalling for an AWS ARN.
type ARN arn.ARN

// UnmarshalText implements encoding.TextUnmarshaler.
func (a *ARN) UnmarshalText(arnStr []byte) error {
	return a.set(string(arnStr))
}

// Set implements pflag.Value.
func (a *ARN) Set(arnStr string) error {
	return a.set(arnStr)
}

// String returns the string representation of the ARN.
func (a ARN) String() string {
	return arn.ARN(a).String()
}

// Type returns the type.
func (a *ARN) Type() string {
	return "string"
}

// MarshalJSON implements json.Marshaler.
func (a ARN) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// IsZero reports whether a is the zero value.
func (a ARN) IsZero() bool {
	return a.Partition == ""
}

func (a *ARN) set(arnStr string) error {
	parsed, err := arn.Parse(arnStr)
	if err != nil {
		return fmt.Errorf("invalid ARN %q: %w", arnStr, err)
	}
	*a = ARN(parsed)
	return nil
}

// MustParseARN returns the parsed ARN or panics if the ARN cannot be parsed.
func MustParseARN(a string) ARN {
	parsed, err := arn.Parse(a)
	if err != nil {
		panic(err)
	}
	return ARN(parsed)
}

// validateAccessEntries validates accessEntries.
func validateAccessEntries(accessEntries []AccessEntry) error {
	seen := make(map[ARN]struct{})
	for i, ae := range accessEntries {
		path := fmt.Sprintf("accessEntries[%d]", i)
		if ae.PrincipalARN.IsZero() {
			return fmt.Errorf("%s.principalARN must be set to a valid AWS ARN", path)
		}

		for _, ap := range ae.AccessPolicies {
			if ap.PolicyARN.IsZero() {
				return fmt.Errorf("%s.policyARN must be set to a valid AWS ARN", path)
			}

			if parts := strings.Split(ap.PolicyARN.Resource, "/"); len(parts) > 1 {
				if parts[0] != "cluster-access-policy" {
					return fmt.Errorf("%s.policyARN must be a cluster-access-policy resource", path)
				}
			} else {
				return fmt.Errorf("invalid %s.policyARN", path)
			}

			// TODO: use SDK enums.
			switch typ := ap.AccessScope.Type; typ {
			case "":
				return fmt.Errorf("%s.accessScope.type must be set to either %q or %q", path, "namespace", "cluster")
			case "cluster":
				if len(ap.AccessScope.Namespaces) > 0 {
					return fmt.Errorf("cannot specify %s.accessScope.namespaces when accessScope is set to %s", path, typ)
				}
			case "namespace":
				if len(ap.AccessScope.Namespaces) == 0 {
					return fmt.Errorf("at least one namespace must be specified when accessScope is set to %s: (%s)", typ, path)
				}
			default:
				return fmt.Errorf("invalid access scope type %q for %s", typ, path)
			}
		}
		if _, exists := seen[ae.PrincipalARN]; exists {
			return fmt.Errorf("duplicate access entry %s with principal ARN %q", path, ae.PrincipalARN.String())
		}
		seen[ae.PrincipalARN] = struct{}{}
	}
	return nil
}
