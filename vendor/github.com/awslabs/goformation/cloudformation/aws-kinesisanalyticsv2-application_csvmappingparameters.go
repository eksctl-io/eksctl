package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_CSVMappingParameters AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.CSVMappingParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-csvmappingparameters.html
type AWSKinesisAnalyticsV2Application_CSVMappingParameters struct {

	// RecordColumnDelimiter AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-csvmappingparameters.html#cfn-kinesisanalyticsv2-application-csvmappingparameters-recordcolumndelimiter
	RecordColumnDelimiter *Value `json:"RecordColumnDelimiter,omitempty"`

	// RecordRowDelimiter AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-csvmappingparameters.html#cfn-kinesisanalyticsv2-application-csvmappingparameters-recordrowdelimiter
	RecordRowDelimiter *Value `json:"RecordRowDelimiter,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_CSVMappingParameters) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.CSVMappingParameters"
}

func (r *AWSKinesisAnalyticsV2Application_CSVMappingParameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
