package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// VPNConnection_CloudwatchLogOptionsSpecification AWS CloudFormation Resource (AWS::EC2::VPNConnection.CloudwatchLogOptionsSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-vpnconnection-cloudwatchlogoptionsspecification.html
type VPNConnection_CloudwatchLogOptionsSpecification struct {

	// LogEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-vpnconnection-cloudwatchlogoptionsspecification.html#cfn-ec2-vpnconnection-cloudwatchlogoptionsspecification-logenabled
	LogEnabled *types.Value `json:"LogEnabled,omitempty"`

	// LogGroupArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-vpnconnection-cloudwatchlogoptionsspecification.html#cfn-ec2-vpnconnection-cloudwatchlogoptionsspecification-loggrouparn
	LogGroupArn *types.Value `json:"LogGroupArn,omitempty"`

	// LogOutputFormat AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-vpnconnection-cloudwatchlogoptionsspecification.html#cfn-ec2-vpnconnection-cloudwatchlogoptionsspecification-logoutputformat
	LogOutputFormat *types.Value `json:"LogOutputFormat,omitempty"`

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
func (r *VPNConnection_CloudwatchLogOptionsSpecification) AWSCloudFormationType() string {
	return "AWS::EC2::VPNConnection.CloudwatchLogOptionsSpecification"
}
