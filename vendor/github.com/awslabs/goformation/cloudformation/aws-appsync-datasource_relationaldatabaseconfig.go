package cloudformation

import (
	"encoding/json"
)

// AWSAppSyncDataSource_RelationalDatabaseConfig AWS CloudFormation Resource (AWS::AppSync::DataSource.RelationalDatabaseConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-relationaldatabaseconfig.html
type AWSAppSyncDataSource_RelationalDatabaseConfig struct {

	// RdsHttpEndpointConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-relationaldatabaseconfig.html#cfn-appsync-datasource-relationaldatabaseconfig-rdshttpendpointconfig
	RdsHttpEndpointConfig *AWSAppSyncDataSource_RdsHttpEndpointConfig `json:"RdsHttpEndpointConfig,omitempty"`

	// RelationalDatabaseSourceType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-relationaldatabaseconfig.html#cfn-appsync-datasource-relationaldatabaseconfig-relationaldatabasesourcetype
	RelationalDatabaseSourceType *Value `json:"RelationalDatabaseSourceType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppSyncDataSource_RelationalDatabaseConfig) AWSCloudFormationType() string {
	return "AWS::AppSync::DataSource.RelationalDatabaseConfig"
}

func (r *AWSAppSyncDataSource_RelationalDatabaseConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
