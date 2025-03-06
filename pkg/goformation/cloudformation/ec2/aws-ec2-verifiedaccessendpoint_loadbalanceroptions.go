package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// VerifiedAccessEndpoint_LoadBalancerOptions AWS CloudFormation Resource (AWS::EC2::VerifiedAccessEndpoint.LoadBalancerOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-loadbalanceroptions.html
type VerifiedAccessEndpoint_LoadBalancerOptions struct {

	// LoadBalancerArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-loadbalanceroptions.html#cfn-ec2-verifiedaccessendpoint-loadbalanceroptions-loadbalancerarn
	LoadBalancerArn *types.Value `json:"LoadBalancerArn,omitempty"`

	// Port AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-loadbalanceroptions.html#cfn-ec2-verifiedaccessendpoint-loadbalanceroptions-port
	Port *types.Value `json:"Port,omitempty"`

	// PortRanges AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-loadbalanceroptions.html#cfn-ec2-verifiedaccessendpoint-loadbalanceroptions-portranges
	PortRanges []VerifiedAccessEndpoint_PortRange `json:"PortRanges,omitempty"`

	// Protocol AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-loadbalanceroptions.html#cfn-ec2-verifiedaccessendpoint-loadbalanceroptions-protocol
	Protocol *types.Value `json:"Protocol,omitempty"`

	// SubnetIds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-loadbalanceroptions.html#cfn-ec2-verifiedaccessendpoint-loadbalanceroptions-subnetids
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
func (r *VerifiedAccessEndpoint_LoadBalancerOptions) AWSCloudFormationType() string {
	return "AWS::EC2::VerifiedAccessEndpoint.LoadBalancerOptions"
}
