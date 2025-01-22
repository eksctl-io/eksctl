package medialive

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Channel_CaptionDestinationSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.CaptionDestinationSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-captiondestinationsettings.html
type Channel_CaptionDestinationSettings struct {

	// AribDestinationSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-captiondestinationsettings.html#cfn-medialive-channel-captiondestinationsettings-aribdestinationsettings
	AribDestinationSettings *Channel_AribDestinationSettings `json:"AribDestinationSettings,omitempty"`

	// BurnInDestinationSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-captiondestinationsettings.html#cfn-medialive-channel-captiondestinationsettings-burnindestinationsettings
	BurnInDestinationSettings *Channel_BurnInDestinationSettings `json:"BurnInDestinationSettings,omitempty"`

	// DvbSubDestinationSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-captiondestinationsettings.html#cfn-medialive-channel-captiondestinationsettings-dvbsubdestinationsettings
	DvbSubDestinationSettings *Channel_DvbSubDestinationSettings `json:"DvbSubDestinationSettings,omitempty"`

	// EbuTtDDestinationSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-captiondestinationsettings.html#cfn-medialive-channel-captiondestinationsettings-ebuttddestinationsettings
	EbuTtDDestinationSettings *Channel_EbuTtDDestinationSettings `json:"EbuTtDDestinationSettings,omitempty"`

	// EmbeddedDestinationSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-captiondestinationsettings.html#cfn-medialive-channel-captiondestinationsettings-embeddeddestinationsettings
	EmbeddedDestinationSettings *Channel_EmbeddedDestinationSettings `json:"EmbeddedDestinationSettings,omitempty"`

	// EmbeddedPlusScte20DestinationSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-captiondestinationsettings.html#cfn-medialive-channel-captiondestinationsettings-embeddedplusscte20destinationsettings
	EmbeddedPlusScte20DestinationSettings *Channel_EmbeddedPlusScte20DestinationSettings `json:"EmbeddedPlusScte20DestinationSettings,omitempty"`

	// RtmpCaptionInfoDestinationSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-captiondestinationsettings.html#cfn-medialive-channel-captiondestinationsettings-rtmpcaptioninfodestinationsettings
	RtmpCaptionInfoDestinationSettings *Channel_RtmpCaptionInfoDestinationSettings `json:"RtmpCaptionInfoDestinationSettings,omitempty"`

	// Scte20PlusEmbeddedDestinationSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-captiondestinationsettings.html#cfn-medialive-channel-captiondestinationsettings-scte20plusembeddeddestinationsettings
	Scte20PlusEmbeddedDestinationSettings *Channel_Scte20PlusEmbeddedDestinationSettings `json:"Scte20PlusEmbeddedDestinationSettings,omitempty"`

	// Scte27DestinationSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-captiondestinationsettings.html#cfn-medialive-channel-captiondestinationsettings-scte27destinationsettings
	Scte27DestinationSettings *Channel_Scte27DestinationSettings `json:"Scte27DestinationSettings,omitempty"`

	// SmpteTtDestinationSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-captiondestinationsettings.html#cfn-medialive-channel-captiondestinationsettings-smptettdestinationsettings
	SmpteTtDestinationSettings *Channel_SmpteTtDestinationSettings `json:"SmpteTtDestinationSettings,omitempty"`

	// TeletextDestinationSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-captiondestinationsettings.html#cfn-medialive-channel-captiondestinationsettings-teletextdestinationsettings
	TeletextDestinationSettings *Channel_TeletextDestinationSettings `json:"TeletextDestinationSettings,omitempty"`

	// TtmlDestinationSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-captiondestinationsettings.html#cfn-medialive-channel-captiondestinationsettings-ttmldestinationsettings
	TtmlDestinationSettings *Channel_TtmlDestinationSettings `json:"TtmlDestinationSettings,omitempty"`

	// WebvttDestinationSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-captiondestinationsettings.html#cfn-medialive-channel-captiondestinationsettings-webvttdestinationsettings
	WebvttDestinationSettings *Channel_WebvttDestinationSettings `json:"WebvttDestinationSettings,omitempty"`

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
func (r *Channel_CaptionDestinationSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.CaptionDestinationSettings"
}
