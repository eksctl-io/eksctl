package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_AudioOnlyHlsSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.AudioOnlyHlsSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-audioonlyhlssettings.html
type Channel_AudioOnlyHlsSettings struct {

	// AudioGroupId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-audioonlyhlssettings.html#cfn-medialive-channel-audioonlyhlssettings-audiogroupid
	AudioGroupId *types.Value `json:"AudioGroupId,omitempty"`

	// AudioOnlyImage AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-audioonlyhlssettings.html#cfn-medialive-channel-audioonlyhlssettings-audioonlyimage
	AudioOnlyImage *Channel_InputLocation `json:"AudioOnlyImage,omitempty"`

	// AudioTrackType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-audioonlyhlssettings.html#cfn-medialive-channel-audioonlyhlssettings-audiotracktype
	AudioTrackType *types.Value `json:"AudioTrackType,omitempty"`

	// SegmentType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-audioonlyhlssettings.html#cfn-medialive-channel-audioonlyhlssettings-segmenttype
	SegmentType *types.Value `json:"SegmentType,omitempty"`

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
func (r *Channel_AudioOnlyHlsSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.AudioOnlyHlsSettings"
}
