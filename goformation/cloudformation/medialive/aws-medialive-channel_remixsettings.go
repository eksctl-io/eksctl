package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_RemixSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.RemixSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-remixsettings.html
type Channel_RemixSettings struct {

	// ChannelMappings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-remixsettings.html#cfn-medialive-channel-remixsettings-channelmappings
	ChannelMappings []Channel_AudioChannelMapping `json:"ChannelMappings,omitempty"`

	// ChannelsIn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-remixsettings.html#cfn-medialive-channel-remixsettings-channelsin
	ChannelsIn *types.Value `json:"ChannelsIn,omitempty"`

	// ChannelsOut AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-remixsettings.html#cfn-medialive-channel-remixsettings-channelsout
	ChannelsOut *types.Value `json:"ChannelsOut,omitempty"`

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
func (r *Channel_RemixSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.RemixSettings"
}
