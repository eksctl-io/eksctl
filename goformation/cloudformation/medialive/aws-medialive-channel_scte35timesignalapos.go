package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_Scte35TimeSignalApos AWS CloudFormation Resource (AWS::MediaLive::Channel.Scte35TimeSignalApos)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-scte35timesignalapos.html
type Channel_Scte35TimeSignalApos struct {

	// AdAvailOffset AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-scte35timesignalapos.html#cfn-medialive-channel-scte35timesignalapos-adavailoffset
	AdAvailOffset *types.Value `json:"AdAvailOffset,omitempty"`

	// NoRegionalBlackoutFlag AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-scte35timesignalapos.html#cfn-medialive-channel-scte35timesignalapos-noregionalblackoutflag
	NoRegionalBlackoutFlag *types.Value `json:"NoRegionalBlackoutFlag,omitempty"`

	// WebDeliveryAllowedFlag AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-scte35timesignalapos.html#cfn-medialive-channel-scte35timesignalapos-webdeliveryallowedflag
	WebDeliveryAllowedFlag *types.Value `json:"WebDeliveryAllowedFlag,omitempty"`

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
func (r *Channel_Scte35TimeSignalApos) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.Scte35TimeSignalApos"
}
