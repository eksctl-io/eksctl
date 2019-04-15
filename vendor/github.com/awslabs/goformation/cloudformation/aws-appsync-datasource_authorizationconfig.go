package cloudformation

import (
	"encoding/json"
)

// AWSAppSyncDataSource_AuthorizationConfig AWS CloudFormation Resource (AWS::AppSync::DataSource.AuthorizationConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-authorizationconfig.html
type AWSAppSyncDataSource_AuthorizationConfig struct {

	// AuthorizationType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-authorizationconfig.html#cfn-appsync-datasource-authorizationconfig-authorizationtype
	AuthorizationType *Value `json:"AuthorizationType,omitempty"`

	// AwsIamConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-authorizationconfig.html#cfn-appsync-datasource-authorizationconfig-awsiamconfig
	AwsIamConfig *AWSAppSyncDataSource_AwsIamConfig `json:"AwsIamConfig,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppSyncDataSource_AuthorizationConfig) AWSCloudFormationType() string {
	return "AWS::AppSync::DataSource.AuthorizationConfig"
}

func (r *AWSAppSyncDataSource_AuthorizationConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
