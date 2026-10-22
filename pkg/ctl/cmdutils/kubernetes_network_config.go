package cmdutils

import (
	"errors"
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// NewKubernetesNetworkConfigLoader creates a new loader for the Kubernetes network config.
func NewKubernetesNetworkConfigLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	l.flagsIncompatibleWithConfigFile.Insert("cluster")

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
		return errors.New("a config file is required; set kubernetesNetworkConfig.elasticLoadBalancing in it")
	}
	return l
}
