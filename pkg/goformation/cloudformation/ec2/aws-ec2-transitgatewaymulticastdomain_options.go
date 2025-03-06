package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// TransitGatewayMulticastDomain_Options AWS CloudFormation Resource (AWS::EC2::TransitGatewayMulticastDomain.Options)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-transitgatewaymulticastdomain-options.html
type TransitGatewayMulticastDomain_Options struct {

	// AutoAcceptSharedAssociations AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-transitgatewaymulticastdomain-options.html#cfn-ec2-transitgatewaymulticastdomain-options-autoacceptsharedassociations
	AutoAcceptSharedAssociations *types.Value `json:"AutoAcceptSharedAssociations,omitempty"`

	// Igmpv2Support AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-transitgatewaymulticastdomain-options.html#cfn-ec2-transitgatewaymulticastdomain-options-igmpv2support
	Igmpv2Support *types.Value `json:"Igmpv2Support,omitempty"`

	// StaticSourcesSupport AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-transitgatewaymulticastdomain-options.html#cfn-ec2-transitgatewaymulticastdomain-options-staticsourcessupport
	StaticSourcesSupport *types.Value `json:"StaticSourcesSupport,omitempty"`

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
func (r *TransitGatewayMulticastDomain_Options) AWSCloudFormationType() string {
	return "AWS::EC2::TransitGatewayMulticastDomain.Options"
}
