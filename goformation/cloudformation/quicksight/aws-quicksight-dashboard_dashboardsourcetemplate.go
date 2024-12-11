package quicksight

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Dashboard_DashboardSourceTemplate AWS CloudFormation Resource (AWS::QuickSight::Dashboard.DashboardSourceTemplate)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-dashboardsourcetemplate.html
type Dashboard_DashboardSourceTemplate struct {

	// Arn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-dashboardsourcetemplate.html#cfn-quicksight-dashboard-dashboardsourcetemplate-arn
	Arn *types.Value `json:"Arn,omitempty"`

	// DataSetReferences AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-dashboardsourcetemplate.html#cfn-quicksight-dashboard-dashboardsourcetemplate-datasetreferences
	DataSetReferences []Dashboard_DataSetReference `json:"DataSetReferences,omitempty"`

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
func (r *Dashboard_DashboardSourceTemplate) AWSCloudFormationType() string {
	return "AWS::QuickSight::Dashboard.DashboardSourceTemplate"
}
