package cloudformation

// AWSSESReceiptRule_StopAction AWS CloudFormation Resource (AWS::SES::ReceiptRule.StopAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-stopaction.html
type AWSSESReceiptRule_StopAction struct {

	// Scope AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-stopaction.html#cfn-ses-receiptrule-stopaction-scope
	Scope *StringIntrinsic `json:"Scope,omitempty"`

	// TopicArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-stopaction.html#cfn-ses-receiptrule-stopaction-topicarn
	TopicArn *StringIntrinsic `json:"TopicArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSESReceiptRule_StopAction) AWSCloudFormationType() string {
	return "AWS::SES::ReceiptRule.StopAction"
}
