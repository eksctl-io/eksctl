package quicksight

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataSet_LogicalTable AWS CloudFormation Resource (AWS::QuickSight::DataSet.LogicalTable)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dataset-logicaltable.html
type DataSet_LogicalTable struct {

	// Alias AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dataset-logicaltable.html#cfn-quicksight-dataset-logicaltable-alias
	Alias *types.Value `json:"Alias,omitempty"`

	// DataTransforms AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dataset-logicaltable.html#cfn-quicksight-dataset-logicaltable-datatransforms
	DataTransforms []DataSet_TransformOperation `json:"DataTransforms,omitempty"`

	// Source AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dataset-logicaltable.html#cfn-quicksight-dataset-logicaltable-source
	Source *DataSet_LogicalTableSource `json:"Source,omitempty"`

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
func (r *DataSet_LogicalTable) AWSCloudFormationType() string {
	return "AWS::QuickSight::DataSet.LogicalTable"
}
