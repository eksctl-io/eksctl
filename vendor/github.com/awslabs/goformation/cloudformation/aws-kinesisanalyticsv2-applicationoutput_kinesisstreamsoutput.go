package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2ApplicationOutput_KinesisStreamsOutput AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::ApplicationOutput.KinesisStreamsOutput)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationoutput-kinesisstreamsoutput.html
type AWSKinesisAnalyticsV2ApplicationOutput_KinesisStreamsOutput struct {

	// ResourceARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationoutput-kinesisstreamsoutput.html#cfn-kinesisanalyticsv2-applicationoutput-kinesisstreamsoutput-resourcearn
	ResourceARN *Value `json:"ResourceARN,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2ApplicationOutput_KinesisStreamsOutput) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::ApplicationOutput.KinesisStreamsOutput"
}

func (r *AWSKinesisAnalyticsV2ApplicationOutput_KinesisStreamsOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
