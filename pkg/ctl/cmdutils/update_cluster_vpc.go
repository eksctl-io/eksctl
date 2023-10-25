package cmdutils

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// UpdateClusterVPCOptions holds the options for updating the VPC config.
type UpdateClusterVPCOptions struct {
	// PrivateAccess configures access for the private endpoint.
	PrivateAccess bool
	// PublicAccess configures access for the public endpoint.
	PublicAccess bool
	// PublicAccessCIDRs configures the public access CIDRs.
	PublicAccessCIDRs []string
	// ControlPlaneSubnetIDs configures the subnets for the control plane.
	ControlPlaneSubnetIDs []string
	// ControlPlaneSecurityGroupIDs configures the security group IDs for the control plane.
	ControlPlaneSecurityGroupIDs []string
}

// NewUpdateClusterVPCLoader will load config or use flags for 'eksctl utils update-cluster-vpc-config'.
func NewUpdateClusterVPCLoader(cmd *Cmd, options UpdateClusterVPCOptions) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	supportedOptions := []string{
		"private-access",
		"public-access",
		"public-access-cidrs",
		"control-plane-subnet-ids",
		"control-plane-security-group-ids",
	}

	l.flagsIncompatibleWithConfigFile.Insert(supportedOptions...)

	l.validateWithoutConfigFile = func() error {
		if err := l.validateMetadataWithoutConfigFile(); err != nil {
			return err
		}
		hasRequiredOptions := false
		for _, option := range supportedOptions {
			if flag := l.CobraCommand.Flag(option); flag != nil && flag.Changed {
				hasRequiredOptions = true
				break
			}
		}
		if !hasRequiredOptions {
			options := make([]string, 0, len(supportedOptions))
			for _, o := range supportedOptions {
				options = append(options, "--"+o)
			}
			return fmt.Errorf("at least one of these options must be specified: %s", strings.Join(options, ", "))
		}
		clusterConfig := cmd.ClusterConfig
		if flag := l.CobraCommand.Flag("private-access"); flag != nil && flag.Changed {
			clusterConfig.VPC.ClusterEndpoints = &api.ClusterEndpoints{
				PrivateAccess: &options.PrivateAccess,
			}
		}
		if flag := l.CobraCommand.Flag("public-access"); flag != nil && flag.Changed {
			if clusterConfig.VPC.ClusterEndpoints == nil {
				clusterConfig.VPC.ClusterEndpoints = &api.ClusterEndpoints{
					PublicAccess: &options.PublicAccess,
				}
			} else {
				clusterConfig.VPC.ClusterEndpoints.PublicAccess = &options.PublicAccess
			}
		}
		clusterConfig.VPC.PublicAccessCIDRs = options.PublicAccessCIDRs
		clusterConfig.VPC.ControlPlaneSubnetIDs = options.ControlPlaneSubnetIDs
		clusterConfig.VPC.ControlPlaneSecurityGroupIDs = options.ControlPlaneSecurityGroupIDs
		return nil
	}

	l.validateWithConfigFile = func() error {
		logger.Info("only changes to vpc.clusterEndpoints, vpc.publicAccessCIDRs, vpc.controlPlaneSubnetIDs and vpc.controlPlaneSecurityGroupIDs are updated in the EKS API, changes to any other fields will be ignored")
		if l.ClusterConfig.VPC == nil {
			l.ClusterConfig.VPC = api.NewClusterVPC(false)
		}
		api.SetClusterEndpointAccessDefaults(l.ClusterConfig.VPC)
		return nil
	}

	return l
}
