package eks

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// Cluster_RemoteNetworkConfig AWS CloudFormation Resource (AWS::EKS::Cluster.RemoteNetworkConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-cluster-remotenetworkconfig.html
type Cluster_RemoteNetworkConfig struct {

	// RemoteNodeNetworks AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-cluster-remotenetworkconfig.html#cfn-eks-cluster-remotenetworkconfig-remotenodenetworks
	RemoteNodeNetworks []Cluster_RemoteNodeNetwork `json:"RemoteNodeNetworks,omitempty"`

	// RemotePodNetworks AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-cluster-remotenetworkconfig.html#cfn-eks-cluster-remotenetworkconfig-remotepodnetworks
	RemotePodNetworks []Cluster_RemotePodNetwork `json:"RemotePodNetworks,omitempty"`

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
func (r *Cluster_RemoteNetworkConfig) AWSCloudFormationType() string {
	return "AWS::EKS::Cluster.RemoteNetworkConfig"
}
