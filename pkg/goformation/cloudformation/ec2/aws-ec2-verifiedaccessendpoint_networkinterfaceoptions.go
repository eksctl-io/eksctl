package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// VerifiedAccessEndpoint_NetworkInterfaceOptions AWS CloudFormation Resource (AWS::EC2::VerifiedAccessEndpoint.NetworkInterfaceOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-networkinterfaceoptions.html
type VerifiedAccessEndpoint_NetworkInterfaceOptions struct {

	// NetworkInterfaceId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-networkinterfaceoptions.html#cfn-ec2-verifiedaccessendpoint-networkinterfaceoptions-networkinterfaceid
	NetworkInterfaceId *types.Value `json:"NetworkInterfaceId,omitempty"`

	// Port AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-networkinterfaceoptions.html#cfn-ec2-verifiedaccessendpoint-networkinterfaceoptions-port
	Port *types.Value `json:"Port,omitempty"`

	// PortRanges AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-networkinterfaceoptions.html#cfn-ec2-verifiedaccessendpoint-networkinterfaceoptions-portranges
	PortRanges []VerifiedAccessEndpoint_PortRange `json:"PortRanges,omitempty"`

	// Protocol AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessendpoint-networkinterfaceoptions.html#cfn-ec2-verifiedaccessendpoint-networkinterfaceoptions-protocol
	Protocol *types.Value `json:"Protocol,omitempty"`

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
func (r *VerifiedAccessEndpoint_NetworkInterfaceOptions) AWSCloudFormationType() string {
	return "AWS::EC2::VerifiedAccessEndpoint.NetworkInterfaceOptions"
}
