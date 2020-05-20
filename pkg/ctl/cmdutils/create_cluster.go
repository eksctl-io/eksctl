package cmdutils

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// CreateClusterCmdParams groups CLI options for the create cluster command.
type CreateClusterCmdParams struct {
	WriteKubeconfig             bool
	KubeconfigPath              string
	AutoKubeconfigPath          bool
	AuthenticatorRoleARN        string
	SetContext                  bool
	AvailabilityZones           []string
	InstallWindowsVPCController bool
	InstallNeuronDevicePlugin   bool
	KopsClusterNameForVPC       string
	Subnets                     map[api.SubnetTopology]*[]string
	WithoutNodeGroup            bool
	Managed                     bool
	Fargate                     bool
}
