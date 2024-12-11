package eks

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Cluster_KubernetesNetworkConfig AWS CloudFormation Resource (AWS::EKS::Cluster.KubernetesNetworkConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-cluster-kubernetesnetworkconfig.html
type Cluster_KubernetesNetworkConfig struct {

	// IpFamily AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-cluster-kubernetesnetworkconfig.html#cfn-eks-cluster-kubernetesnetworkconfig-ipfamily
	IpFamily *types.Value `json:"IpFamily,omitempty"`

	// ServiceIpv4Cidr AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-cluster-kubernetesnetworkconfig.html#cfn-eks-cluster-kubernetesnetworkconfig-serviceipv4cidr
	ServiceIpv4Cidr *types.Value `json:"ServiceIpv4Cidr,omitempty"`

	// ServiceIpv6Cidr AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-cluster-kubernetesnetworkconfig.html#cfn-eks-cluster-kubernetesnetworkconfig-serviceipv6cidr
	ServiceIpv6Cidr *types.Value `json:"ServiceIpv6Cidr,omitempty"`

	ElasticLoadBalancing *Cluster_ElasticLoadBalancing `json:"ElasticLoadBalancing,omitempty"`

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
func (r *Cluster_KubernetesNetworkConfig) AWSCloudFormationType() string {
	return "AWS::EKS::Cluster.KubernetesNetworkConfig"
}
