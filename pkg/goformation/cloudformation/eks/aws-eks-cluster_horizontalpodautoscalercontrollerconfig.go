package eks

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"
)

// Cluster_HorizontalPodAutoscalerControllerConfig AWS CloudFormation Resource (AWS::EKS::Cluster.HorizontalPodAutoscalerControllerConfig)
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-properties-eks-cluster-horizontalpodautoscalercontrollerconfig.html
type Cluster_HorizontalPodAutoscalerControllerConfig struct {

	// HorizontalPodAutoscalerSyncPeriod AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-properties-eks-cluster-horizontalpodautoscalercontrollerconfig.html#cfn-eks-cluster-horizontalpodautoscalercontrollerconfig-horizontalpodautoscalersyncperiod
	HorizontalPodAutoscalerSyncPeriod *types.Value `json:"HorizontalPodAutoscalerSyncPeriod,omitempty"`

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
func (r *Cluster_HorizontalPodAutoscalerControllerConfig) AWSCloudFormationType() string {
	return "AWS::EKS::Cluster.HorizontalPodAutoscalerControllerConfig"
}
