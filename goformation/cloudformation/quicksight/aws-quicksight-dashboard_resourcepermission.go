package quicksight

import (
	"github.com/awslabs/goformation/v4/cloudformation/types"

	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Dashboard_ResourcePermission AWS CloudFormation Resource (AWS::QuickSight::Dashboard.ResourcePermission)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-resourcepermission.html
type Dashboard_ResourcePermission struct {

	// Actions AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-resourcepermission.html#cfn-quicksight-dashboard-resourcepermission-actions
	Actions *types.Value `json:"Actions,omitempty"`

	// Principal AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-resourcepermission.html#cfn-quicksight-dashboard-resourcepermission-principal
	Principal *types.Value `json:"Principal,omitempty"`

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
func (r *Dashboard_ResourcePermission) AWSCloudFormationType() string {
	return "AWS::QuickSight::Dashboard.ResourcePermission"
}
