package cloudformation

// AWSCloudFrontDistribution_Logging AWS CloudFormation Resource (AWS::CloudFront::Distribution.Logging)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-logging.html
type AWSCloudFrontDistribution_Logging struct {

	// Bucket AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-logging.html#cfn-cloudfront-distribution-logging-bucket
	Bucket *StringIntrinsic `json:"Bucket,omitempty"`

	// IncludeCookies AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-logging.html#cfn-cloudfront-distribution-logging-includecookies
	IncludeCookies bool `json:"IncludeCookies,omitempty"`

	// Prefix AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-logging.html#cfn-cloudfront-distribution-logging-prefix
	Prefix *StringIntrinsic `json:"Prefix,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCloudFrontDistribution_Logging) AWSCloudFormationType() string {
	return "AWS::CloudFront::Distribution.Logging"
}
