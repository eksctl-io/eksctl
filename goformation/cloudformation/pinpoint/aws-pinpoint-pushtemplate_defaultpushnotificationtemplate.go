package pinpoint

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// PushTemplate_DefaultPushNotificationTemplate AWS CloudFormation Resource (AWS::Pinpoint::PushTemplate.DefaultPushNotificationTemplate)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-pushtemplate-defaultpushnotificationtemplate.html
type PushTemplate_DefaultPushNotificationTemplate struct {

	// Action AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-pushtemplate-defaultpushnotificationtemplate.html#cfn-pinpoint-pushtemplate-defaultpushnotificationtemplate-action
	Action *types.Value `json:"Action,omitempty"`

	// Body AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-pushtemplate-defaultpushnotificationtemplate.html#cfn-pinpoint-pushtemplate-defaultpushnotificationtemplate-body
	Body *types.Value `json:"Body,omitempty"`

	// Sound AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-pushtemplate-defaultpushnotificationtemplate.html#cfn-pinpoint-pushtemplate-defaultpushnotificationtemplate-sound
	Sound *types.Value `json:"Sound,omitempty"`

	// Title AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-pushtemplate-defaultpushnotificationtemplate.html#cfn-pinpoint-pushtemplate-defaultpushnotificationtemplate-title
	Title *types.Value `json:"Title,omitempty"`

	// Url AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-pushtemplate-defaultpushnotificationtemplate.html#cfn-pinpoint-pushtemplate-defaultpushnotificationtemplate-url
	Url *types.Value `json:"Url,omitempty"`

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
func (r *PushTemplate_DefaultPushNotificationTemplate) AWSCloudFormationType() string {
	return "AWS::Pinpoint::PushTemplate.DefaultPushNotificationTemplate"
}
