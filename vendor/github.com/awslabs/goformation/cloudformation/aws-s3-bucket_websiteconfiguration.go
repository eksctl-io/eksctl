package cloudformation

import (
	"encoding/json"
)

// AWSS3Bucket_WebsiteConfiguration AWS CloudFormation Resource (AWS::S3::Bucket.WebsiteConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-websiteconfiguration.html
type AWSS3Bucket_WebsiteConfiguration struct {

	// ErrorDocument AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-websiteconfiguration.html#cfn-s3-websiteconfiguration-errordocument
	ErrorDocument *Value `json:"ErrorDocument,omitempty"`

	// IndexDocument AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-websiteconfiguration.html#cfn-s3-websiteconfiguration-indexdocument
	IndexDocument *Value `json:"IndexDocument,omitempty"`

	// RedirectAllRequestsTo AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-websiteconfiguration.html#cfn-s3-websiteconfiguration-redirectallrequeststo
	RedirectAllRequestsTo *AWSS3Bucket_RedirectAllRequestsTo `json:"RedirectAllRequestsTo,omitempty"`

	// RoutingRules AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-websiteconfiguration.html#cfn-s3-websiteconfiguration-routingrules
	RoutingRules []AWSS3Bucket_RoutingRule `json:"RoutingRules,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSS3Bucket_WebsiteConfiguration) AWSCloudFormationType() string {
	return "AWS::S3::Bucket.WebsiteConfiguration"
}

func (r *AWSS3Bucket_WebsiteConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
