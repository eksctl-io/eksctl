package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_AudioChannelMapping AWS CloudFormation Resource (AWS::MediaLive::Channel.AudioChannelMapping)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-audiochannelmapping.html
type Channel_AudioChannelMapping struct {

	// InputChannelLevels AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-audiochannelmapping.html#cfn-medialive-channel-audiochannelmapping-inputchannellevels
	InputChannelLevels []Channel_InputChannelLevel `json:"InputChannelLevels,omitempty"`

	// OutputChannel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-audiochannelmapping.html#cfn-medialive-channel-audiochannelmapping-outputchannel
	OutputChannel *types.Value `json:"OutputChannel,omitempty"`

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
func (r *Channel_AudioChannelMapping) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.AudioChannelMapping"
}
