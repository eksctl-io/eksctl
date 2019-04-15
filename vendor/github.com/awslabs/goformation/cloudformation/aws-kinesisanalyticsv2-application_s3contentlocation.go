package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_S3ContentLocation AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.S3ContentLocation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-s3contentlocation.html
type AWSKinesisAnalyticsV2Application_S3ContentLocation struct {

	// BucketARN AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-s3contentlocation.html#cfn-kinesisanalyticsv2-application-s3contentlocation-bucketarn
	BucketARN *Value `json:"BucketARN,omitempty"`

	// FileKey AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-s3contentlocation.html#cfn-kinesisanalyticsv2-application-s3contentlocation-filekey
	FileKey *Value `json:"FileKey,omitempty"`

	// ObjectVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-s3contentlocation.html#cfn-kinesisanalyticsv2-application-s3contentlocation-objectversion
	ObjectVersion *Value `json:"ObjectVersion,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_S3ContentLocation) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.S3ContentLocation"
}

func (r *AWSKinesisAnalyticsV2Application_S3ContentLocation) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
