package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_AudioNormalizationSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.AudioNormalizationSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-audionormalizationsettings.html
type Channel_AudioNormalizationSettings struct {

	// Algorithm AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-audionormalizationsettings.html#cfn-medialive-channel-audionormalizationsettings-algorithm
	Algorithm *types.Value `json:"Algorithm,omitempty"`

	// AlgorithmControl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-audionormalizationsettings.html#cfn-medialive-channel-audionormalizationsettings-algorithmcontrol
	AlgorithmControl *types.Value `json:"AlgorithmControl,omitempty"`

	// TargetLkfs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-audionormalizationsettings.html#cfn-medialive-channel-audionormalizationsettings-targetlkfs
	TargetLkfs *types.Value `json:"TargetLkfs,omitempty"`

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
func (r *Channel_AudioNormalizationSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.AudioNormalizationSettings"
}
