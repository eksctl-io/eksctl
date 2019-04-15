package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_InputSchema AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.InputSchema)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-inputschema.html
type AWSKinesisAnalyticsV2Application_InputSchema struct {

	// RecordColumns AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-inputschema.html#cfn-kinesisanalyticsv2-application-inputschema-recordcolumns
	RecordColumns []AWSKinesisAnalyticsV2Application_RecordColumn `json:"RecordColumns,omitempty"`

	// RecordEncoding AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-inputschema.html#cfn-kinesisanalyticsv2-application-inputschema-recordencoding
	RecordEncoding *Value `json:"RecordEncoding,omitempty"`

	// RecordFormat AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-inputschema.html#cfn-kinesisanalyticsv2-application-inputschema-recordformat
	RecordFormat *AWSKinesisAnalyticsV2Application_RecordFormat `json:"RecordFormat,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_InputSchema) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.InputSchema"
}

func (r *AWSKinesisAnalyticsV2Application_InputSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
