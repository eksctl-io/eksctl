package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2ApplicationReferenceDataSource_RecordFormat AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource.RecordFormat)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-recordformat.html
type AWSKinesisAnalyticsV2ApplicationReferenceDataSource_RecordFormat struct {

	// MappingParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-recordformat.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-recordformat-mappingparameters
	MappingParameters *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_MappingParameters `json:"MappingParameters,omitempty"`

	// RecordFormatType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-recordformat.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-recordformat-recordformattype
	RecordFormatType *Value `json:"RecordFormatType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_RecordFormat) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource.RecordFormat"
}

func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_RecordFormat) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
