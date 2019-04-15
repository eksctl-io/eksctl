package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_RecordFormat AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.RecordFormat)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-recordformat.html
type AWSKinesisAnalyticsV2Application_RecordFormat struct {

	// MappingParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-recordformat.html#cfn-kinesisanalyticsv2-application-recordformat-mappingparameters
	MappingParameters *AWSKinesisAnalyticsV2Application_MappingParameters `json:"MappingParameters,omitempty"`

	// RecordFormatType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-recordformat.html#cfn-kinesisanalyticsv2-application-recordformat-recordformattype
	RecordFormatType *Value `json:"RecordFormatType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_RecordFormat) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.RecordFormat"
}

func (r *AWSKinesisAnalyticsV2Application_RecordFormat) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
