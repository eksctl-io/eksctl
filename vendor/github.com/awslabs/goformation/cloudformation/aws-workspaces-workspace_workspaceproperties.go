package cloudformation

import (
	"encoding/json"
)

// AWSWorkSpacesWorkspace_WorkspaceProperties AWS CloudFormation Resource (AWS::WorkSpaces::Workspace.WorkspaceProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-workspaces-workspace-workspaceproperties.html
type AWSWorkSpacesWorkspace_WorkspaceProperties struct {

	// ComputeTypeName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-workspaces-workspace-workspaceproperties.html#cfn-workspaces-workspace-workspaceproperties-computetypename
	ComputeTypeName *Value `json:"ComputeTypeName,omitempty"`

	// RootVolumeSizeGib AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-workspaces-workspace-workspaceproperties.html#cfn-workspaces-workspace-workspaceproperties-rootvolumesizegib
	RootVolumeSizeGib *Value `json:"RootVolumeSizeGib,omitempty"`

	// RunningMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-workspaces-workspace-workspaceproperties.html#cfn-workspaces-workspace-workspaceproperties-runningmode
	RunningMode *Value `json:"RunningMode,omitempty"`

	// RunningModeAutoStopTimeoutInMinutes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-workspaces-workspace-workspaceproperties.html#cfn-workspaces-workspace-workspaceproperties-runningmodeautostoptimeoutinminutes
	RunningModeAutoStopTimeoutInMinutes *Value `json:"RunningModeAutoStopTimeoutInMinutes,omitempty"`

	// UserVolumeSizeGib AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-workspaces-workspace-workspaceproperties.html#cfn-workspaces-workspace-workspaceproperties-uservolumesizegib
	UserVolumeSizeGib *Value `json:"UserVolumeSizeGib,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSWorkSpacesWorkspace_WorkspaceProperties) AWSCloudFormationType() string {
	return "AWS::WorkSpaces::Workspace.WorkspaceProperties"
}

func (r *AWSWorkSpacesWorkspace_WorkspaceProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
