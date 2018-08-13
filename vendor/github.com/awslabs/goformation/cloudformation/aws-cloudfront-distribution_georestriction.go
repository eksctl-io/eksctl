package cloudformation

import (
	"encoding/json"
)

// AWSCloudFrontDistribution_GeoRestriction AWS CloudFormation Resource (AWS::CloudFront::Distribution.GeoRestriction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-georestriction.html
type AWSCloudFrontDistribution_GeoRestriction struct {

	// Locations AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-georestriction.html#cfn-cloudfront-distribution-georestriction-locations
	Locations []*Value `json:"Locations,omitempty"`

	// RestrictionType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-georestriction.html#cfn-cloudfront-distribution-georestriction-restrictiontype
	RestrictionType *Value `json:"RestrictionType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCloudFrontDistribution_GeoRestriction) AWSCloudFormationType() string {
	return "AWS::CloudFront::Distribution.GeoRestriction"
}

func (r *AWSCloudFrontDistribution_GeoRestriction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
