package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2ApplicationReferenceDataSource_ReferenceSchema AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource.ReferenceSchema)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-referenceschema.html
type AWSKinesisAnalyticsV2ApplicationReferenceDataSource_ReferenceSchema struct {

	// RecordColumns AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-referenceschema.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-referenceschema-recordcolumns
	RecordColumns []AWSKinesisAnalyticsV2ApplicationReferenceDataSource_RecordColumn `json:"RecordColumns,omitempty"`

	// RecordEncoding AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-referenceschema.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-referenceschema-recordencoding
	RecordEncoding *Value `json:"RecordEncoding,omitempty"`

	// RecordFormat AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-referenceschema.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-referenceschema-recordformat
	RecordFormat *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_RecordFormat `json:"RecordFormat,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_ReferenceSchema) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource.ReferenceSchema"
}

func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_ReferenceSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
