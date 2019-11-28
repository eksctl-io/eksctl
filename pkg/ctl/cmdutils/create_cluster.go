package cmdutils

import (
	"github.com/pkg/errors"
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
	KopsClusterNameForVPC       string
	Subnets                     map[api.SubnetTopology]*[]string
	WithoutNodeGroup            bool
	Managed                     bool
	Fargate                     bool
}

// Validate validates this CreateClusterCmdParams.
func (p CreateClusterCmdParams) Validate() error {
	if p.Managed && p.Fargate {
		return errors.New("--managed and --fargate are mutually exclusive: please provide either one of these flags, but not both")
	}
	return nil
}
