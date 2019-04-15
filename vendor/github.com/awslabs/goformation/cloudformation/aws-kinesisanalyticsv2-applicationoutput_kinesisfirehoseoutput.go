package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2ApplicationOutput_KinesisFirehoseOutput AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::ApplicationOutput.KinesisFirehoseOutput)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationoutput-kinesisfirehoseoutput.html
type AWSKinesisAnalyticsV2ApplicationOutput_KinesisFirehoseOutput struct {

	// ResourceARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationoutput-kinesisfirehoseoutput.html#cfn-kinesisanalyticsv2-applicationoutput-kinesisfirehoseoutput-resourcearn
	ResourceARN *Value `json:"ResourceARN,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2ApplicationOutput_KinesisFirehoseOutput) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::ApplicationOutput.KinesisFirehoseOutput"
}

func (r *AWSKinesisAnalyticsV2ApplicationOutput_KinesisFirehoseOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
