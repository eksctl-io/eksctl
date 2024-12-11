package glue

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Table_SchemaReference AWS CloudFormation Resource (AWS::Glue::Table.SchemaReference)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-schemareference.html
type Table_SchemaReference struct {

	// SchemaId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-schemareference.html#cfn-glue-table-schemareference-schemaid
	SchemaId *Table_SchemaId `json:"SchemaId,omitempty"`

	// SchemaVersionId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-schemareference.html#cfn-glue-table-schemareference-schemaversionid
	SchemaVersionId *types.Value `json:"SchemaVersionId,omitempty"`

	// SchemaVersionNumber AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-schemareference.html#cfn-glue-table-schemareference-schemaversionnumber
	SchemaVersionNumber *types.Value `json:"SchemaVersionNumber,omitempty"`

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
func (r *Table_SchemaReference) AWSCloudFormationType() string {
	return "AWS::Glue::Table.SchemaReference"
}
