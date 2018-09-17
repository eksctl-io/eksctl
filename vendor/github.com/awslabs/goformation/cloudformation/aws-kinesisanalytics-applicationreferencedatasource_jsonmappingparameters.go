package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsApplicationReferenceDataSource_JSONMappingParameters AWS CloudFormation Resource (AWS::KinesisAnalytics::ApplicationReferenceDataSource.JSONMappingParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-applicationreferencedatasource-jsonmappingparameters.html
type AWSKinesisAnalyticsApplicationReferenceDataSource_JSONMappingParameters struct {

	// RecordRowPath AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-applicationreferencedatasource-jsonmappingparameters.html#cfn-kinesisanalytics-applicationreferencedatasource-jsonmappingparameters-recordrowpath
	RecordRowPath *Value `json:"RecordRowPath,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsApplicationReferenceDataSource_JSONMappingParameters) AWSCloudFormationType() string {
	return "AWS::KinesisAnalytics::ApplicationReferenceDataSource.JSONMappingParameters"
}

func (r *AWSKinesisAnalyticsApplicationReferenceDataSource_JSONMappingParameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
