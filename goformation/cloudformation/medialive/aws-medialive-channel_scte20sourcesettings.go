package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_Scte20SourceSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.Scte20SourceSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-scte20sourcesettings.html
type Channel_Scte20SourceSettings struct {

	// Convert608To708 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-scte20sourcesettings.html#cfn-medialive-channel-scte20sourcesettings-convert608to708
	Convert608To708 *types.Value `json:"Convert608To708,omitempty"`

	// Source608ChannelNumber AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-scte20sourcesettings.html#cfn-medialive-channel-scte20sourcesettings-source608channelnumber
	Source608ChannelNumber *types.Value `json:"Source608ChannelNumber,omitempty"`

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
func (r *Channel_Scte20SourceSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.Scte20SourceSettings"
}
