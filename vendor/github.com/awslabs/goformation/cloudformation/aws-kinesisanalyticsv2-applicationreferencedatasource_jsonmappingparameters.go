package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2ApplicationReferenceDataSource_JSONMappingParameters AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource.JSONMappingParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-jsonmappingparameters.html
type AWSKinesisAnalyticsV2ApplicationReferenceDataSource_JSONMappingParameters struct {

	// RecordRowPath AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationreferencedatasource-jsonmappingparameters.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-jsonmappingparameters-recordrowpath
	RecordRowPath *Value `json:"RecordRowPath,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_JSONMappingParameters) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource.JSONMappingParameters"
}

func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_JSONMappingParameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
