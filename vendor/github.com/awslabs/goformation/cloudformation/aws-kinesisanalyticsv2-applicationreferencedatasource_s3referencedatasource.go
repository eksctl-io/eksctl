package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2ApplicationReferenceDataSource_S3ReferenceDataSource AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource.S3ReferenceDataSource)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-s3referencedatasource.html
type AWSKinesisAnalyticsV2ApplicationReferenceDataSource_S3ReferenceDataSource struct {

	// BucketARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-s3referencedatasource.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-s3referencedatasource-bucketarn
	BucketARN *Value `json:"BucketARN,omitempty"`

	// FileKey AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-s3referencedatasource.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-s3referencedatasource-filekey
	FileKey *Value `json:"FileKey,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_S3ReferenceDataSource) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource.S3ReferenceDataSource"
}

func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_S3ReferenceDataSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
