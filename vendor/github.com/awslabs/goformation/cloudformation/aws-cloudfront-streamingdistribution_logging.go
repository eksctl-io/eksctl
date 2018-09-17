package cloudformation

import (
	"encoding/json"
)

// AWSCloudFrontStreamingDistribution_Logging AWS CloudFormation Resource (AWS::CloudFront::StreamingDistribution.Logging)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-streamingdistribution-logging.html
type AWSCloudFrontStreamingDistribution_Logging struct {

	// Bucket AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-streamingdistribution-logging.html#cfn-cloudfront-streamingdistribution-logging-bucket
	Bucket *Value `json:"Bucket,omitempty"`

	// Enabled AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-streamingdistribution-logging.html#cfn-cloudfront-streamingdistribution-logging-enabled
	Enabled *Value `json:"Enabled,omitempty"`

	// Prefix AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-streamingdistribution-logging.html#cfn-cloudfront-streamingdistribution-logging-prefix
	Prefix *Value `json:"Prefix,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCloudFrontStreamingDistribution_Logging) AWSCloudFormationType() string {
	return "AWS::CloudFront::StreamingDistribution.Logging"
}

func (r *AWSCloudFrontStreamingDistribution_Logging) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
