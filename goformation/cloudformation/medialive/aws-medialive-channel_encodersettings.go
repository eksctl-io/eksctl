package medialive

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Channel_EncoderSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.EncoderSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-encodersettings.html
type Channel_EncoderSettings struct {

	// AudioDescriptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-encodersettings.html#cfn-medialive-channel-encodersettings-audiodescriptions
	AudioDescriptions []Channel_AudioDescription `json:"AudioDescriptions,omitempty"`

	// AvailBlanking AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-encodersettings.html#cfn-medialive-channel-encodersettings-availblanking
	AvailBlanking *Channel_AvailBlanking `json:"AvailBlanking,omitempty"`

	// AvailConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-encodersettings.html#cfn-medialive-channel-encodersettings-availconfiguration
	AvailConfiguration *Channel_AvailConfiguration `json:"AvailConfiguration,omitempty"`

	// BlackoutSlate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-encodersettings.html#cfn-medialive-channel-encodersettings-blackoutslate
	BlackoutSlate *Channel_BlackoutSlate `json:"BlackoutSlate,omitempty"`

	// CaptionDescriptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-encodersettings.html#cfn-medialive-channel-encodersettings-captiondescriptions
	CaptionDescriptions []Channel_CaptionDescription `json:"CaptionDescriptions,omitempty"`

	// FeatureActivations AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-encodersettings.html#cfn-medialive-channel-encodersettings-featureactivations
	FeatureActivations *Channel_FeatureActivations `json:"FeatureActivations,omitempty"`

	// GlobalConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-encodersettings.html#cfn-medialive-channel-encodersettings-globalconfiguration
	GlobalConfiguration *Channel_GlobalConfiguration `json:"GlobalConfiguration,omitempty"`

	// MotionGraphicsConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-encodersettings.html#cfn-medialive-channel-encodersettings-motiongraphicsconfiguration
	MotionGraphicsConfiguration *Channel_MotionGraphicsConfiguration `json:"MotionGraphicsConfiguration,omitempty"`

	// NielsenConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-encodersettings.html#cfn-medialive-channel-encodersettings-nielsenconfiguration
	NielsenConfiguration *Channel_NielsenConfiguration `json:"NielsenConfiguration,omitempty"`

	// OutputGroups AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-encodersettings.html#cfn-medialive-channel-encodersettings-outputgroups
	OutputGroups []Channel_OutputGroup `json:"OutputGroups,omitempty"`

	// TimecodeConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-encodersettings.html#cfn-medialive-channel-encodersettings-timecodeconfig
	TimecodeConfig *Channel_TimecodeConfig `json:"TimecodeConfig,omitempty"`

	// VideoDescriptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-encodersettings.html#cfn-medialive-channel-encodersettings-videodescriptions
	VideoDescriptions []Channel_VideoDescription `json:"VideoDescriptions,omitempty"`

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
func (r *Channel_EncoderSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.EncoderSettings"
}
