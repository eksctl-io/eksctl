package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_SqlApplicationConfiguration AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.SqlApplicationConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-sqlapplicationconfiguration.html
type AWSKinesisAnalyticsV2Application_SqlApplicationConfiguration struct {

	// Inputs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-sqlapplicationconfiguration.html#cfn-kinesisanalyticsv2-application-sqlapplicationconfiguration-inputs
	Inputs []AWSKinesisAnalyticsV2Application_Input `json:"Inputs,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_SqlApplicationConfiguration) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.SqlApplicationConfiguration"
}

func (r *AWSKinesisAnalyticsV2Application_SqlApplicationConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
