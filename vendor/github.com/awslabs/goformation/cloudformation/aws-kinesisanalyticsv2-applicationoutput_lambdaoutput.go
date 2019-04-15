package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2ApplicationOutput_LambdaOutput AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::ApplicationOutput.LambdaOutput)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationoutput-lambdaoutput.html
type AWSKinesisAnalyticsV2ApplicationOutput_LambdaOutput struct {

	// ResourceARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationoutput-lambdaoutput.html#cfn-kinesisanalyticsv2-applicationoutput-lambdaoutput-resourcearn
	ResourceARN *Value `json:"ResourceARN,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2ApplicationOutput_LambdaOutput) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::ApplicationOutput.LambdaOutput"
}

func (r *AWSKinesisAnalyticsV2ApplicationOutput_LambdaOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
