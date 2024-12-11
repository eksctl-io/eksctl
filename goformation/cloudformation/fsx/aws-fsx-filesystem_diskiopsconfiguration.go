package fsx

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// FileSystem_DiskIopsConfiguration AWS CloudFormation Resource (AWS::FSx::FileSystem.DiskIopsConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-openzfsconfiguration-diskiopsconfiguration.html
type FileSystem_DiskIopsConfiguration struct {

	// Iops AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-openzfsconfiguration-diskiopsconfiguration.html#cfn-fsx-filesystem-openzfsconfiguration-diskiopsconfiguration-iops
	Iops *types.Value `json:"Iops,omitempty"`

	// Mode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-openzfsconfiguration-diskiopsconfiguration.html#cfn-fsx-filesystem-openzfsconfiguration-diskiopsconfiguration-mode
	Mode *types.Value `json:"Mode,omitempty"`

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
func (r *FileSystem_DiskIopsConfiguration) AWSCloudFormationType() string {
	return "AWS::FSx::FileSystem.DiskIopsConfiguration"
}
