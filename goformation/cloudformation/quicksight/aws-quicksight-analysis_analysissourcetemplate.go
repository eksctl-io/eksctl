package quicksight

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Analysis_AnalysisSourceTemplate AWS CloudFormation Resource (AWS::QuickSight::Analysis.AnalysisSourceTemplate)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-analysis-analysissourcetemplate.html
type Analysis_AnalysisSourceTemplate struct {

	// Arn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-analysis-analysissourcetemplate.html#cfn-quicksight-analysis-analysissourcetemplate-arn
	Arn *types.Value `json:"Arn,omitempty"`

	// DataSetReferences AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-analysis-analysissourcetemplate.html#cfn-quicksight-analysis-analysissourcetemplate-datasetreferences
	DataSetReferences []Analysis_DataSetReference `json:"DataSetReferences,omitempty"`

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
func (r *Analysis_AnalysisSourceTemplate) AWSCloudFormationType() string {
	return "AWS::QuickSight::Analysis.AnalysisSourceTemplate"
}
