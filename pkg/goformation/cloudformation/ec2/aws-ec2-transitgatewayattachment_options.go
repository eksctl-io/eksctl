package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// TransitGatewayAttachment_Options AWS CloudFormation Resource (AWS::EC2::TransitGatewayAttachment.Options)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-transitgatewayattachment-options.html
type TransitGatewayAttachment_Options struct {

	// ApplianceModeSupport AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-transitgatewayattachment-options.html#cfn-ec2-transitgatewayattachment-options-appliancemodesupport
	ApplianceModeSupport *types.Value `json:"ApplianceModeSupport,omitempty"`

	// DnsSupport AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-transitgatewayattachment-options.html#cfn-ec2-transitgatewayattachment-options-dnssupport
	DnsSupport *types.Value `json:"DnsSupport,omitempty"`

	// Ipv6Support AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-transitgatewayattachment-options.html#cfn-ec2-transitgatewayattachment-options-ipv6support
	Ipv6Support *types.Value `json:"Ipv6Support,omitempty"`

	// SecurityGroupReferencingSupport AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-transitgatewayattachment-options.html#cfn-ec2-transitgatewayattachment-options-securitygroupreferencingsupport
	SecurityGroupReferencingSupport *types.Value `json:"SecurityGroupReferencingSupport,omitempty"`

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
func (r *TransitGatewayAttachment_Options) AWSCloudFormationType() string {
	return "AWS::EC2::TransitGatewayAttachment.Options"
}
