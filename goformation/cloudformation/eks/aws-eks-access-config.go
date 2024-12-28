package eks

import (
	"goformation/v4/cloudformation/types"
)

// Cluster_AccessConfig describes the access configuration for the cluster.
type Cluster_AccessConfig struct {

	// AuthenticationMode specifies the desired authentication mode for the cluster.
	AuthenticationMode *types.Value

	// BootstrapClusterCreatorAdminPermissions specifies whether the cluster creator IAM principal was set as a cluster
	// admin access entry during cluster creation time.
	BootstrapClusterCreatorAdminPermissions *types.Value
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *Cluster_AccessConfig) AWSCloudFormationType() string {
	return "AWS::EKS::Cluster.AccessConfig"
}
