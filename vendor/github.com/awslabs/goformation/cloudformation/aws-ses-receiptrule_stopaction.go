package cloudformation

import (
	"encoding/json"
)

// AWSSESReceiptRule_StopAction AWS CloudFormation Resource (AWS::SES::ReceiptRule.StopAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-stopaction.html
type AWSSESReceiptRule_StopAction struct {

	// Scope AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-stopaction.html#cfn-ses-receiptrule-stopaction-scope
	Scope *Value `json:"Scope,omitempty"`

	// TopicArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-stopaction.html#cfn-ses-receiptrule-stopaction-topicarn
	TopicArn *Value `json:"TopicArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSESReceiptRule_StopAction) AWSCloudFormationType() string {
	return "AWS::SES::ReceiptRule.StopAction"
}

func (r *AWSSESReceiptRule_StopAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
