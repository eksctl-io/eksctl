package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsPipeline_Lambda AWS CloudFormation Resource (AWS::IoTAnalytics::Pipeline.Lambda)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-lambda.html
type AWSIoTAnalyticsPipeline_Lambda struct {

	// BatchSize AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-lambda.html#cfn-iotanalytics-pipeline-lambda-batchsize
	BatchSize *Value `json:"BatchSize,omitempty"`

	// LambdaName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-lambda.html#cfn-iotanalytics-pipeline-lambda-lambdaname
	LambdaName *Value `json:"LambdaName,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-lambda.html#cfn-iotanalytics-pipeline-lambda-name
	Name *Value `json:"Name,omitempty"`

	// Next AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-lambda.html#cfn-iotanalytics-pipeline-lambda-next
	Next *Value `json:"Next,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsPipeline_Lambda) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Pipeline.Lambda"
}

func (r *AWSIoTAnalyticsPipeline_Lambda) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
