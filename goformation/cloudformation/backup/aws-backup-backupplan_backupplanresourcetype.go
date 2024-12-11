package backup

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// BackupPlan_BackupPlanResourceType AWS CloudFormation Resource (AWS::Backup::BackupPlan.BackupPlanResourceType)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-backup-backupplan-backupplanresourcetype.html
type BackupPlan_BackupPlanResourceType struct {

	// AdvancedBackupSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-backup-backupplan-backupplanresourcetype.html#cfn-backup-backupplan-backupplanresourcetype-advancedbackupsettings
	AdvancedBackupSettings []BackupPlan_AdvancedBackupSettingResourceType `json:"AdvancedBackupSettings,omitempty"`

	// BackupPlanName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-backup-backupplan-backupplanresourcetype.html#cfn-backup-backupplan-backupplanresourcetype-backupplanname
	BackupPlanName *types.Value `json:"BackupPlanName,omitempty"`

	// BackupPlanRule AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-backup-backupplan-backupplanresourcetype.html#cfn-backup-backupplan-backupplanresourcetype-backupplanrule
	BackupPlanRule []BackupPlan_BackupRuleResourceType `json:"BackupPlanRule,omitempty"`

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
func (r *BackupPlan_BackupPlanResourceType) AWSCloudFormationType() string {
	return "AWS::Backup::BackupPlan.BackupPlanResourceType"
}
