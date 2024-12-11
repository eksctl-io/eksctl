package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_InputLossBehavior AWS CloudFormation Resource (AWS::MediaLive::Channel.InputLossBehavior)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-inputlossbehavior.html
type Channel_InputLossBehavior struct {

	// BlackFrameMsec AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-inputlossbehavior.html#cfn-medialive-channel-inputlossbehavior-blackframemsec
	BlackFrameMsec *types.Value `json:"BlackFrameMsec,omitempty"`

	// InputLossImageColor AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-inputlossbehavior.html#cfn-medialive-channel-inputlossbehavior-inputlossimagecolor
	InputLossImageColor *types.Value `json:"InputLossImageColor,omitempty"`

	// InputLossImageSlate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-inputlossbehavior.html#cfn-medialive-channel-inputlossbehavior-inputlossimageslate
	InputLossImageSlate *Channel_InputLocation `json:"InputLossImageSlate,omitempty"`

	// InputLossImageType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-inputlossbehavior.html#cfn-medialive-channel-inputlossbehavior-inputlossimagetype
	InputLossImageType *types.Value `json:"InputLossImageType,omitempty"`

	// RepeatFrameMsec AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-inputlossbehavior.html#cfn-medialive-channel-inputlossbehavior-repeatframemsec
	RepeatFrameMsec *types.Value `json:"RepeatFrameMsec,omitempty"`

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
func (r *Channel_InputLossBehavior) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.InputLossBehavior"
}
