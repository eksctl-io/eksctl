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

	KopsClusterNameForVPC string
	Subnets               map[api.SubnetTopology]*[]string
	WithoutNodeGroup      bool
	Fargate               bool
	DryRun                bool
	CreateNGOptions
	CreateManagedNGOptions
}

// CreateManagedNGOptions holds options for creating a managed nodegroup
type CreateManagedNGOptions struct {
	Managed       bool
	Spot          bool
	InstanceTypes []string
}

// CreateNGOptions holds options for creating a nodegroup
type CreateNGOptions struct {
	InstallNeuronDevicePlugin bool
	InstallNvidiaDevicePlugin bool
	DryRun                    bool
}
