package cloudformation

// AWSSESReceiptRule_SNSAction AWS CloudFormation Resource (AWS::SES::ReceiptRule.SNSAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-snsaction.html
type AWSSESReceiptRule_SNSAction struct {

	// Encoding AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-snsaction.html#cfn-ses-receiptrule-snsaction-encoding
	Encoding *StringIntrinsic `json:"Encoding,omitempty"`

	// TopicArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-snsaction.html#cfn-ses-receiptrule-snsaction-topicarn
	TopicArn *StringIntrinsic `json:"TopicArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSESReceiptRule_SNSAction) AWSCloudFormationType() string {
	return "AWS::SES::ReceiptRule.SNSAction"
}
