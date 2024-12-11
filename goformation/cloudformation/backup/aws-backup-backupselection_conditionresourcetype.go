package backup

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// BackupSelection_ConditionResourceType AWS CloudFormation Resource (AWS::Backup::BackupSelection.ConditionResourceType)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-backup-backupselection-conditionresourcetype.html
type BackupSelection_ConditionResourceType struct {

	// ConditionKey AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-backup-backupselection-conditionresourcetype.html#cfn-backup-backupselection-conditionresourcetype-conditionkey
	ConditionKey *types.Value `json:"ConditionKey,omitempty"`

	// ConditionType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-backup-backupselection-conditionresourcetype.html#cfn-backup-backupselection-conditionresourcetype-conditiontype
	ConditionType *types.Value `json:"ConditionType,omitempty"`

	// ConditionValue AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-backup-backupselection-conditionresourcetype.html#cfn-backup-backupselection-conditionresourcetype-conditionvalue
	ConditionValue *types.Value `json:"ConditionValue,omitempty"`

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
func (r *BackupSelection_ConditionResourceType) AWSCloudFormationType() string {
	return "AWS::Backup::BackupSelection.ConditionResourceType"
}
