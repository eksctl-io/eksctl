package lightsail

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Instance_Disk AWS CloudFormation Resource (AWS::Lightsail::Instance.Disk)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-disk.html
type Instance_Disk struct {

	// AttachedTo AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-disk.html#cfn-lightsail-instance-disk-attachedto
	AttachedTo *types.Value `json:"AttachedTo,omitempty"`

	// AttachmentState AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-disk.html#cfn-lightsail-instance-disk-attachmentstate
	AttachmentState *types.Value `json:"AttachmentState,omitempty"`

	// DiskName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-disk.html#cfn-lightsail-instance-disk-diskname
	DiskName *types.Value `json:"DiskName,omitempty"`

	// IOPS AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-disk.html#cfn-lightsail-instance-disk-iops
	IOPS *types.Value `json:"IOPS,omitempty"`

	// IsSystemDisk AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-disk.html#cfn-lightsail-instance-disk-issystemdisk
	IsSystemDisk *types.Value `json:"IsSystemDisk,omitempty"`

	// Path AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-disk.html#cfn-lightsail-instance-disk-path
	Path *types.Value `json:"Path,omitempty"`

	// SizeInGb AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-disk.html#cfn-lightsail-instance-disk-sizeingb
	SizeInGb *types.Value `json:"SizeInGb,omitempty"`

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
func (r *Instance_Disk) AWSCloudFormationType() string {
	return "AWS::Lightsail::Instance.Disk"
}
