package cloudformation

import (
	"encoding/json"
)

// AWSCloudFrontDistribution_Restrictions AWS CloudFormation Resource (AWS::CloudFront::Distribution.Restrictions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-restrictions.html
type AWSCloudFrontDistribution_Restrictions struct {

	// GeoRestriction AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-restrictions.html#cfn-cloudfront-distribution-restrictions-georestriction
	GeoRestriction *AWSCloudFrontDistribution_GeoRestriction `json:"GeoRestriction,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCloudFrontDistribution_Restrictions) AWSCloudFormationType() string {
	return "AWS::CloudFront::Distribution.Restrictions"
}

func (r *AWSCloudFrontDistribution_Restrictions) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
