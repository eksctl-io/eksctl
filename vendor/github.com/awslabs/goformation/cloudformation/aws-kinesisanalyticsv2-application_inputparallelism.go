package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_InputParallelism AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.InputParallelism)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-inputparallelism.html
type AWSKinesisAnalyticsV2Application_InputParallelism struct {

	// Count AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-inputparallelism.html#cfn-kinesisanalyticsv2-application-inputparallelism-count
	Count *Value `json:"Count,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_InputParallelism) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.InputParallelism"
}

func (r *AWSKinesisAnalyticsV2Application_InputParallelism) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
