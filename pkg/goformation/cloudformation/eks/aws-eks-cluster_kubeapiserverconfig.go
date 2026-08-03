package eks

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"
)

// Cluster_KubeApiServerConfig AWS CloudFormation Resource (AWS::EKS::Cluster.KubeApiServerConfig)
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-properties-eks-cluster-kubeapiserverconfig.html
type Cluster_KubeApiServerConfig struct {

	// EventTtl AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-properties-eks-cluster-kubeapiserverconfig.html#cfn-eks-cluster-kubeapiserverconfig-eventttl
	EventTtl *types.Value `json:"EventTtl,omitempty"`

	// ServiceNodePortRange AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-properties-eks-cluster-kubeapiserverconfig.html#cfn-eks-cluster-kubeapiserverconfig-servicenodeportrange
	ServiceNodePortRange *Cluster_ServiceNodePortRange `json:"ServiceNodePortRange,omitempty"`

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
func (r *Cluster_KubeApiServerConfig) AWSCloudFormationType() string {
	return "AWS::EKS::Cluster.KubeApiServerConfig"
}
