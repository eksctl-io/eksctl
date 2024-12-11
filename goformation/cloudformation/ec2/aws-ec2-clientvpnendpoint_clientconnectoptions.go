package ec2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ClientVpnEndpoint_ClientConnectOptions AWS CloudFormation Resource (AWS::EC2::ClientVpnEndpoint.ClientConnectOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-clientvpnendpoint-clientconnectoptions.html
type ClientVpnEndpoint_ClientConnectOptions struct {

	// Enabled AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-clientvpnendpoint-clientconnectoptions.html#cfn-ec2-clientvpnendpoint-clientconnectoptions-enabled
	Enabled *types.Value `json:"Enabled"`

	// LambdaFunctionArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-clientvpnendpoint-clientconnectoptions.html#cfn-ec2-clientvpnendpoint-clientconnectoptions-lambdafunctionarn
	LambdaFunctionArn *types.Value `json:"LambdaFunctionArn,omitempty"`

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
func (r *ClientVpnEndpoint_ClientConnectOptions) AWSCloudFormationType() string {
	return "AWS::EC2::ClientVpnEndpoint.ClientConnectOptions"
}
