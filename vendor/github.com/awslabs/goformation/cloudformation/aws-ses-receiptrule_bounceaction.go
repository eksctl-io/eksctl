package cloudformation

// AWSSESReceiptRule_BounceAction AWS CloudFormation Resource (AWS::SES::ReceiptRule.BounceAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-bounceaction.html
type AWSSESReceiptRule_BounceAction struct {

	// Message AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-bounceaction.html#cfn-ses-receiptrule-bounceaction-message
	Message *StringIntrinsic `json:"Message,omitempty"`

	// Sender AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-bounceaction.html#cfn-ses-receiptrule-bounceaction-sender
	Sender *StringIntrinsic `json:"Sender,omitempty"`

	// SmtpReplyCode AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-bounceaction.html#cfn-ses-receiptrule-bounceaction-smtpreplycode
	SmtpReplyCode *StringIntrinsic `json:"SmtpReplyCode,omitempty"`

	// StatusCode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-bounceaction.html#cfn-ses-receiptrule-bounceaction-statuscode
	StatusCode *StringIntrinsic `json:"StatusCode,omitempty"`

	// TopicArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-bounceaction.html#cfn-ses-receiptrule-bounceaction-topicarn
	TopicArn *StringIntrinsic `json:"TopicArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSESReceiptRule_BounceAction) AWSCloudFormationType() string {
	return "AWS::SES::ReceiptRule.BounceAction"
}
