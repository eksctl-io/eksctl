package cloudformation

// AWSCloudFrontDistribution_ViewerCertificate AWS CloudFormation Resource (AWS::CloudFront::Distribution.ViewerCertificate)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-viewercertificate.html
type AWSCloudFrontDistribution_ViewerCertificate struct {

	// AcmCertificateArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-viewercertificate.html#cfn-cloudfront-distribution-viewercertificate-acmcertificatearn
	AcmCertificateArn string `json:"AcmCertificateArn,omitempty"`

	// CloudFrontDefaultCertificate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-viewercertificate.html#cfn-cloudfront-distribution-viewercertificate-cloudfrontdefaultcertificate
	CloudFrontDefaultCertificate bool `json:"CloudFrontDefaultCertificate,omitempty"`

	// IamCertificateId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-viewercertificate.html#cfn-cloudfront-distribution-viewercertificate-iamcertificateid
	IamCertificateId string `json:"IamCertificateId,omitempty"`

	// MinimumProtocolVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-viewercertificate.html#cfn-cloudfront-distribution-viewercertificate-minimumprotocolversion
	MinimumProtocolVersion string `json:"MinimumProtocolVersion,omitempty"`

	// SslSupportMethod AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-viewercertificate.html#cfn-cloudfront-distribution-viewercertificate-sslsupportmethod
	SslSupportMethod string `json:"SslSupportMethod,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCloudFrontDistribution_ViewerCertificate) AWSCloudFormationType() string {
	return "AWS::CloudFront::Distribution.ViewerCertificate"
}
