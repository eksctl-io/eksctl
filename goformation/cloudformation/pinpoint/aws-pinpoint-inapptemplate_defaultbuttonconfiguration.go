package pinpoint

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// InAppTemplate_DefaultButtonConfiguration AWS CloudFormation Resource (AWS::Pinpoint::InAppTemplate.DefaultButtonConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-inapptemplate-defaultbuttonconfiguration.html
type InAppTemplate_DefaultButtonConfiguration struct {

	// BackgroundColor AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-inapptemplate-defaultbuttonconfiguration.html#cfn-pinpoint-inapptemplate-defaultbuttonconfiguration-backgroundcolor
	BackgroundColor *types.Value `json:"BackgroundColor,omitempty"`

	// BorderRadius AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-inapptemplate-defaultbuttonconfiguration.html#cfn-pinpoint-inapptemplate-defaultbuttonconfiguration-borderradius
	BorderRadius *types.Value `json:"BorderRadius,omitempty"`

	// ButtonAction AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-inapptemplate-defaultbuttonconfiguration.html#cfn-pinpoint-inapptemplate-defaultbuttonconfiguration-buttonaction
	ButtonAction *types.Value `json:"ButtonAction,omitempty"`

	// Link AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-inapptemplate-defaultbuttonconfiguration.html#cfn-pinpoint-inapptemplate-defaultbuttonconfiguration-link
	Link *types.Value `json:"Link,omitempty"`

	// Text AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-inapptemplate-defaultbuttonconfiguration.html#cfn-pinpoint-inapptemplate-defaultbuttonconfiguration-text
	Text *types.Value `json:"Text,omitempty"`

	// TextColor AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-inapptemplate-defaultbuttonconfiguration.html#cfn-pinpoint-inapptemplate-defaultbuttonconfiguration-textcolor
	TextColor *types.Value `json:"TextColor,omitempty"`

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
func (r *InAppTemplate_DefaultButtonConfiguration) AWSCloudFormationType() string {
	return "AWS::Pinpoint::InAppTemplate.DefaultButtonConfiguration"
}
