package medialive

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Channel_OutputSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.OutputSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputsettings.html
type Channel_OutputSettings struct {

	// ArchiveOutputSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputsettings.html#cfn-medialive-channel-outputsettings-archiveoutputsettings
	ArchiveOutputSettings *Channel_ArchiveOutputSettings `json:"ArchiveOutputSettings,omitempty"`

	// FrameCaptureOutputSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputsettings.html#cfn-medialive-channel-outputsettings-framecaptureoutputsettings
	FrameCaptureOutputSettings *Channel_FrameCaptureOutputSettings `json:"FrameCaptureOutputSettings,omitempty"`

	// HlsOutputSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputsettings.html#cfn-medialive-channel-outputsettings-hlsoutputsettings
	HlsOutputSettings *Channel_HlsOutputSettings `json:"HlsOutputSettings,omitempty"`

	// MediaPackageOutputSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputsettings.html#cfn-medialive-channel-outputsettings-mediapackageoutputsettings
	MediaPackageOutputSettings *Channel_MediaPackageOutputSettings `json:"MediaPackageOutputSettings,omitempty"`

	// MsSmoothOutputSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputsettings.html#cfn-medialive-channel-outputsettings-mssmoothoutputsettings
	MsSmoothOutputSettings *Channel_MsSmoothOutputSettings `json:"MsSmoothOutputSettings,omitempty"`

	// MultiplexOutputSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputsettings.html#cfn-medialive-channel-outputsettings-multiplexoutputsettings
	MultiplexOutputSettings *Channel_MultiplexOutputSettings `json:"MultiplexOutputSettings,omitempty"`

	// RtmpOutputSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputsettings.html#cfn-medialive-channel-outputsettings-rtmpoutputsettings
	RtmpOutputSettings *Channel_RtmpOutputSettings `json:"RtmpOutputSettings,omitempty"`

	// UdpOutputSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputsettings.html#cfn-medialive-channel-outputsettings-udpoutputsettings
	UdpOutputSettings *Channel_UdpOutputSettings `json:"UdpOutputSettings,omitempty"`

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
func (r *Channel_OutputSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.OutputSettings"
}
