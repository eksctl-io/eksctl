package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsApplicationOutput_KinesisFirehoseOutput AWS CloudFormation Resource (AWS::KinesisAnalytics::ApplicationOutput.KinesisFirehoseOutput)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-applicationoutput-kinesisfirehoseoutput.html
type AWSKinesisAnalyticsApplicationOutput_KinesisFirehoseOutput struct {

	// ResourceARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-applicationoutput-kinesisfirehoseoutput.html#cfn-kinesisanalytics-applicationoutput-kinesisfirehoseoutput-resourcearn
	ResourceARN *Value `json:"ResourceARN,omitempty"`

	// RoleARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-applicationoutput-kinesisfirehoseoutput.html#cfn-kinesisanalytics-applicationoutput-kinesisfirehoseoutput-rolearn
	RoleARN *Value `json:"RoleARN,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsApplicationOutput_KinesisFirehoseOutput) AWSCloudFormationType() string {
	return "AWS::KinesisAnalytics::ApplicationOutput.KinesisFirehoseOutput"
}

func (r *AWSKinesisAnalyticsApplicationOutput_KinesisFirehoseOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
