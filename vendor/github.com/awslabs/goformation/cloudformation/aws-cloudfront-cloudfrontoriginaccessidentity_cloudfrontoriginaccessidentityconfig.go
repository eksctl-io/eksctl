package cloudformation

import (
	"encoding/json"
)

// AWSCloudFrontCloudFrontOriginAccessIdentity_CloudFrontOriginAccessIdentityConfig AWS CloudFormation Resource (AWS::CloudFront::CloudFrontOriginAccessIdentity.CloudFrontOriginAccessIdentityConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-cloudfrontoriginaccessidentity-cloudfrontoriginaccessidentityconfig.html
type AWSCloudFrontCloudFrontOriginAccessIdentity_CloudFrontOriginAccessIdentityConfig struct {

	// Comment AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-cloudfrontoriginaccessidentity-cloudfrontoriginaccessidentityconfig.html#cfn-cloudfront-cloudfrontoriginaccessidentity-cloudfrontoriginaccessidentityconfig-comment
	Comment *Value `json:"Comment,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCloudFrontCloudFrontOriginAccessIdentity_CloudFrontOriginAccessIdentityConfig) AWSCloudFormationType() string {
	return "AWS::CloudFront::CloudFrontOriginAccessIdentity.CloudFrontOriginAccessIdentityConfig"
}

func (r *AWSCloudFrontCloudFrontOriginAccessIdentity_CloudFrontOriginAccessIdentityConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
