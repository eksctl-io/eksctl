package cloudformation

import (
	"encoding/json"
)

// AWSCloudFrontDistribution_LambdaFunctionAssociation AWS CloudFormation Resource (AWS::CloudFront::Distribution.LambdaFunctionAssociation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-lambdafunctionassociation.html
type AWSCloudFrontDistribution_LambdaFunctionAssociation struct {

	// EventType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-lambdafunctionassociation.html#cfn-cloudfront-distribution-lambdafunctionassociation-eventtype
	EventType *Value `json:"EventType,omitempty"`

	// LambdaFunctionARN AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-lambdafunctionassociation.html#cfn-cloudfront-distribution-lambdafunctionassociation-lambdafunctionarn
	LambdaFunctionARN *Value `json:"LambdaFunctionARN,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCloudFrontDistribution_LambdaFunctionAssociation) AWSCloudFormationType() string {
	return "AWS::CloudFront::Distribution.LambdaFunctionAssociation"
}

func (r *AWSCloudFrontDistribution_LambdaFunctionAssociation) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
