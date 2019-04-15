package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2ApplicationReferenceDataSource_CSVMappingParameters AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource.CSVMappingParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-csvmappingparameters.html
type AWSKinesisAnalyticsV2ApplicationReferenceDataSource_CSVMappingParameters struct {

	// RecordColumnDelimiter AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-csvmappingparameters.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-csvmappingparameters-recordcolumndelimiter
	RecordColumnDelimiter *Value `json:"RecordColumnDelimiter,omitempty"`

	// RecordRowDelimiter AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-csvmappingparameters.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-csvmappingparameters-recordrowdelimiter
	RecordRowDelimiter *Value `json:"RecordRowDelimiter,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_CSVMappingParameters) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource.CSVMappingParameters"
}

func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_CSVMappingParameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
