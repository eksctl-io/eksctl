package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption_CloudWatchLoggingOption AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::ApplicationCloudWatchLoggingOption.CloudWatchLoggingOption)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationcloudwatchloggingoption-cloudwatchloggingoption.html
type AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption_CloudWatchLoggingOption struct {

	// LogStreamARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationcloudwatchloggingoption-cloudwatchloggingoption.html#cfn-kinesisanalyticsv2-applicationcloudwatchloggingoption-cloudwatchloggingoption-logstreamarn
	LogStreamARN *Value `json:"LogStreamARN,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption_CloudWatchLoggingOption) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::ApplicationCloudWatchLoggingOption.CloudWatchLoggingOption"
}

func (r *AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption_CloudWatchLoggingOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
