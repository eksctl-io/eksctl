package cloudformation

import (
	"encoding/json"
)

// AWSSESReceiptRule_BounceAction AWS CloudFormation Resource (AWS::SES::ReceiptRule.BounceAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-bounceaction.html
type AWSSESReceiptRule_BounceAction struct {

	// Message AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-bounceaction.html#cfn-ses-receiptrule-bounceaction-message
	Message *Value `json:"Message,omitempty"`

	// Sender AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-bounceaction.html#cfn-ses-receiptrule-bounceaction-sender
	Sender *Value `json:"Sender,omitempty"`

	// SmtpReplyCode AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-bounceaction.html#cfn-ses-receiptrule-bounceaction-smtpreplycode
	SmtpReplyCode *Value `json:"SmtpReplyCode,omitempty"`

	// StatusCode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-bounceaction.html#cfn-ses-receiptrule-bounceaction-statuscode
	StatusCode *Value `json:"StatusCode,omitempty"`

	// TopicArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-bounceaction.html#cfn-ses-receiptrule-bounceaction-topicarn
	TopicArn *Value `json:"TopicArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSESReceiptRule_BounceAction) AWSCloudFormationType() string {
	return "AWS::SES::ReceiptRule.BounceAction"
}

func (r *AWSSESReceiptRule_BounceAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
