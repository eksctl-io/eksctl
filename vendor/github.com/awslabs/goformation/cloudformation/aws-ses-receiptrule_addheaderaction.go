package cloudformation

// AWSSESReceiptRule_AddHeaderAction AWS CloudFormation Resource (AWS::SES::ReceiptRule.AddHeaderAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-addheaderaction.html
type AWSSESReceiptRule_AddHeaderAction struct {

	// HeaderName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-addheaderaction.html#cfn-ses-receiptrule-addheaderaction-headername
	HeaderName *StringIntrinsic `json:"HeaderName,omitempty"`

	// HeaderValue AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-addheaderaction.html#cfn-ses-receiptrule-addheaderaction-headervalue
	HeaderValue *StringIntrinsic `json:"HeaderValue,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSESReceiptRule_AddHeaderAction) AWSCloudFormationType() string {
	return "AWS::SES::ReceiptRule.AddHeaderAction"
}
