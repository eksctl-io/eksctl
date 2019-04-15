package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_KinesisFirehoseInput AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.KinesisFirehoseInput)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-kinesisfirehoseinput.html
type AWSKinesisAnalyticsV2Application_KinesisFirehoseInput struct {

	// ResourceARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-kinesisfirehoseinput.html#cfn-kinesisanalyticsv2-application-kinesisfirehoseinput-resourcearn
	ResourceARN *Value `json:"ResourceARN,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_KinesisFirehoseInput) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.KinesisFirehoseInput"
}

func (r *AWSKinesisAnalyticsV2Application_KinesisFirehoseInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
