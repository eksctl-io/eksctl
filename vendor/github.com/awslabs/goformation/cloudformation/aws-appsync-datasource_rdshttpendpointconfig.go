package cloudformation

import (
	"encoding/json"
)

// AWSAppSyncDataSource_RdsHttpEndpointConfig AWS CloudFormation Resource (AWS::AppSync::DataSource.RdsHttpEndpointConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-rdshttpendpointconfig.html
type AWSAppSyncDataSource_RdsHttpEndpointConfig struct {

	// AwsRegion AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-rdshttpendpointconfig.html#cfn-appsync-datasource-rdshttpendpointconfig-awsregion
	AwsRegion *Value `json:"AwsRegion,omitempty"`

	// AwsSecretStoreArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-rdshttpendpointconfig.html#cfn-appsync-datasource-rdshttpendpointconfig-awssecretstorearn
	AwsSecretStoreArn *Value `json:"AwsSecretStoreArn,omitempty"`

	// DatabaseName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-rdshttpendpointconfig.html#cfn-appsync-datasource-rdshttpendpointconfig-databasename
	DatabaseName *Value `json:"DatabaseName,omitempty"`

	// DbClusterIdentifier AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-rdshttpendpointconfig.html#cfn-appsync-datasource-rdshttpendpointconfig-dbclusteridentifier
	DbClusterIdentifier *Value `json:"DbClusterIdentifier,omitempty"`

	// Schema AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-rdshttpendpointconfig.html#cfn-appsync-datasource-rdshttpendpointconfig-schema
	Schema *Value `json:"Schema,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppSyncDataSource_RdsHttpEndpointConfig) AWSCloudFormationType() string {
	return "AWS::AppSync::DataSource.RdsHttpEndpointConfig"
}

func (r *AWSAppSyncDataSource_RdsHttpEndpointConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
