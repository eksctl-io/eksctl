package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_InputProcessingConfiguration AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.InputProcessingConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-inputprocessingconfiguration.html
type AWSKinesisAnalyticsV2Application_InputProcessingConfiguration struct {

	// InputLambdaProcessor AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-inputprocessingconfiguration.html#cfn-kinesisanalyticsv2-application-inputprocessingconfiguration-inputlambdaprocessor
	InputLambdaProcessor *AWSKinesisAnalyticsV2Application_InputLambdaProcessor `json:"InputLambdaProcessor,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_InputProcessingConfiguration) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.InputProcessingConfiguration"
}

func (r *AWSKinesisAnalyticsV2Application_InputProcessingConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
