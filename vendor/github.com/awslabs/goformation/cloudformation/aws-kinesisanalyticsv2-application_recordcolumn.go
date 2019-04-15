package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_RecordColumn AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.RecordColumn)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-recordcolumn.html
type AWSKinesisAnalyticsV2Application_RecordColumn struct {

	// Mapping AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-recordcolumn.html#cfn-kinesisanalyticsv2-application-recordcolumn-mapping
	Mapping *Value `json:"Mapping,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-recordcolumn.html#cfn-kinesisanalyticsv2-application-recordcolumn-name
	Name *Value `json:"Name,omitempty"`

	// SqlType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-recordcolumn.html#cfn-kinesisanalyticsv2-application-recordcolumn-sqltype
	SqlType *Value `json:"SqlType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_RecordColumn) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.RecordColumn"
}

func (r *AWSKinesisAnalyticsV2Application_RecordColumn) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
