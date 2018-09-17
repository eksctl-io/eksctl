package cloudformation

import (
	"encoding/json"
)

// AWSS3Bucket_MetricsConfiguration AWS CloudFormation Resource (AWS::S3::Bucket.MetricsConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-metricsconfiguration.html
type AWSS3Bucket_MetricsConfiguration struct {

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-metricsconfiguration.html#cfn-s3-bucket-metricsconfiguration-id
	Id *Value `json:"Id,omitempty"`

	// Prefix AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-metricsconfiguration.html#cfn-s3-bucket-metricsconfiguration-prefix
	Prefix *Value `json:"Prefix,omitempty"`

	// TagFilters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-metricsconfiguration.html#cfn-s3-bucket-metricsconfiguration-tagfilters
	TagFilters []AWSS3Bucket_TagFilter `json:"TagFilters,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSS3Bucket_MetricsConfiguration) AWSCloudFormationType() string {
	return "AWS::S3::Bucket.MetricsConfiguration"
}

func (r *AWSS3Bucket_MetricsConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
