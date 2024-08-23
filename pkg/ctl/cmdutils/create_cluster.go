package cmdutils

import (
	"io"

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

	ConfigReader io.Reader
}

// NodeGroupOptions holds options for creating nodegroups.
type NodeGroupOptions struct {
	CreateNGOptions
	CreateManagedNGOptions
	UpdateAuthConfigMap     *bool
	SkipOutdatedAddonsCheck bool
	SubnetIDs               []string
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
	NodeGroupParallelism      int
}
