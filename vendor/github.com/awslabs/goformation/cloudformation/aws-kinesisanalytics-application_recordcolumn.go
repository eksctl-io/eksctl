package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsApplication_RecordColumn AWS CloudFormation Resource (AWS::KinesisAnalytics::Application.RecordColumn)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-application-recordcolumn.html
type AWSKinesisAnalyticsApplication_RecordColumn struct {

	// Mapping AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-application-recordcolumn.html#cfn-kinesisanalytics-application-recordcolumn-mapping
	Mapping *Value `json:"Mapping,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-application-recordcolumn.html#cfn-kinesisanalytics-application-recordcolumn-name
	Name *Value `json:"Name,omitempty"`

	// SqlType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-application-recordcolumn.html#cfn-kinesisanalytics-application-recordcolumn-sqltype
	SqlType *Value `json:"SqlType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsApplication_RecordColumn) AWSCloudFormationType() string {
	return "AWS::KinesisAnalytics::Application.RecordColumn"
}

func (r *AWSKinesisAnalyticsApplication_RecordColumn) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
