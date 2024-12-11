package cognito

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// UserPool_VerificationMessageTemplate AWS CloudFormation Resource (AWS::Cognito::UserPool.VerificationMessageTemplate)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpool-verificationmessagetemplate.html
type UserPool_VerificationMessageTemplate struct {

	// DefaultEmailOption AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpool-verificationmessagetemplate.html#cfn-cognito-userpool-verificationmessagetemplate-defaultemailoption
	DefaultEmailOption *types.Value `json:"DefaultEmailOption,omitempty"`

	// EmailMessage AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpool-verificationmessagetemplate.html#cfn-cognito-userpool-verificationmessagetemplate-emailmessage
	EmailMessage *types.Value `json:"EmailMessage,omitempty"`

	// EmailMessageByLink AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpool-verificationmessagetemplate.html#cfn-cognito-userpool-verificationmessagetemplate-emailmessagebylink
	EmailMessageByLink *types.Value `json:"EmailMessageByLink,omitempty"`

	// EmailSubject AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpool-verificationmessagetemplate.html#cfn-cognito-userpool-verificationmessagetemplate-emailsubject
	EmailSubject *types.Value `json:"EmailSubject,omitempty"`

	// EmailSubjectByLink AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpool-verificationmessagetemplate.html#cfn-cognito-userpool-verificationmessagetemplate-emailsubjectbylink
	EmailSubjectByLink *types.Value `json:"EmailSubjectByLink,omitempty"`

	// SmsMessage AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpool-verificationmessagetemplate.html#cfn-cognito-userpool-verificationmessagetemplate-smsmessage
	SmsMessage *types.Value `json:"SmsMessage,omitempty"`

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
func (r *UserPool_VerificationMessageTemplate) AWSCloudFormationType() string {
	return "AWS::Cognito::UserPool.VerificationMessageTemplate"
}
