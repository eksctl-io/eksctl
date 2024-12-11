package workspaces

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Workspace_WorkspaceProperties AWS CloudFormation Resource (AWS::WorkSpaces::Workspace.WorkspaceProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-workspaces-workspace-workspaceproperties.html
type Workspace_WorkspaceProperties struct {

	// ComputeTypeName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-workspaces-workspace-workspaceproperties.html#cfn-workspaces-workspace-workspaceproperties-computetypename
	ComputeTypeName *types.Value `json:"ComputeTypeName,omitempty"`

	// RootVolumeSizeGib AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-workspaces-workspace-workspaceproperties.html#cfn-workspaces-workspace-workspaceproperties-rootvolumesizegib
	RootVolumeSizeGib *types.Value `json:"RootVolumeSizeGib,omitempty"`

	// RunningMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-workspaces-workspace-workspaceproperties.html#cfn-workspaces-workspace-workspaceproperties-runningmode
	RunningMode *types.Value `json:"RunningMode,omitempty"`

	// RunningModeAutoStopTimeoutInMinutes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-workspaces-workspace-workspaceproperties.html#cfn-workspaces-workspace-workspaceproperties-runningmodeautostoptimeoutinminutes
	RunningModeAutoStopTimeoutInMinutes *types.Value `json:"RunningModeAutoStopTimeoutInMinutes,omitempty"`

	// UserVolumeSizeGib AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-workspaces-workspace-workspaceproperties.html#cfn-workspaces-workspace-workspaceproperties-uservolumesizegib
	UserVolumeSizeGib *types.Value `json:"UserVolumeSizeGib,omitempty"`

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
func (r *Workspace_WorkspaceProperties) AWSCloudFormationType() string {
	return "AWS::WorkSpaces::Workspace.WorkspaceProperties"
}
