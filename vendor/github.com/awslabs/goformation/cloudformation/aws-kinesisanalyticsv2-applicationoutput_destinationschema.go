package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2ApplicationOutput_DestinationSchema AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::ApplicationOutput.DestinationSchema)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationoutput-destinationschema.html
type AWSKinesisAnalyticsV2ApplicationOutput_DestinationSchema struct {

	// RecordFormatType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationoutput-destinationschema.html#cfn-kinesisanalyticsv2-applicationoutput-destinationschema-recordformattype
	RecordFormatType *Value `json:"RecordFormatType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2ApplicationOutput_DestinationSchema) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::ApplicationOutput.DestinationSchema"
}

func (r *AWSKinesisAnalyticsV2ApplicationOutput_DestinationSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
