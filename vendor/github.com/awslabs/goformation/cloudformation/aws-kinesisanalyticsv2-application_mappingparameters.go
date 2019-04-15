package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_MappingParameters AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.MappingParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-mappingparameters.html
type AWSKinesisAnalyticsV2Application_MappingParameters struct {

	// CSVMappingParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-mappingparameters.html#cfn-kinesisanalyticsv2-application-mappingparameters-csvmappingparameters
	CSVMappingParameters *AWSKinesisAnalyticsV2Application_CSVMappingParameters `json:"CSVMappingParameters,omitempty"`

	// JSONMappingParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-mappingparameters.html#cfn-kinesisanalyticsv2-application-mappingparameters-jsonmappingparameters
	JSONMappingParameters *AWSKinesisAnalyticsV2Application_JSONMappingParameters `json:"JSONMappingParameters,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_MappingParameters) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.MappingParameters"
}

func (r *AWSKinesisAnalyticsV2Application_MappingParameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
