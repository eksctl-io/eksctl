package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// VerifiedAccessEndpoint_PortRange AWS CloudFormation Resource (AWS::EC2::VerifiedAccessEndpoint.PortRange)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-portrange.html
type VerifiedAccessEndpoint_PortRange struct {

	// FromPort AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-portrange.html#cfn-ec2-verifiedaccessendpoint-portrange-fromport
	FromPort *types.Value `json:"FromPort,omitempty"`

	// ToPort AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-portrange.html#cfn-ec2-verifiedaccessendpoint-portrange-toport
	ToPort *types.Value `json:"ToPort,omitempty"`

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
func (r *VerifiedAccessEndpoint_PortRange) AWSCloudFormationType() string {
	return "AWS::EC2::VerifiedAccessEndpoint.PortRange"
}
