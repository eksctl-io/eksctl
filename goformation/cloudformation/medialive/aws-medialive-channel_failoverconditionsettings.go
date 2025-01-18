package medialive

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Channel_FailoverConditionSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.FailoverConditionSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-failoverconditionsettings.html
type Channel_FailoverConditionSettings struct {

	// AudioSilenceSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-failoverconditionsettings.html#cfn-medialive-channel-failoverconditionsettings-audiosilencesettings
	AudioSilenceSettings *Channel_AudioSilenceFailoverSettings `json:"AudioSilenceSettings,omitempty"`

	// InputLossSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-failoverconditionsettings.html#cfn-medialive-channel-failoverconditionsettings-inputlosssettings
	InputLossSettings *Channel_InputLossFailoverSettings `json:"InputLossSettings,omitempty"`

	// VideoBlackSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-failoverconditionsettings.html#cfn-medialive-channel-failoverconditionsettings-videoblacksettings
	VideoBlackSettings *Channel_VideoBlackFailoverSettings `json:"VideoBlackSettings,omitempty"`

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
func (r *Channel_FailoverConditionSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.FailoverConditionSettings"
}
