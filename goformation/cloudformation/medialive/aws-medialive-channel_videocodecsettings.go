package medialive

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Channel_VideoCodecSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.VideoCodecSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-videocodecsettings.html
type Channel_VideoCodecSettings struct {

	// FrameCaptureSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-videocodecsettings.html#cfn-medialive-channel-videocodecsettings-framecapturesettings
	FrameCaptureSettings *Channel_FrameCaptureSettings `json:"FrameCaptureSettings,omitempty"`

	// H264Settings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-videocodecsettings.html#cfn-medialive-channel-videocodecsettings-h264settings
	H264Settings *Channel_H264Settings `json:"H264Settings,omitempty"`

	// H265Settings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-videocodecsettings.html#cfn-medialive-channel-videocodecsettings-h265settings
	H265Settings *Channel_H265Settings `json:"H265Settings,omitempty"`

	// Mpeg2Settings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-videocodecsettings.html#cfn-medialive-channel-videocodecsettings-mpeg2settings
	Mpeg2Settings *Channel_Mpeg2Settings `json:"Mpeg2Settings,omitempty"`

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
func (r *Channel_VideoCodecSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.VideoCodecSettings"
}
