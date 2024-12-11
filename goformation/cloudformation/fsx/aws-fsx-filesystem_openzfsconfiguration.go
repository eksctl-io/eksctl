package fsx

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// FileSystem_OpenZFSConfiguration AWS CloudFormation Resource (AWS::FSx::FileSystem.OpenZFSConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-openzfsconfiguration.html
type FileSystem_OpenZFSConfiguration struct {

	// AutomaticBackupRetentionDays AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-openzfsconfiguration.html#cfn-fsx-filesystem-openzfsconfiguration-automaticbackupretentiondays
	AutomaticBackupRetentionDays *types.Value `json:"AutomaticBackupRetentionDays,omitempty"`

	// CopyTagsToBackups AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-openzfsconfiguration.html#cfn-fsx-filesystem-openzfsconfiguration-copytagstobackups
	CopyTagsToBackups *types.Value `json:"CopyTagsToBackups,omitempty"`

	// CopyTagsToVolumes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-openzfsconfiguration.html#cfn-fsx-filesystem-openzfsconfiguration-copytagstovolumes
	CopyTagsToVolumes *types.Value `json:"CopyTagsToVolumes,omitempty"`

	// DailyAutomaticBackupStartTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-openzfsconfiguration.html#cfn-fsx-filesystem-openzfsconfiguration-dailyautomaticbackupstarttime
	DailyAutomaticBackupStartTime *types.Value `json:"DailyAutomaticBackupStartTime,omitempty"`

	// DeploymentType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-openzfsconfiguration.html#cfn-fsx-filesystem-openzfsconfiguration-deploymenttype
	DeploymentType *types.Value `json:"DeploymentType,omitempty"`

	// DiskIopsConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-openzfsconfiguration.html#cfn-fsx-filesystem-openzfsconfiguration-diskiopsconfiguration
	DiskIopsConfiguration *FileSystem_DiskIopsConfiguration `json:"DiskIopsConfiguration,omitempty"`

	// RootVolumeConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-openzfsconfiguration.html#cfn-fsx-filesystem-openzfsconfiguration-rootvolumeconfiguration
	RootVolumeConfiguration *FileSystem_RootVolumeConfiguration `json:"RootVolumeConfiguration,omitempty"`

	// ThroughputCapacity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-openzfsconfiguration.html#cfn-fsx-filesystem-openzfsconfiguration-throughputcapacity
	ThroughputCapacity *types.Value `json:"ThroughputCapacity,omitempty"`

	// WeeklyMaintenanceStartTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-openzfsconfiguration.html#cfn-fsx-filesystem-openzfsconfiguration-weeklymaintenancestarttime
	WeeklyMaintenanceStartTime *types.Value `json:"WeeklyMaintenanceStartTime,omitempty"`

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
func (r *FileSystem_OpenZFSConfiguration) AWSCloudFormationType() string {
	return "AWS::FSx::FileSystem.OpenZFSConfiguration"
}
