package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsApplicationReferenceDataSource_RecordColumn AWS CloudFormation Resource (AWS::KinesisAnalytics::ApplicationReferenceDataSource.RecordColumn)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-applicationreferencedatasource-recordcolumn.html
type AWSKinesisAnalyticsApplicationReferenceDataSource_RecordColumn struct {

	// Mapping AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-applicationreferencedatasource-recordcolumn.html#cfn-kinesisanalytics-applicationreferencedatasource-recordcolumn-mapping
	Mapping *Value `json:"Mapping,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-applicationreferencedatasource-recordcolumn.html#cfn-kinesisanalytics-applicationreferencedatasource-recordcolumn-name
	Name *Value `json:"Name,omitempty"`

	// SqlType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-applicationreferencedatasource-recordcolumn.html#cfn-kinesisanalytics-applicationreferencedatasource-recordcolumn-sqltype
	SqlType *Value `json:"SqlType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsApplicationReferenceDataSource_RecordColumn) AWSCloudFormationType() string {
	return "AWS::KinesisAnalytics::ApplicationReferenceDataSource.RecordColumn"
}

func (r *AWSKinesisAnalyticsApplicationReferenceDataSource_RecordColumn) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
