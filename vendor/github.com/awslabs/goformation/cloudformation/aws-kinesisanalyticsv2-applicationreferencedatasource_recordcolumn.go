package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2ApplicationReferenceDataSource_RecordColumn AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource.RecordColumn)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-recordcolumn.html
type AWSKinesisAnalyticsV2ApplicationReferenceDataSource_RecordColumn struct {

	// Mapping AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-recordcolumn.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-recordcolumn-mapping
	Mapping *Value `json:"Mapping,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-recordcolumn.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-recordcolumn-name
	Name *Value `json:"Name,omitempty"`

	// SqlType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-recordcolumn.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-recordcolumn-sqltype
	SqlType *Value `json:"SqlType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_RecordColumn) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource.RecordColumn"
}

func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_RecordColumn) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
