package pinpoint

import (
	"goformation/v4/cloudformation/policies"
)

// Campaign_InAppMessageButton AWS CloudFormation Resource (AWS::Pinpoint::Campaign.InAppMessageButton)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-campaign-inappmessagebutton.html
type Campaign_InAppMessageButton struct {

	// Android AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-campaign-inappmessagebutton.html#cfn-pinpoint-campaign-inappmessagebutton-android
	Android *Campaign_OverrideButtonConfiguration `json:"Android,omitempty"`

	// DefaultConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-campaign-inappmessagebutton.html#cfn-pinpoint-campaign-inappmessagebutton-defaultconfig
	DefaultConfig *Campaign_DefaultButtonConfiguration `json:"DefaultConfig,omitempty"`

	// IOS AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-campaign-inappmessagebutton.html#cfn-pinpoint-campaign-inappmessagebutton-ios
	IOS *Campaign_OverrideButtonConfiguration `json:"IOS,omitempty"`

	// Web AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-campaign-inappmessagebutton.html#cfn-pinpoint-campaign-inappmessagebutton-web
	Web *Campaign_OverrideButtonConfiguration `json:"Web,omitempty"`

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
func (r *Campaign_InAppMessageButton) AWSCloudFormationType() string {
	return "AWS::Pinpoint::Campaign.InAppMessageButton"
}
