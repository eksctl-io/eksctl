package quicksight

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Dashboard_DashboardSourceEntity AWS CloudFormation Resource (AWS::QuickSight::Dashboard.DashboardSourceEntity)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-dashboardsourceentity.html
type Dashboard_DashboardSourceEntity struct {

	// SourceTemplate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-dashboardsourceentity.html#cfn-quicksight-dashboard-dashboardsourceentity-sourcetemplate
	SourceTemplate *Dashboard_DashboardSourceTemplate `json:"SourceTemplate,omitempty"`

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
func (r *Dashboard_DashboardSourceEntity) AWSCloudFormationType() string {
	return "AWS::QuickSight::Dashboard.DashboardSourceEntity"
}
