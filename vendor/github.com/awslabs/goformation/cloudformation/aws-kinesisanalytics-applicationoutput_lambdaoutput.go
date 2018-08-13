package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsApplicationOutput_LambdaOutput AWS CloudFormation Resource (AWS::KinesisAnalytics::ApplicationOutput.LambdaOutput)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-applicationoutput-lambdaoutput.html
type AWSKinesisAnalyticsApplicationOutput_LambdaOutput struct {

	// ResourceARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-applicationoutput-lambdaoutput.html#cfn-kinesisanalytics-applicationoutput-lambdaoutput-resourcearn
	ResourceARN *Value `json:"ResourceARN,omitempty"`

	// RoleARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-applicationoutput-lambdaoutput.html#cfn-kinesisanalytics-applicationoutput-lambdaoutput-rolearn
	RoleARN *Value `json:"RoleARN,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsApplicationOutput_LambdaOutput) AWSCloudFormationType() string {
	return "AWS::KinesisAnalytics::ApplicationOutput.LambdaOutput"
}

func (r *AWSKinesisAnalyticsApplicationOutput_LambdaOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
