package cmdutils

import (
	"errors"

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
}

// NewUpdateClusterVPCLoader will load config or use flags for 'eksctl utils update-cluster-vpc-config'.
func NewUpdateClusterVPCLoader(cmd *Cmd, options UpdateClusterVPCOptions) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	supportedOptions := []string{
		"private-access",
		"public-access",
		"public-access-cidrs",
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
			return errors.New("at least one of --public-access, --private-access and --public-access-cidrs must be specified")
		}
		clusterConfig := cmd.ClusterConfig

		if clusterConfig.VPC.ClusterEndpoints == nil {
			clusterConfig.VPC.ClusterEndpoints = api.ClusterEndpointAccessDefaults()
		}
		if flag := l.CobraCommand.Flag("private-access"); flag != nil && flag.Changed {
			clusterConfig.VPC.ClusterEndpoints.PrivateAccess = &options.PrivateAccess
		} else {
			clusterConfig.VPC.ClusterEndpoints.PrivateAccess = nil
		}

		if flag := l.CobraCommand.Flag("public-access"); flag != nil && flag.Changed {
			clusterConfig.VPC.ClusterEndpoints.PublicAccess = &options.PublicAccess
		} else {
			clusterConfig.VPC.ClusterEndpoints.PublicAccess = nil
		}
		clusterConfig.VPC.PublicAccessCIDRs = options.PublicAccessCIDRs
		return nil
	}

	l.validateWithConfigFile = func() error {
		logger.Info("only changes to vpc.clusterEndpoints and vpc.publicAccessCIDRs are updated in the EKS API, changes to any other fields will be ignored")
		if l.ClusterConfig.VPC == nil {
			l.ClusterConfig.VPC = api.NewClusterVPC(false)
		}
		api.SetClusterEndpointAccessDefaults(l.ClusterConfig.VPC)
		return nil
	}

	return l
}
