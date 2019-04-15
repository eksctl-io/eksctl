package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_JSONMappingParameters AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.JSONMappingParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-jsonmappingparameters.html
type AWSKinesisAnalyticsV2Application_JSONMappingParameters struct {

	// RecordRowPath AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-jsonmappingparameters.html#cfn-kinesisanalyticsv2-application-jsonmappingparameters-recordrowpath
	RecordRowPath *Value `json:"RecordRowPath,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_JSONMappingParameters) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.JSONMappingParameters"
}

func (r *AWSKinesisAnalyticsV2Application_JSONMappingParameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
