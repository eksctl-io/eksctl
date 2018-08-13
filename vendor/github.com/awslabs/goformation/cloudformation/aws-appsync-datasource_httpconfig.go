package cloudformation

import (
	"encoding/json"
)

// AWSAppSyncDataSource_HttpConfig AWS CloudFormation Resource (AWS::AppSync::DataSource.HttpConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-httpconfig.html
type AWSAppSyncDataSource_HttpConfig struct {

	// Endpoint AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-httpconfig.html#cfn-appsync-datasource-httpconfig-endpoint
	Endpoint *Value `json:"Endpoint,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppSyncDataSource_HttpConfig) AWSCloudFormationType() string {
	return "AWS::AppSync::DataSource.HttpConfig"
}

func (r *AWSAppSyncDataSource_HttpConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
