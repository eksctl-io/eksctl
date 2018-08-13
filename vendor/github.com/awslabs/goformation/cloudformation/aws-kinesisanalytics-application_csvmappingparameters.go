package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsApplication_CSVMappingParameters AWS CloudFormation Resource (AWS::KinesisAnalytics::Application.CSVMappingParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-application-csvmappingparameters.html
type AWSKinesisAnalyticsApplication_CSVMappingParameters struct {

	// RecordColumnDelimiter AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-application-csvmappingparameters.html#cfn-kinesisanalytics-application-csvmappingparameters-recordcolumndelimiter
	RecordColumnDelimiter *Value `json:"RecordColumnDelimiter,omitempty"`

	// RecordRowDelimiter AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-application-csvmappingparameters.html#cfn-kinesisanalytics-application-csvmappingparameters-recordrowdelimiter
	RecordRowDelimiter *Value `json:"RecordRowDelimiter,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsApplication_CSVMappingParameters) AWSCloudFormationType() string {
	return "AWS::KinesisAnalytics::Application.CSVMappingParameters"
}

func (r *AWSKinesisAnalyticsApplication_CSVMappingParameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
