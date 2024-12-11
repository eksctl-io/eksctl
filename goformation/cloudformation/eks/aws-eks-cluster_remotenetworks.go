package eks

import "goformation/v4/cloudformation/types"

type Cluster_RemoteNetworkConfig struct {
	RemotePodNetworks  []RemoteNetworks `json:"RemotePodNetworks,omitempty"`
	RemoteNodeNetworks []RemoteNetworks `json:"RemoteNodeNetworks,omitempty"`
}

type RemoteNetworks struct {
	CIDRs *types.Value `json:"Cidrs"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *Cluster_RemoteNetworkConfig) AWSCloudFormationType() string {
	return "AWS::EKS::Cluster.RemoteNetworkConfig"
}
