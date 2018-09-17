package cloudformation

import (
	"encoding/json"
)

// AWSCloudTrailTrail_DataResource AWS CloudFormation Resource (AWS::CloudTrail::Trail.DataResource)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudtrail-trail-dataresource.html
type AWSCloudTrailTrail_DataResource struct {

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudtrail-trail-dataresource.html#cfn-cloudtrail-trail-dataresource-type
	Type *Value `json:"Type,omitempty"`

	// Values AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudtrail-trail-dataresource.html#cfn-cloudtrail-trail-dataresource-values
	Values []*Value `json:"Values,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCloudTrailTrail_DataResource) AWSCloudFormationType() string {
	return "AWS::CloudTrail::Trail.DataResource"
}

func (r *AWSCloudTrailTrail_DataResource) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
