package cmdutils

import (
	"errors"
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// NewKubernetesNetworkConfigLoader creates a new loader for the Kubernetes network config.
func NewKubernetesNetworkConfigLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	l.flagsIncompatibleWithConfigFile.Insert(
		"enable-elastic-load-balancing",
		"cluster",
	)

	l.validateWithConfigFile = func() error {
		if cmd.NameArg != "" {
			return fmt.Errorf("config file and cluster name argument %s", IncompatibleFlags)
		}
		knc := l.ClusterConfig.KubernetesNetworkConfig
		if knc == nil || knc.ElasticLoadBalancing == nil || knc.ElasticLoadBalancing.Enabled == nil {
			return errors.New("field kubernetesNetworkConfig.elasticLoadBalancing.enabled is required")
		}
		return api.ValidateElasticLoadBalancing(l.ClusterConfig)
	}

	l.validateWithoutConfigFile = func() error {
		if !cmd.CobraCommand.Flag("enable-elastic-load-balancing").Changed {
			return errors.New("--enable-elastic-load-balancing is required when a config file is not specified")
		}
		return nil
	}
	return l
}
