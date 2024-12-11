package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_TemporalFilterSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.TemporalFilterSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-temporalfiltersettings.html
type Channel_TemporalFilterSettings struct {

	// PostFilterSharpening AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-temporalfiltersettings.html#cfn-medialive-channel-temporalfiltersettings-postfiltersharpening
	PostFilterSharpening *types.Value `json:"PostFilterSharpening,omitempty"`

	// Strength AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-temporalfiltersettings.html#cfn-medialive-channel-temporalfiltersettings-strength
	Strength *types.Value `json:"Strength,omitempty"`

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
func (r *Channel_TemporalFilterSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.TemporalFilterSettings"
}
