package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// LaunchTemplate_Ebs AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.Ebs)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-ebs.html
type LaunchTemplate_Ebs struct {

	// DeleteOnTermination AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-ebs.html#cfn-ec2-launchtemplate-ebs-deleteontermination
	DeleteOnTermination *types.Value `json:"DeleteOnTermination,omitempty"`

	// Encrypted AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-ebs.html#cfn-ec2-launchtemplate-ebs-encrypted
	Encrypted *types.Value `json:"Encrypted,omitempty"`

	// Iops AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-ebs.html#cfn-ec2-launchtemplate-ebs-iops
	Iops *types.Value `json:"Iops,omitempty"`

	// KmsKeyId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-ebs.html#cfn-ec2-launchtemplate-ebs-kmskeyid
	KmsKeyId *types.Value `json:"KmsKeyId,omitempty"`

	// SnapshotId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-ebs.html#cfn-ec2-launchtemplate-ebs-snapshotid
	SnapshotId *types.Value `json:"SnapshotId,omitempty"`

	// Throughput AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-ebs.html#cfn-ec2-launchtemplate-ebs-throughput
	Throughput *types.Value `json:"Throughput,omitempty"`

	// VolumeSize AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-ebs.html#cfn-ec2-launchtemplate-ebs-volumesize
	VolumeSize *types.Value `json:"VolumeSize,omitempty"`

	// VolumeType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-ebs.html#cfn-ec2-launchtemplate-ebs-volumetype
	VolumeType *types.Value `json:"VolumeType,omitempty"`

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
func (r *LaunchTemplate_Ebs) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.Ebs"
}
