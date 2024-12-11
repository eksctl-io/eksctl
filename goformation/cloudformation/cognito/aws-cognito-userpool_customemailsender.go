package cognito

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// UserPool_CustomEmailSender AWS CloudFormation Resource (AWS::Cognito::UserPool.CustomEmailSender)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpool-customemailsender.html
type UserPool_CustomEmailSender struct {

	// LambdaArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpool-customemailsender.html#cfn-cognito-userpool-customemailsender-lambdaarn
	LambdaArn *types.Value `json:"LambdaArn,omitempty"`

	// LambdaVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpool-customemailsender.html#cfn-cognito-userpool-customemailsender-lambdaversion
	LambdaVersion *types.Value `json:"LambdaVersion,omitempty"`

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
func (r *UserPool_CustomEmailSender) AWSCloudFormationType() string {
	return "AWS::Cognito::UserPool.CustomEmailSender"
}
