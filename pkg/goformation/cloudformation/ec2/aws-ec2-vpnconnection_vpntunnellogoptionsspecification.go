package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// VPNConnection_VpnTunnelLogOptionsSpecification AWS CloudFormation Resource (AWS::EC2::VPNConnection.VpnTunnelLogOptionsSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-vpnconnection-vpntunnellogoptionsspecification.html
type VPNConnection_VpnTunnelLogOptionsSpecification struct {

	// CloudwatchLogOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-vpnconnection-vpntunnellogoptionsspecification.html#cfn-ec2-vpnconnection-vpntunnellogoptionsspecification-cloudwatchlogoptions
	CloudwatchLogOptions *VPNConnection_CloudwatchLogOptionsSpecification `json:"CloudwatchLogOptions,omitempty"`

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
func (r *VPNConnection_VpnTunnelLogOptionsSpecification) AWSCloudFormationType() string {
	return "AWS::EC2::VPNConnection.VpnTunnelLogOptionsSpecification"
}
