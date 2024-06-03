package v1alpha5

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
)

// AccessEntry represents an access entry for managing access to a cluster.
type AccessEntry struct {
	// existing IAM principal ARN to associate with an access entry
	PrincipalARN ARN `json:"principalARN"`
	// `EC2_LINUX`, `EC2_WINDOWS`, `FARGATE_LINUX` or `STANDARD`
	// +optional
	Type string `json:"type,omitempty"`
	// set of Kubernetes groups to map to the principal ARN
	// +optional
	KubernetesGroups []string `json:"kubernetesGroups,omitempty"`
	// username to map to the principal ARN
	// +optional
	KubernetesUsername string `json:"kubernetesUsername,omitempty"`
	// set of policies to associate with an access entry
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
	// `namespace` or `cluster`
	Type ekstypes.AccessScopeType `json:"type"`
	// Scope access to namespace(s)
	// +optional
	Namespaces []string `json:"namespaces,omitempty"`
}

// AccessEntryType represents the type of access entry.
type AccessEntryType string

const (
	// AccessEntryTypeLinux specifies the EC2 Linux access entry type.
	AccessEntryTypeLinux AccessEntryType = "EC2_LINUX"
	// AccessEntryTypeWindows specifies the Windows access entry type.
	AccessEntryTypeWindows AccessEntryType = "EC2_WINDOWS"
	// AccessEntryTypeFargateLinux specifies the Fargate Linux access entry type.
	AccessEntryTypeFargateLinux AccessEntryType = "FARGATE_LINUX"
	// AccessEntryTypeStandard specifies a standard access entry type.
	AccessEntryTypeStandard AccessEntryType = "STANDARD"
)

// GetAccessEntryType returns the access entry type for the specified AMI family.
func GetAccessEntryType(ng *NodeGroup) AccessEntryType {
	if IsWindowsImage(ng.GetAMIFamily()) {
		return AccessEntryTypeWindows
	}
	return AccessEntryTypeLinux
}

type ARN arn.ARN

// ARN provides custom unmarshalling for an AWS ARN.

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

// RoleNameFromARN returns the role name for roleARN.
func RoleNameFromARN(roleARN string) (string, error) {
	parsed, err := arn.Parse(roleARN)
	if err != nil {
		return "", err
	}
	parts := strings.Split(parsed.Resource, "/")
	if len(parts) != 2 {
		return "", errors.New("invalid format for role ARN")
	}
	if parts[0] != "role" {
		return "", fmt.Errorf("expected resource type to be %q; got %q", "role", parts[0])
	}
	return parts[1], nil
}

// validateAccessEntries validates accessEntries.
func validateAccessEntries(accessEntries []AccessEntry) error {
	seen := make(map[ARN]struct{})
	for i, ae := range accessEntries {
		path := fmt.Sprintf("accessEntries[%d]", i)
		if ae.PrincipalARN.IsZero() {
			return fmt.Errorf("%s.principalARN must be set to a valid AWS ARN", path)
		}

		switch AccessEntryType(ae.Type) {
		case "", AccessEntryTypeStandard:
		case AccessEntryTypeLinux, AccessEntryTypeWindows, AccessEntryTypeFargateLinux:
			if len(ae.KubernetesGroups) > 0 || ae.KubernetesUsername != "" {
				return fmt.Errorf("cannot specify %s.kubernetesGroups nor %s.kubernetesUsername when type is set to %s", path, path, ae.Type)
			}
			if len(ae.AccessPolicies) > 0 {
				return fmt.Errorf("cannot specify %s.accessPolicies when type is set to %s", path, ae.Type)
			}
		default:
			return fmt.Errorf("invalid access entry type %q for %s", ae.Type, path)
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

			switch typ := ap.AccessScope.Type; typ {
			case "":
				return fmt.Errorf("%s.accessScope.type must be set to either %q or %q", path, ekstypes.AccessScopeTypeNamespace, ekstypes.AccessScopeTypeCluster)
			case ekstypes.AccessScopeTypeCluster:
				if len(ap.AccessScope.Namespaces) > 0 {
					return fmt.Errorf("cannot specify %s.accessScope.namespaces when accessScope is set to %s", path, typ)
				}
			case ekstypes.AccessScopeTypeNamespace:
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
