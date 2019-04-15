package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2ApplicationReferenceDataSource_MappingParameters AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource.MappingParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-mappingparameters.html
type AWSKinesisAnalyticsV2ApplicationReferenceDataSource_MappingParameters struct {

	// CSVMappingParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-mappingparameters.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-mappingparameters-csvmappingparameters
	CSVMappingParameters *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_CSVMappingParameters `json:"CSVMappingParameters,omitempty"`

	// JSONMappingParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-mappingparameters.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-mappingparameters-jsonmappingparameters
	JSONMappingParameters *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_JSONMappingParameters `json:"JSONMappingParameters,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_MappingParameters) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource.MappingParameters"
}

func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_MappingParameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
