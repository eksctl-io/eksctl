package cloudformation

import (
	"encoding/json"
)

// AWSAppSyncDataSource_AwsIamConfig AWS CloudFormation Resource (AWS::AppSync::DataSource.AwsIamConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-awsiamconfig.html
type AWSAppSyncDataSource_AwsIamConfig struct {

	// SigningRegion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-awsiamconfig.html#cfn-appsync-datasource-awsiamconfig-signingregion
	SigningRegion *Value `json:"SigningRegion,omitempty"`

	// SigningServiceName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-awsiamconfig.html#cfn-appsync-datasource-awsiamconfig-signingservicename
	SigningServiceName *Value `json:"SigningServiceName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppSyncDataSource_AwsIamConfig) AWSCloudFormationType() string {
	return "AWS::AppSync::DataSource.AwsIamConfig"
}

func (r *AWSAppSyncDataSource_AwsIamConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
