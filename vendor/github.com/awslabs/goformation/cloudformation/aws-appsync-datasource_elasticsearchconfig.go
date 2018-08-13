package cloudformation

import (
	"encoding/json"
)

// AWSAppSyncDataSource_ElasticsearchConfig AWS CloudFormation Resource (AWS::AppSync::DataSource.ElasticsearchConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-elasticsearchconfig.html
type AWSAppSyncDataSource_ElasticsearchConfig struct {

	// AwsRegion AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-elasticsearchconfig.html#cfn-appsync-datasource-elasticsearchconfig-awsregion
	AwsRegion *Value `json:"AwsRegion,omitempty"`

	// Endpoint AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-elasticsearchconfig.html#cfn-appsync-datasource-elasticsearchconfig-endpoint
	Endpoint *Value `json:"Endpoint,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppSyncDataSource_ElasticsearchConfig) AWSCloudFormationType() string {
	return "AWS::AppSync::DataSource.ElasticsearchConfig"
}

func (r *AWSAppSyncDataSource_ElasticsearchConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
