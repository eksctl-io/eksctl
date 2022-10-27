package v1alpha5

import (
	"fmt"
	"strings"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
)

// Addon holds the EKS addon configuration
type Addon struct {
	// +required
	Name string `json:"name,omitempty"`
	// +optional
	Version string `json:"version,omitempty"`
	// +optional
	ServiceAccountRoleARN string `json:"serviceAccountRoleARN,omitempty"`
	// list of ARNs of the IAM policies to attach
	// +optional
	AttachPolicyARNs []string `json:"attachPolicyARNs,omitempty"`
	// AttachPolicy holds a policy document to attach
	// +optional
	AttachPolicy InlineDocument `json:"attachPolicy,omitempty"`
	// ARN of the permissions' boundary to associate
	// +optional
	PermissionsBoundary string `json:"permissionsBoundary,omitempty"`
	// WellKnownPolicies for attaching common IAM policies
	WellKnownPolicies WellKnownPolicies `json:"wellKnownPolicies,omitempty"`
	// The metadata to apply to the cluster to assist with categorization and organization.
	// Each tag consists of a key and an optional value, both of which you define.
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
	// ResolveConflicts determines how to resolve field value conflicts for an EKS add-on
	// if a value was changed from default
	ResolveConflicts ekstypes.ResolveConflicts `json:"resolveConflicts,omitempty"`
	// Force overwrites an existing self-managed add-on with an EKS managed add-on.
	// Force is intended to be used when migrating an existing self-managed add-on to an EKS managed add-on.
	Force bool `json:"-"`
}

func (a Addon) CanonicalName() string {
	return strings.ToLower(a.Name)
}

func (a Addon) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("name required")
	}

	return a.checkOnlyOnePolicyProviderIsSet()
}

func (a Addon) checkOnlyOnePolicyProviderIsSet() error {
	setPolicyProviders := 0
	if a.AttachPolicy != nil {
		setPolicyProviders++
	}

	if a.AttachPolicyARNs != nil && len(a.AttachPolicyARNs) > 0 {
		setPolicyProviders++
	}

	if a.ServiceAccountRoleARN != "" {
		setPolicyProviders++
	}

	if a.WellKnownPolicies.HasPolicy() {
		setPolicyProviders++
	}

	if setPolicyProviders > 1 {
		return fmt.Errorf("at most one of wellKnownPolicies, serviceAccountRoleARN, attachPolicyARNs and attachPolicy can be specified")
	}
	return nil
}
