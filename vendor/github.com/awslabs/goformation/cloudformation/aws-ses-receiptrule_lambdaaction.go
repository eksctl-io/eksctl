package cloudformation

import (
	"encoding/json"
)

// AWSSESReceiptRule_LambdaAction AWS CloudFormation Resource (AWS::SES::ReceiptRule.LambdaAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-lambdaaction.html
type AWSSESReceiptRule_LambdaAction struct {

	// FunctionArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-lambdaaction.html#cfn-ses-receiptrule-lambdaaction-functionarn
	FunctionArn *Value `json:"FunctionArn,omitempty"`

	// InvocationType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-lambdaaction.html#cfn-ses-receiptrule-lambdaaction-invocationtype
	InvocationType *Value `json:"InvocationType,omitempty"`

	// TopicArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-lambdaaction.html#cfn-ses-receiptrule-lambdaaction-topicarn
	TopicArn *Value `json:"TopicArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSESReceiptRule_LambdaAction) AWSCloudFormationType() string {
	return "AWS::SES::ReceiptRule.LambdaAction"
}

func (r *AWSSESReceiptRule_LambdaAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
