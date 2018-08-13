package cloudformation

import (
	"encoding/json"
)

// AWSAppSyncDataSource_DynamoDBConfig AWS CloudFormation Resource (AWS::AppSync::DataSource.DynamoDBConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-dynamodbconfig.html
type AWSAppSyncDataSource_DynamoDBConfig struct {

	// AwsRegion AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-dynamodbconfig.html#cfn-appsync-datasource-dynamodbconfig-awsregion
	AwsRegion *Value `json:"AwsRegion,omitempty"`

	// TableName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-dynamodbconfig.html#cfn-appsync-datasource-dynamodbconfig-tablename
	TableName *Value `json:"TableName,omitempty"`

	// UseCallerCredentials AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-dynamodbconfig.html#cfn-appsync-datasource-dynamodbconfig-usecallercredentials
	UseCallerCredentials *Value `json:"UseCallerCredentials,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppSyncDataSource_DynamoDBConfig) AWSCloudFormationType() string {
	return "AWS::AppSync::DataSource.DynamoDBConfig"
}

func (r *AWSAppSyncDataSource_DynamoDBConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
