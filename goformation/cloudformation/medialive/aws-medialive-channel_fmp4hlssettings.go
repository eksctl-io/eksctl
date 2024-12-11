package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_Fmp4HlsSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.Fmp4HlsSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-fmp4hlssettings.html
type Channel_Fmp4HlsSettings struct {

	// AudioRenditionSets AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-fmp4hlssettings.html#cfn-medialive-channel-fmp4hlssettings-audiorenditionsets
	AudioRenditionSets *types.Value `json:"AudioRenditionSets,omitempty"`

	// NielsenId3Behavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-fmp4hlssettings.html#cfn-medialive-channel-fmp4hlssettings-nielsenid3behavior
	NielsenId3Behavior *types.Value `json:"NielsenId3Behavior,omitempty"`

	// TimedMetadataBehavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-fmp4hlssettings.html#cfn-medialive-channel-fmp4hlssettings-timedmetadatabehavior
	TimedMetadataBehavior *types.Value `json:"TimedMetadataBehavior,omitempty"`

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
func (r *Channel_Fmp4HlsSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.Fmp4HlsSettings"
}
