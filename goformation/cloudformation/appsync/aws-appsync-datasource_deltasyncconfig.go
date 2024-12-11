package appsync

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataSource_DeltaSyncConfig AWS CloudFormation Resource (AWS::AppSync::DataSource.DeltaSyncConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-deltasyncconfig.html
type DataSource_DeltaSyncConfig struct {

	// BaseTableTTL AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-deltasyncconfig.html#cfn-appsync-datasource-deltasyncconfig-basetablettl
	BaseTableTTL *types.Value `json:"BaseTableTTL,omitempty"`

	// DeltaSyncTableName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-deltasyncconfig.html#cfn-appsync-datasource-deltasyncconfig-deltasynctablename
	DeltaSyncTableName *types.Value `json:"DeltaSyncTableName,omitempty"`

	// DeltaSyncTableTTL AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-deltasyncconfig.html#cfn-appsync-datasource-deltasyncconfig-deltasynctablettl
	DeltaSyncTableTTL *types.Value `json:"DeltaSyncTableTTL,omitempty"`

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
func (r *DataSource_DeltaSyncConfig) AWSCloudFormationType() string {
	return "AWS::AppSync::DataSource.DeltaSyncConfig"
}
