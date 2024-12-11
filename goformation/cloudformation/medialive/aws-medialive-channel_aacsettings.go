package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_AacSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.AacSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-aacsettings.html
type Channel_AacSettings struct {

	// Bitrate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-aacsettings.html#cfn-medialive-channel-aacsettings-bitrate
	Bitrate *types.Value `json:"Bitrate,omitempty"`

	// CodingMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-aacsettings.html#cfn-medialive-channel-aacsettings-codingmode
	CodingMode *types.Value `json:"CodingMode,omitempty"`

	// InputType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-aacsettings.html#cfn-medialive-channel-aacsettings-inputtype
	InputType *types.Value `json:"InputType,omitempty"`

	// Profile AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-aacsettings.html#cfn-medialive-channel-aacsettings-profile
	Profile *types.Value `json:"Profile,omitempty"`

	// RateControlMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-aacsettings.html#cfn-medialive-channel-aacsettings-ratecontrolmode
	RateControlMode *types.Value `json:"RateControlMode,omitempty"`

	// RawFormat AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-aacsettings.html#cfn-medialive-channel-aacsettings-rawformat
	RawFormat *types.Value `json:"RawFormat,omitempty"`

	// SampleRate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-aacsettings.html#cfn-medialive-channel-aacsettings-samplerate
	SampleRate *types.Value `json:"SampleRate,omitempty"`

	// Spec AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-aacsettings.html#cfn-medialive-channel-aacsettings-spec
	Spec *types.Value `json:"Spec,omitempty"`

	// VbrQuality AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-aacsettings.html#cfn-medialive-channel-aacsettings-vbrquality
	VbrQuality *types.Value `json:"VbrQuality,omitempty"`

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
func (r *Channel_AacSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.AacSettings"
}
