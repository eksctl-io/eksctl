package kendra

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataSource_DatabaseConfiguration AWS CloudFormation Resource (AWS::Kendra::DataSource.DatabaseConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-databaseconfiguration.html
type DataSource_DatabaseConfiguration struct {

	// AclConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-databaseconfiguration.html#cfn-kendra-datasource-databaseconfiguration-aclconfiguration
	AclConfiguration *DataSource_AclConfiguration `json:"AclConfiguration,omitempty"`

	// ColumnConfiguration AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-databaseconfiguration.html#cfn-kendra-datasource-databaseconfiguration-columnconfiguration
	ColumnConfiguration *DataSource_ColumnConfiguration `json:"ColumnConfiguration,omitempty"`

	// ConnectionConfiguration AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-databaseconfiguration.html#cfn-kendra-datasource-databaseconfiguration-connectionconfiguration
	ConnectionConfiguration *DataSource_ConnectionConfiguration `json:"ConnectionConfiguration,omitempty"`

	// DatabaseEngineType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-databaseconfiguration.html#cfn-kendra-datasource-databaseconfiguration-databaseenginetype
	DatabaseEngineType *types.Value `json:"DatabaseEngineType,omitempty"`

	// SqlConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-databaseconfiguration.html#cfn-kendra-datasource-databaseconfiguration-sqlconfiguration
	SqlConfiguration *DataSource_SqlConfiguration `json:"SqlConfiguration,omitempty"`

	// VpcConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-databaseconfiguration.html#cfn-kendra-datasource-databaseconfiguration-vpcconfiguration
	VpcConfiguration *DataSource_DataSourceVpcConfiguration `json:"VpcConfiguration,omitempty"`

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
func (r *DataSource_DatabaseConfiguration) AWSCloudFormationType() string {
	return "AWS::Kendra::DataSource.DatabaseConfiguration"
}
