package databrew

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Dataset_DatabaseInputDefinition AWS CloudFormation Resource (AWS::DataBrew::Dataset.DatabaseInputDefinition)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-dataset-databaseinputdefinition.html
type Dataset_DatabaseInputDefinition struct {

	// DatabaseTableName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-dataset-databaseinputdefinition.html#cfn-databrew-dataset-databaseinputdefinition-databasetablename
	DatabaseTableName *types.Value `json:"DatabaseTableName,omitempty"`

	// GlueConnectionName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-dataset-databaseinputdefinition.html#cfn-databrew-dataset-databaseinputdefinition-glueconnectionname
	GlueConnectionName *types.Value `json:"GlueConnectionName,omitempty"`

	// QueryString AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-dataset-databaseinputdefinition.html#cfn-databrew-dataset-databaseinputdefinition-querystring
	QueryString *types.Value `json:"QueryString,omitempty"`

	// TempDirectory AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-dataset-databaseinputdefinition.html#cfn-databrew-dataset-databaseinputdefinition-tempdirectory
	TempDirectory *Dataset_S3Location `json:"TempDirectory,omitempty"`

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
func (r *Dataset_DatabaseInputDefinition) AWSCloudFormationType() string {
	return "AWS::DataBrew::Dataset.DatabaseInputDefinition"
}
