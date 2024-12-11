package ec2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// IPAM_IpamOperatingRegion AWS CloudFormation Resource (AWS::EC2::IPAM.IpamOperatingRegion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ipam-ipamoperatingregion.html
type IPAM_IpamOperatingRegion struct {

	// RegionName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ipam-ipamoperatingregion.html#cfn-ec2-ipam-ipamoperatingregion-regionname
	RegionName *types.Value `json:"RegionName,omitempty"`

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
func (r *IPAM_IpamOperatingRegion) AWSCloudFormationType() string {
	return "AWS::EC2::IPAM.IpamOperatingRegion"
}
