package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// IPAM_IpamOrganizationalUnitExclusion AWS CloudFormation Resource (AWS::EC2::IPAM.IpamOrganizationalUnitExclusion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ipam-ipamorganizationalunitexclusion.html
type IPAM_IpamOrganizationalUnitExclusion struct {

	// OrganizationsEntityPath AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ipam-ipamorganizationalunitexclusion.html#cfn-ec2-ipam-ipamorganizationalunitexclusion-organizationsentitypath
	OrganizationsEntityPath *types.Value `json:"OrganizationsEntityPath,omitempty"`

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
func (r *IPAM_IpamOrganizationalUnitExclusion) AWSCloudFormationType() string {
	return "AWS::EC2::IPAM.IpamOrganizationalUnitExclusion"
}
