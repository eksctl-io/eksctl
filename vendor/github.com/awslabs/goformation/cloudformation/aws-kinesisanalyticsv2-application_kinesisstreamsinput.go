package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_KinesisStreamsInput AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.KinesisStreamsInput)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-kinesisstreamsinput.html
type AWSKinesisAnalyticsV2Application_KinesisStreamsInput struct {

	// ResourceARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-kinesisstreamsinput.html#cfn-kinesisanalyticsv2-application-kinesisstreamsinput-resourcearn
	ResourceARN *Value `json:"ResourceARN,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_KinesisStreamsInput) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.KinesisStreamsInput"
}

func (r *AWSKinesisAnalyticsV2Application_KinesisStreamsInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
