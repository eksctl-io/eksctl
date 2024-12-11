package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_RtmpGroupSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.RtmpGroupSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-rtmpgroupsettings.html
type Channel_RtmpGroupSettings struct {

	// AdMarkers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-rtmpgroupsettings.html#cfn-medialive-channel-rtmpgroupsettings-admarkers
	AdMarkers *types.Value `json:"AdMarkers,omitempty"`

	// AuthenticationScheme AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-rtmpgroupsettings.html#cfn-medialive-channel-rtmpgroupsettings-authenticationscheme
	AuthenticationScheme *types.Value `json:"AuthenticationScheme,omitempty"`

	// CacheFullBehavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-rtmpgroupsettings.html#cfn-medialive-channel-rtmpgroupsettings-cachefullbehavior
	CacheFullBehavior *types.Value `json:"CacheFullBehavior,omitempty"`

	// CacheLength AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-rtmpgroupsettings.html#cfn-medialive-channel-rtmpgroupsettings-cachelength
	CacheLength *types.Value `json:"CacheLength,omitempty"`

	// CaptionData AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-rtmpgroupsettings.html#cfn-medialive-channel-rtmpgroupsettings-captiondata
	CaptionData *types.Value `json:"CaptionData,omitempty"`

	// InputLossAction AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-rtmpgroupsettings.html#cfn-medialive-channel-rtmpgroupsettings-inputlossaction
	InputLossAction *types.Value `json:"InputLossAction,omitempty"`

	// RestartDelay AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-rtmpgroupsettings.html#cfn-medialive-channel-rtmpgroupsettings-restartdelay
	RestartDelay *types.Value `json:"RestartDelay,omitempty"`

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
func (r *Channel_RtmpGroupSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.RtmpGroupSettings"
}
