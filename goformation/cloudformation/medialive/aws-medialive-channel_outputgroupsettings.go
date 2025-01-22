package medialive

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Channel_OutputGroupSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.OutputGroupSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputgroupsettings.html
type Channel_OutputGroupSettings struct {

	// ArchiveGroupSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputgroupsettings.html#cfn-medialive-channel-outputgroupsettings-archivegroupsettings
	ArchiveGroupSettings *Channel_ArchiveGroupSettings `json:"ArchiveGroupSettings,omitempty"`

	// FrameCaptureGroupSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputgroupsettings.html#cfn-medialive-channel-outputgroupsettings-framecapturegroupsettings
	FrameCaptureGroupSettings *Channel_FrameCaptureGroupSettings `json:"FrameCaptureGroupSettings,omitempty"`

	// HlsGroupSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputgroupsettings.html#cfn-medialive-channel-outputgroupsettings-hlsgroupsettings
	HlsGroupSettings *Channel_HlsGroupSettings `json:"HlsGroupSettings,omitempty"`

	// MediaPackageGroupSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputgroupsettings.html#cfn-medialive-channel-outputgroupsettings-mediapackagegroupsettings
	MediaPackageGroupSettings *Channel_MediaPackageGroupSettings `json:"MediaPackageGroupSettings,omitempty"`

	// MsSmoothGroupSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputgroupsettings.html#cfn-medialive-channel-outputgroupsettings-mssmoothgroupsettings
	MsSmoothGroupSettings *Channel_MsSmoothGroupSettings `json:"MsSmoothGroupSettings,omitempty"`

	// MultiplexGroupSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputgroupsettings.html#cfn-medialive-channel-outputgroupsettings-multiplexgroupsettings
	MultiplexGroupSettings *Channel_MultiplexGroupSettings `json:"MultiplexGroupSettings,omitempty"`

	// RtmpGroupSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputgroupsettings.html#cfn-medialive-channel-outputgroupsettings-rtmpgroupsettings
	RtmpGroupSettings *Channel_RtmpGroupSettings `json:"RtmpGroupSettings,omitempty"`

	// UdpGroupSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputgroupsettings.html#cfn-medialive-channel-outputgroupsettings-udpgroupsettings
	UdpGroupSettings *Channel_UdpGroupSettings `json:"UdpGroupSettings,omitempty"`

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
func (r *Channel_OutputGroupSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.OutputGroupSettings"
}
