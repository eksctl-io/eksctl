package cloudformation

import (
	"encoding/json"
)

// AWSFSxFileSystem_WindowsConfiguration AWS CloudFormation Resource (AWS::FSx::FileSystem.WindowsConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-windowsconfiguration.html
type AWSFSxFileSystem_WindowsConfiguration struct {

	// ActiveDirectoryId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-windowsconfiguration.html#cfn-fsx-filesystem-windowsconfiguration-activedirectoryid
	ActiveDirectoryId *Value `json:"ActiveDirectoryId,omitempty"`

	// AutomaticBackupRetentionDays AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-windowsconfiguration.html#cfn-fsx-filesystem-windowsconfiguration-automaticbackupretentiondays
	AutomaticBackupRetentionDays *Value `json:"AutomaticBackupRetentionDays,omitempty"`

	// CopyTagsToBackups AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-windowsconfiguration.html#cfn-fsx-filesystem-windowsconfiguration-copytagstobackups
	CopyTagsToBackups *Value `json:"CopyTagsToBackups,omitempty"`

	// DailyAutomaticBackupStartTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-windowsconfiguration.html#cfn-fsx-filesystem-windowsconfiguration-dailyautomaticbackupstarttime
	DailyAutomaticBackupStartTime *Value `json:"DailyAutomaticBackupStartTime,omitempty"`

	// ThroughputCapacity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-windowsconfiguration.html#cfn-fsx-filesystem-windowsconfiguration-throughputcapacity
	ThroughputCapacity *Value `json:"ThroughputCapacity,omitempty"`

	// WeeklyMaintenanceStartTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-windowsconfiguration.html#cfn-fsx-filesystem-windowsconfiguration-weeklymaintenancestarttime
	WeeklyMaintenanceStartTime *Value `json:"WeeklyMaintenanceStartTime,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSFSxFileSystem_WindowsConfiguration) AWSCloudFormationType() string {
	return "AWS::FSx::FileSystem.WindowsConfiguration"
}

func (r *AWSFSxFileSystem_WindowsConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
