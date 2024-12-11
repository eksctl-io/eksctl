package pinpoint

import (
	"goformation/v4/cloudformation/policies"
)

// InAppTemplate_ButtonConfig AWS CloudFormation Resource (AWS::Pinpoint::InAppTemplate.ButtonConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-inapptemplate-buttonconfig.html
type InAppTemplate_ButtonConfig struct {

	// Android AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-inapptemplate-buttonconfig.html#cfn-pinpoint-inapptemplate-buttonconfig-android
	Android *InAppTemplate_OverrideButtonConfiguration `json:"Android,omitempty"`

	// DefaultConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-inapptemplate-buttonconfig.html#cfn-pinpoint-inapptemplate-buttonconfig-defaultconfig
	DefaultConfig *InAppTemplate_DefaultButtonConfiguration `json:"DefaultConfig,omitempty"`

	// IOS AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-inapptemplate-buttonconfig.html#cfn-pinpoint-inapptemplate-buttonconfig-ios
	IOS *InAppTemplate_OverrideButtonConfiguration `json:"IOS,omitempty"`

	// Web AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-inapptemplate-buttonconfig.html#cfn-pinpoint-inapptemplate-buttonconfig-web
	Web *InAppTemplate_OverrideButtonConfiguration `json:"Web,omitempty"`

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
func (r *InAppTemplate_ButtonConfig) AWSCloudFormationType() string {
	return "AWS::Pinpoint::InAppTemplate.ButtonConfig"
}
