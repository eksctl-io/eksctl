package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2ApplicationReferenceDataSource_ReferenceDataSource AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource.ReferenceDataSource)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-referencedatasource.html
type AWSKinesisAnalyticsV2ApplicationReferenceDataSource_ReferenceDataSource struct {

	// ReferenceSchema AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-referencedatasource.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-referencedatasource-referenceschema
	ReferenceSchema *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_ReferenceSchema `json:"ReferenceSchema,omitempty"`

	// S3ReferenceDataSource AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-referencedatasource.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-referencedatasource-s3referencedatasource
	S3ReferenceDataSource *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_S3ReferenceDataSource `json:"S3ReferenceDataSource,omitempty"`

	// TableName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-referencedatasource.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-referencedatasource-tablename
	TableName *Value `json:"TableName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_ReferenceDataSource) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource.ReferenceDataSource"
}

func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_ReferenceDataSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
