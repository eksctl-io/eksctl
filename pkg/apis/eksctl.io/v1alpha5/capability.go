package v1alpha5

import "fmt"

// Capability represents an EKS capability configuration
type Capability struct {
	// Name of the capability
	// +required
	Name string `json:"name"`

	// Type of the capability (ACK, KRO, ARGOCD)
	// +required
	Type string `json:"type"`

	// RoleARN is the IAM role ARN for the capability
	// +optional
	RoleARN string `json:"roleArn,omitempty"`

	// DeletePropagationPolicy specifies the delete propagation policy
	// +optional
	DeletePropagationPolicy string `json:"deletePropagationPolicy,omitempty"`

	// Configuration holds capability-specific configuration. Only applicable for ArgoCD
	// +optional
	Configuration *CapabilityConfiguration `json:"configuration,omitempty"`

	// Tags are used to tag AWS resources created by the capability
	// +optional
	Tags map[string]string `json:"tags,omitempty"`

	// AccessPolicies list of access policies to associate with the access entry
	// +optional
	AccessPolicies []AccessPolicy `json:"accessPolicies,omitempty"`

	// AttachPolicyARNs list of ARNs of the IAM policies to attach
	// +optional
	AttachPolicyARNs []string `json:"attachPolicyARNs,omitempty"`

	// AttachPolicy holds a policy document to attach
	// +optional
	AttachPolicy InlineDocument `json:"attachPolicy,omitempty"`

	// PermissionsBoundary is the ARN of the permissions boundary policy
	// +optional
	PermissionsBoundary string `json:"permissionsBoundary,omitempty"`
}

// CapabilityConfiguration holds capability-specific configuration
type CapabilityConfiguration struct {
	// ArgoCD configuration for ARGOCD capability type
	// +optional
	ArgoCD *ArgoCDConfiguration `json:"argocd,omitempty"`
}

// ArgoCDConfiguration holds ArgoCD-specific configuration
type ArgoCDConfiguration struct {
	// Namespace for ArgoCD installation
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// NetworkAccess configuration
	// +optional
	NetworkAccess *ArgoCDNetworkAccess `json:"networkAccess,omitempty"`

	// RBACRoleMappings for ArgoCD RBAC
	// +optional
	RBACRoleMappings []ArgoCDRoleMapping `json:"rbacRoleMappings,omitempty"`

	// AWSIDC configuration
	// +optional
	AWSIDC *ArgoCDAWSIDC `json:"awsIdc,omitempty"`
}

// ArgoCDNetworkAccess holds network access configuration for ArgoCD
type ArgoCDNetworkAccess struct {
	// VPCEIDs for VPC endpoint access
	// +optional
	VPCEIDs []string `json:"vpceIds,omitempty"`
}

// ArgoCDRoleMapping holds RBAC role mapping for ArgoCD
type ArgoCDRoleMapping struct {
	// Role is the ArgoCD role (ADMIN, EDITOR, VIEWER)
	// +required
	Role string `json:"role"`

	// Identities are the SSO identities to map to the role
	// +required
	Identities []SSOIdentity `json:"identities"`
}

// SSOIdentity represents an SSO identity
type SSOIdentity struct {
	// ID of the SSO identity
	// +required
	ID string `json:"id"`

	// Type of the SSO identity (SSO_USER, SSO_GROUP)
	// +required
	Type string `json:"type"`
}

// ArgoCDAWSIDC holds AWS IDC configuration for ArgoCD
type ArgoCDAWSIDC struct {
	// IDCInstanceARN is the ARN of the IDC instance
	// +required
	IDCInstanceARN string `json:"idcInstanceArn"`

	// IDCRegion is the region of the IDC instance
	// +optional
	IDCRegion string `json:"idcRegion,omitempty"`
}

// Validate validates the capability configuration
func (c Capability) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("capability name is required")
	}
	if c.Type == "" {
		return fmt.Errorf("capability type is required")
	}
	if c.Type == "ARGOCD" {
		if c.Configuration == nil {
			return fmt.Errorf("configuration is required for ARGOCD capability")
		}
		if c.Configuration.ArgoCD == nil {
			return fmt.Errorf("argocd configuration is required for ARGOCD capability")
		}
		if c.Configuration.ArgoCD.AWSIDC == nil {
			return fmt.Errorf("awsIdc configuration is required for ARGOCD capability")
		}
	}
	return nil
}
