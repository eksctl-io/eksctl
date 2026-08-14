package eks

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// Cluster_KubeControllerManagerConfig AWS CloudFormation Resource (AWS::EKS::Cluster.KubeControllerManagerConfig)
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-properties-eks-cluster-kubecontrollermanagerconfig.html
type Cluster_KubeControllerManagerConfig struct {

	// HorizontalPodAutoscalerControllerConfig AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-properties-eks-cluster-kubecontrollermanagerconfig.html#cfn-eks-cluster-kubecontrollermanagerconfig-horizontalpodautoscalercontrollerconfig
	HorizontalPodAutoscalerControllerConfig *Cluster_HorizontalPodAutoscalerControllerConfig `json:"HorizontalPodAutoscalerControllerConfig,omitempty"`

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
func (r *Cluster_KubeControllerManagerConfig) AWSCloudFormationType() string {
	return "AWS::EKS::Cluster.KubeControllerManagerConfig"
}
