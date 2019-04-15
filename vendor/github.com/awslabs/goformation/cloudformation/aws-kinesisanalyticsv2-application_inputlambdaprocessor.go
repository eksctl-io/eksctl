package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_InputLambdaProcessor AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.InputLambdaProcessor)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-inputlambdaprocessor.html
type AWSKinesisAnalyticsV2Application_InputLambdaProcessor struct {

	// ResourceARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-inputlambdaprocessor.html#cfn-kinesisanalyticsv2-application-inputlambdaprocessor-resourcearn
	ResourceARN *Value `json:"ResourceARN,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_InputLambdaProcessor) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.InputLambdaProcessor"
}

func (r *AWSKinesisAnalyticsV2Application_InputLambdaProcessor) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
