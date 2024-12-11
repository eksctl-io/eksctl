package batch

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// JobDefinition_EfsVolumeConfiguration AWS CloudFormation Resource (AWS::Batch::JobDefinition.EfsVolumeConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-efsvolumeconfiguration.html
type JobDefinition_EfsVolumeConfiguration struct {

	// AuthorizationConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-efsvolumeconfiguration.html#cfn-batch-jobdefinition-efsvolumeconfiguration-authorizationconfig
	AuthorizationConfig *JobDefinition_AuthorizationConfig `json:"AuthorizationConfig,omitempty"`

	// FileSystemId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-efsvolumeconfiguration.html#cfn-batch-jobdefinition-efsvolumeconfiguration-filesystemid
	FileSystemId *types.Value `json:"FileSystemId,omitempty"`

	// RootDirectory AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-efsvolumeconfiguration.html#cfn-batch-jobdefinition-efsvolumeconfiguration-rootdirectory
	RootDirectory *types.Value `json:"RootDirectory,omitempty"`

	// TransitEncryption AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-efsvolumeconfiguration.html#cfn-batch-jobdefinition-efsvolumeconfiguration-transitencryption
	TransitEncryption *types.Value `json:"TransitEncryption,omitempty"`

	// TransitEncryptionPort AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-efsvolumeconfiguration.html#cfn-batch-jobdefinition-efsvolumeconfiguration-transitencryptionport
	TransitEncryptionPort *types.Value `json:"TransitEncryptionPort,omitempty"`

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
func (r *JobDefinition_EfsVolumeConfiguration) AWSCloudFormationType() string {
	return "AWS::Batch::JobDefinition.EfsVolumeConfiguration"
}
