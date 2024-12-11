package appsync

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataSource_OpenSearchServiceConfig AWS CloudFormation Resource (AWS::AppSync::DataSource.OpenSearchServiceConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-opensearchserviceconfig.html
type DataSource_OpenSearchServiceConfig struct {

	// AwsRegion AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-opensearchserviceconfig.html#cfn-appsync-datasource-opensearchserviceconfig-awsregion
	AwsRegion *types.Value `json:"AwsRegion,omitempty"`

	// Endpoint AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-opensearchserviceconfig.html#cfn-appsync-datasource-opensearchserviceconfig-endpoint
	Endpoint *types.Value `json:"Endpoint,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *DataSource_OpenSearchServiceConfig) AWSCloudFormationType() string {
	return "AWS::AppSync::DataSource.OpenSearchServiceConfig"
}
