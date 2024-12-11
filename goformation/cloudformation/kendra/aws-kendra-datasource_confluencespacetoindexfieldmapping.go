package kendra

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataSource_ConfluenceSpaceToIndexFieldMapping AWS CloudFormation Resource (AWS::Kendra::DataSource.ConfluenceSpaceToIndexFieldMapping)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluencespacetoindexfieldmapping.html
type DataSource_ConfluenceSpaceToIndexFieldMapping struct {

	// DataSourceFieldName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluencespacetoindexfieldmapping.html#cfn-kendra-datasource-confluencespacetoindexfieldmapping-datasourcefieldname
	DataSourceFieldName *types.Value `json:"DataSourceFieldName,omitempty"`

	// DateFieldFormat AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluencespacetoindexfieldmapping.html#cfn-kendra-datasource-confluencespacetoindexfieldmapping-datefieldformat
	DateFieldFormat *types.Value `json:"DateFieldFormat,omitempty"`

	// IndexFieldName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluencespacetoindexfieldmapping.html#cfn-kendra-datasource-confluencespacetoindexfieldmapping-indexfieldname
	IndexFieldName *types.Value `json:"IndexFieldName,omitempty"`

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
func (r *DataSource_ConfluenceSpaceToIndexFieldMapping) AWSCloudFormationType() string {
	return "AWS::Kendra::DataSource.ConfluenceSpaceToIndexFieldMapping"
}
