package medialive

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Channel_H265FilterSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.H265FilterSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-h265filtersettings.html
type Channel_H265FilterSettings struct {

	// TemporalFilterSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-h265filtersettings.html#cfn-medialive-channel-h265filtersettings-temporalfiltersettings
	TemporalFilterSettings *Channel_TemporalFilterSettings `json:"TemporalFilterSettings,omitempty"`

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
func (r *Channel_H265FilterSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.H265FilterSettings"
}
