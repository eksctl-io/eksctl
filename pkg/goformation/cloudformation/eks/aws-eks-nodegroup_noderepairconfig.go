package eks

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// Nodegroup_NodeRepairConfig AWS CloudFormation Resource (AWS::EKS::Nodegroup.NodeRepairConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-nodegroup-noderepairconfig.html
type Nodegroup_NodeRepairConfig struct {

	// Enabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-nodegroup-noderepairconfig.html#cfn-eks-nodegroup-noderepairconfig-enabled
	Enabled *types.Value `json:"Enabled,omitempty"`

	// MaxUnhealthyNodeThresholdPercentage AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-nodegroup-noderepairconfig.html#cfn-eks-nodegroup-noderepairconfig-maxunhealthynodethresholdpercentage
	MaxUnhealthyNodeThresholdPercentage *types.Value `json:"MaxUnhealthyNodeThresholdPercentage,omitempty"`

	// MaxUnhealthyNodeThresholdCount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-nodegroup-noderepairconfig.html#cfn-eks-nodegroup-noderepairconfig-maxunhealthynodethresholdcount
	MaxUnhealthyNodeThresholdCount *types.Value `json:"MaxUnhealthyNodeThresholdCount,omitempty"`

	// MaxParallelNodesRepairedPercentage AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-nodegroup-noderepairconfig.html#cfn-eks-nodegroup-noderepairconfig-maxparallelnodesrepairedpercentage
	MaxParallelNodesRepairedPercentage *types.Value `json:"MaxParallelNodesRepairedPercentage,omitempty"`

	// MaxParallelNodesRepairedCount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-nodegroup-noderepairconfig.html#cfn-eks-nodegroup-noderepairconfig-maxparallelnodesrepairedcount
	MaxParallelNodesRepairedCount *types.Value `json:"MaxParallelNodesRepairedCount,omitempty"`

	// NodeRepairConfigOverrides AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-nodegroup-noderepairconfig.html#cfn-eks-nodegroup-noderepairconfig-noderepairconfigurations
	NodeRepairConfigOverrides []Nodegroup_NodeRepairConfigOverride `json:"NodeRepairConfigOverrides,omitempty"`

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
func (r *Nodegroup_NodeRepairConfig) AWSCloudFormationType() string {
	return "AWS::EKS::Nodegroup.NodeRepairConfig"
}

// Nodegroup_NodeRepairConfigOverride AWS CloudFormation Resource (AWS::EKS::Nodegroup.NodeRepairConfigOverride)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-nodegroup-noderepairconfigurations.html
type Nodegroup_NodeRepairConfigOverride struct {

	// NodeMonitoringCondition AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-nodegroup-noderepairconfigurations.html#cfn-eks-nodegroup-noderepairconfigurations-nodemonitoringcondition
	NodeMonitoringCondition *types.Value `json:"NodeMonitoringCondition,omitempty"`

	// NodeUnhealthyReason AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-nodegroup-noderepairconfigurations.html#cfn-eks-nodegroup-noderepairconfigurations-nodeunhealthyreason
	NodeUnhealthyReason *types.Value `json:"NodeUnhealthyReason,omitempty"`

	// MinRepairWaitTimeMins AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-nodegroup-noderepairconfigurations.html#cfn-eks-nodegroup-noderepairconfigurations-minrepairwaittimemins
	MinRepairWaitTimeMins *types.Value `json:"MinRepairWaitTimeMins,omitempty"`

	// RepairAction AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-nodegroup-noderepairconfigurations.html#cfn-eks-nodegroup-noderepairconfigurations-repairaction
	RepairAction *types.Value `json:"RepairAction,omitempty"`

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
func (r *Nodegroup_NodeRepairConfigOverride) AWSCloudFormationType() string {
	return "AWS::EKS::Nodegroup.NodeRepairConfigOverride"
}
