package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsApplication_InputLambdaProcessor AWS CloudFormation Resource (AWS::KinesisAnalytics::Application.InputLambdaProcessor)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-application-inputlambdaprocessor.html
type AWSKinesisAnalyticsApplication_InputLambdaProcessor struct {

	// ResourceARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-application-inputlambdaprocessor.html#cfn-kinesisanalytics-application-inputlambdaprocessor-resourcearn
	ResourceARN *Value `json:"ResourceARN,omitempty"`

	// RoleARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-application-inputlambdaprocessor.html#cfn-kinesisanalytics-application-inputlambdaprocessor-rolearn
	RoleARN *Value `json:"RoleARN,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsApplication_InputLambdaProcessor) AWSCloudFormationType() string {
	return "AWS::KinesisAnalytics::Application.InputLambdaProcessor"
}

func (r *AWSKinesisAnalyticsApplication_InputLambdaProcessor) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
