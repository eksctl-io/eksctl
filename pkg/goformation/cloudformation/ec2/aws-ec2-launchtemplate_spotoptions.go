package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// LaunchTemplate_SpotOptions AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.SpotOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-spotoptions.html
type LaunchTemplate_SpotOptions struct {

	// BlockDurationMinutes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-spotoptions.html#cfn-ec2-launchtemplate-spotoptions-blockdurationminutes
	BlockDurationMinutes *types.Value `json:"BlockDurationMinutes,omitempty"`

	// InstanceInterruptionBehavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-spotoptions.html#cfn-ec2-launchtemplate-spotoptions-instanceinterruptionbehavior
	InstanceInterruptionBehavior *types.Value `json:"InstanceInterruptionBehavior,omitempty"`

	// MaxPrice AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-spotoptions.html#cfn-ec2-launchtemplate-spotoptions-maxprice
	MaxPrice *types.Value `json:"MaxPrice,omitempty"`

	// SpotInstanceType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-spotoptions.html#cfn-ec2-launchtemplate-spotoptions-spotinstancetype
	SpotInstanceType *types.Value `json:"SpotInstanceType,omitempty"`

	// ValidUntil AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-spotoptions.html#cfn-ec2-launchtemplate-spotoptions-validuntil
	ValidUntil *types.Value `json:"ValidUntil,omitempty"`

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
func (r *LaunchTemplate_SpotOptions) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.SpotOptions"
}
