package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_NielsenWatermarksSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.NielsenWatermarksSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-nielsenwatermarkssettings.html
type Channel_NielsenWatermarksSettings struct {

	// NielsenCbetSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-nielsenwatermarkssettings.html#cfn-medialive-channel-nielsenwatermarkssettings-nielsencbetsettings
	NielsenCbetSettings *Channel_NielsenCBET `json:"NielsenCbetSettings,omitempty"`

	// NielsenDistributionType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-nielsenwatermarkssettings.html#cfn-medialive-channel-nielsenwatermarkssettings-nielsendistributiontype
	NielsenDistributionType *types.Value `json:"NielsenDistributionType,omitempty"`

	// NielsenNaesIiNwSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-nielsenwatermarkssettings.html#cfn-medialive-channel-nielsenwatermarkssettings-nielsennaesiinwsettings
	NielsenNaesIiNwSettings *Channel_NielsenNaesIiNw `json:"NielsenNaesIiNwSettings,omitempty"`

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
func (r *Channel_NielsenWatermarksSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.NielsenWatermarksSettings"
}
