package cloudformation

import (
	"encoding/json"
)

// AWSCloudFrontDistribution_S3OriginConfig AWS CloudFormation Resource (AWS::CloudFront::Distribution.S3OriginConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-s3originconfig.html
type AWSCloudFrontDistribution_S3OriginConfig struct {

	// OriginAccessIdentity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-s3originconfig.html#cfn-cloudfront-distribution-s3originconfig-originaccessidentity
	OriginAccessIdentity *Value `json:"OriginAccessIdentity,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCloudFrontDistribution_S3OriginConfig) AWSCloudFormationType() string {
	return "AWS::CloudFront::Distribution.S3OriginConfig"
}

func (r *AWSCloudFrontDistribution_S3OriginConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
