package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_AutomaticInputFailoverSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.AutomaticInputFailoverSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-automaticinputfailoversettings.html
type Channel_AutomaticInputFailoverSettings struct {

	// ErrorClearTimeMsec AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-automaticinputfailoversettings.html#cfn-medialive-channel-automaticinputfailoversettings-errorcleartimemsec
	ErrorClearTimeMsec *types.Value `json:"ErrorClearTimeMsec,omitempty"`

	// FailoverConditions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-automaticinputfailoversettings.html#cfn-medialive-channel-automaticinputfailoversettings-failoverconditions
	FailoverConditions []Channel_FailoverCondition `json:"FailoverConditions,omitempty"`

	// InputPreference AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-automaticinputfailoversettings.html#cfn-medialive-channel-automaticinputfailoversettings-inputpreference
	InputPreference *types.Value `json:"InputPreference,omitempty"`

	// SecondaryInputId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-automaticinputfailoversettings.html#cfn-medialive-channel-automaticinputfailoversettings-secondaryinputid
	SecondaryInputId *types.Value `json:"SecondaryInputId,omitempty"`

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
func (r *Channel_AutomaticInputFailoverSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.AutomaticInputFailoverSettings"
}
