package eks

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/cloudformation"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// CapabilityConfiguration holds capability-specific configuration for CloudFormation
type CapabilityConfiguration struct {
	// ArgoCD configuration for ARGOCD capability type
	ArgoCd ArgoCDConfiguration `json:"ArgoCd,omitempty"`
}

// ArgoCDConfiguration holds ArgoCD-specific configuration for CloudFormation
type ArgoCDConfiguration struct {
	// Namespace for ArgoCD installation
	// +optional
	Namespace *types.Value `json:"Namespace,omitempty"`

	// NetworkAccess configuration
	// +optional
	NetworkAccess *ArgoCDNetworkAccess `json:"NetworkAccess,omitempty"`

	// RBACRoleMappings for ArgoCD RBAC
	// +optional
	RBACRoleMappings []ArgoCDRoleMapping `json:"RbacRoleMappings,omitempty"`

	// AWSIDC configuration
	// +optional
	AWSIDC *ArgoCDAWSIDC `json:"AwsIdc,omitempty"`
}

// ArgoCDNetworkAccess holds network access configuration for ArgoCD
type ArgoCDNetworkAccess struct {
	// VPCEIDs for VPC endpoint access
	// +optional
	VPCEIDs *types.Value `json:"VpceIds,omitempty"`
}

// ArgoCDRoleMapping holds RBAC role mapping for ArgoCD
type ArgoCDRoleMapping struct {
	// Role is the ArgoCD role (ADMIN, EDITOR, VIEWER)
	// +required
	Role *types.Value `json:"Role"`

	// Identities are the SSO identities to map to the role
	// +required
	Identities []SSOIdentity `json:"Identities"`
}

// SSOIdentity represents an SSO identity
type SSOIdentity struct {
	// ID of the SSO identity
	// +required
	ID *types.Value `json:"Id"`

	// Type of the SSO identity (SSO_USER, SSO_GROUP)
	// +required
	Type *types.Value `json:"Type"`
}

// ArgoCDAWSIDC holds AWS IDC configuration for ArgoCD
type ArgoCDAWSIDC struct {
	// IDCInstanceARN is the ARN of the IDC instance
	// +required
	IDCInstanceARN *types.Value `json:"IdcInstanceArn"`

	// IDCRegion is the region of the IDC instance
	// +optional
	IDCRegion *types.Value `json:"IdcRegion,omitempty"`
}

// Capability represents an AWS::EKS::Capability resource
type Capability struct {
	// CapabilityName AWS CloudFormation Property
	// Required: true
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-eks-capability.html#cfn-eks-capability-capabilityname
	CapabilityName *types.Value `json:"CapabilityName,omitempty"`

	// ClusterName AWS CloudFormation Property
	// Required: true
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-eks-capability.html#cfn-eks-capability-clustername
	ClusterName *types.Value `json:"ClusterName,omitempty"`

	// Configuration AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-eks-capability.html#cfn-eks-capability-configuration
	Configuration *CapabilityConfiguration `json:"Configuration,omitempty"`

	// DeletePropagationPolicy AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-eks-capability.html#cfn-eks-capability-deletepropagationpolicy
	DeletePropagationPolicy *types.Value `json:"DeletePropagationPolicy,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-eks-capability.html#cfn-eks-capability-rolearn
	RoleArn *types.Value `json:"RoleArn,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-eks-capability.html#cfn-eks-capability-tags
	Tags []cloudformation.Tag `json:"Tags,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-eks-capability.html#cfn-eks-capability-type
	Type *types.Value `json:"Type,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *Capability) AWSCloudFormationType() string {
	return "AWS::EKS::Capability"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r Capability) MarshalJSON() ([]byte, error) {
	type Properties Capability
	return json.Marshal(&struct {
		Type                string
		Properties          Properties
		DependsOn           []string                     `json:"DependsOn,omitempty"`
		Metadata            map[string]interface{}       `json:"Metadata,omitempty"`
		DeletionPolicy      policies.DeletionPolicy      `json:"DeletionPolicy,omitempty"`
		UpdateReplacePolicy policies.UpdateReplacePolicy `json:"UpdateReplacePolicy,omitempty"`
		Condition           string                       `json:"Condition,omitempty"`
	}{
		Type:                r.AWSCloudFormationType(),
		Properties:          (Properties)(r),
		DependsOn:           r.AWSCloudFormationDependsOn,
		Metadata:            r.AWSCloudFormationMetadata,
		DeletionPolicy:      r.AWSCloudFormationDeletionPolicy,
		UpdateReplacePolicy: r.AWSCloudFormationUpdateReplacePolicy,
		Condition:           r.AWSCloudFormationCondition,
	})
}

// UnmarshalJSON is a custom JSON unmarshalling hook that strips the outer
// AWS CloudFormation resource object, and just keeps the 'Properties' field.
func (r *Capability) UnmarshalJSON(b []byte) error {
	type Properties Capability
	res := &struct {
		Type                string
		Properties          *Properties
		DependsOn           []string
		Metadata            map[string]interface{}
		DeletionPolicy      string
		UpdateReplacePolicy string
		Condition           string
	}{}

	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&res); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return err
	}

	// If the resource has no Properties set, it could be nil
	if res.Properties != nil {
		*r = Capability(*res.Properties)
	}
	if res.DependsOn != nil {
		r.AWSCloudFormationDependsOn = res.DependsOn
	}
	if res.Metadata != nil {
		r.AWSCloudFormationMetadata = res.Metadata
	}
	if res.DeletionPolicy != "" {
		r.AWSCloudFormationDeletionPolicy = policies.DeletionPolicy(res.DeletionPolicy)
	}
	if res.UpdateReplacePolicy != "" {
		r.AWSCloudFormationUpdateReplacePolicy = policies.UpdateReplacePolicy(res.UpdateReplacePolicy)
	}
	if res.Condition != "" {
		r.AWSCloudFormationCondition = res.Condition
	}
	return nil
}
