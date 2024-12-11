package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_DvbSubDestinationSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.DvbSubDestinationSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-dvbsubdestinationsettings.html
type Channel_DvbSubDestinationSettings struct {

	// Alignment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-dvbsubdestinationsettings.html#cfn-medialive-channel-dvbsubdestinationsettings-alignment
	Alignment *types.Value `json:"Alignment,omitempty"`

	// BackgroundColor AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-dvbsubdestinationsettings.html#cfn-medialive-channel-dvbsubdestinationsettings-backgroundcolor
	BackgroundColor *types.Value `json:"BackgroundColor,omitempty"`

	// BackgroundOpacity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-dvbsubdestinationsettings.html#cfn-medialive-channel-dvbsubdestinationsettings-backgroundopacity
	BackgroundOpacity *types.Value `json:"BackgroundOpacity,omitempty"`

	// Font AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-dvbsubdestinationsettings.html#cfn-medialive-channel-dvbsubdestinationsettings-font
	Font *Channel_InputLocation `json:"Font,omitempty"`

	// FontColor AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-dvbsubdestinationsettings.html#cfn-medialive-channel-dvbsubdestinationsettings-fontcolor
	FontColor *types.Value `json:"FontColor,omitempty"`

	// FontOpacity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-dvbsubdestinationsettings.html#cfn-medialive-channel-dvbsubdestinationsettings-fontopacity
	FontOpacity *types.Value `json:"FontOpacity,omitempty"`

	// FontResolution AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-dvbsubdestinationsettings.html#cfn-medialive-channel-dvbsubdestinationsettings-fontresolution
	FontResolution *types.Value `json:"FontResolution,omitempty"`

	// FontSize AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-dvbsubdestinationsettings.html#cfn-medialive-channel-dvbsubdestinationsettings-fontsize
	FontSize *types.Value `json:"FontSize,omitempty"`

	// OutlineColor AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-dvbsubdestinationsettings.html#cfn-medialive-channel-dvbsubdestinationsettings-outlinecolor
	OutlineColor *types.Value `json:"OutlineColor,omitempty"`

	// OutlineSize AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-dvbsubdestinationsettings.html#cfn-medialive-channel-dvbsubdestinationsettings-outlinesize
	OutlineSize *types.Value `json:"OutlineSize,omitempty"`

	// ShadowColor AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-dvbsubdestinationsettings.html#cfn-medialive-channel-dvbsubdestinationsettings-shadowcolor
	ShadowColor *types.Value `json:"ShadowColor,omitempty"`

	// ShadowOpacity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-dvbsubdestinationsettings.html#cfn-medialive-channel-dvbsubdestinationsettings-shadowopacity
	ShadowOpacity *types.Value `json:"ShadowOpacity,omitempty"`

	// ShadowXOffset AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-dvbsubdestinationsettings.html#cfn-medialive-channel-dvbsubdestinationsettings-shadowxoffset
	ShadowXOffset *types.Value `json:"ShadowXOffset,omitempty"`

	// ShadowYOffset AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-dvbsubdestinationsettings.html#cfn-medialive-channel-dvbsubdestinationsettings-shadowyoffset
	ShadowYOffset *types.Value `json:"ShadowYOffset,omitempty"`

	// TeletextGridControl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-dvbsubdestinationsettings.html#cfn-medialive-channel-dvbsubdestinationsettings-teletextgridcontrol
	TeletextGridControl *types.Value `json:"TeletextGridControl,omitempty"`

	// XPosition AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-dvbsubdestinationsettings.html#cfn-medialive-channel-dvbsubdestinationsettings-xposition
	XPosition *types.Value `json:"XPosition,omitempty"`

	// YPosition AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-dvbsubdestinationsettings.html#cfn-medialive-channel-dvbsubdestinationsettings-yposition
	YPosition *types.Value `json:"YPosition,omitempty"`

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
func (r *Channel_DvbSubDestinationSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.DvbSubDestinationSettings"
}
