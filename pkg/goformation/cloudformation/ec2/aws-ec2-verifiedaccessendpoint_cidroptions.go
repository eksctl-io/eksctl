package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// VerifiedAccessEndpoint_CidrOptions AWS CloudFormation Resource (AWS::EC2::VerifiedAccessEndpoint.CidrOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-cidroptions.html
type VerifiedAccessEndpoint_CidrOptions struct {

	// Cidr AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-cidroptions.html#cfn-ec2-verifiedaccessendpoint-cidroptions-cidr
	Cidr *types.Value `json:"Cidr,omitempty"`

	// PortRanges AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-cidroptions.html#cfn-ec2-verifiedaccessendpoint-cidroptions-portranges
	PortRanges []VerifiedAccessEndpoint_PortRange `json:"PortRanges,omitempty"`

	// Protocol AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-cidroptions.html#cfn-ec2-verifiedaccessendpoint-cidroptions-protocol
	Protocol *types.Value `json:"Protocol,omitempty"`

	// SubnetIds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-cidroptions.html#cfn-ec2-verifiedaccessendpoint-cidroptions-subnetids
	SubnetIds *types.Value `json:"SubnetIds,omitempty"`

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
func (r *VerifiedAccessEndpoint_CidrOptions) AWSCloudFormationType() string {
	return "AWS::EC2::VerifiedAccessEndpoint.CidrOptions"
}
