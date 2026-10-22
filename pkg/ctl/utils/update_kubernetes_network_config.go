package utils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func updateKubernetesNetworkConfig(cmd *cmdutils.Cmd, handler func(*cmdutils.Cmd) error) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("update-kubernetes-network-config", "update Kubernetes network config", "update the Kubernetes network config on a cluster")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		if err := cmdutils.NewKubernetesNetworkConfigLoader(cmd).Load(); err != nil {
			return err
		}
		return handler(cmd)
	}

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
	})
}

func updateKubernetesNetworkConfigCmd(cmd *cmdutils.Cmd) {
	updateKubernetesNetworkConfig(cmd, doUpdateKubernetesNetworkConfig)
}

func doUpdateKubernetesNetworkConfig(cmd *cmdutils.Cmd) error {
	cfg := cmd.ClusterConfig
	ctx := context.Background()
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}
	// capture the requested state first, as NewProviderForExistingCluster overwrites kubernetesNetworkConfig
	elb := cfg.ElasticLoadBalancingRequest()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}
	cluster := ctl.Status.ClusterInfo.Cluster
	if autoModeEnabled := cluster.ComputeConfig != nil && aws.ToBool(cluster.ComputeConfig.Enabled); aws.ToBool(elb.Enabled) != autoModeEnabled {
		return fmt.Errorf("kubernetesNetworkConfig.elasticLoadBalancing.enabled must match Auto Mode on the cluster (Auto Mode enabled: %t); "+
			"use `eksctl update auto-mode-config` to enable or disable Auto Mode", autoModeEnabled)
	}
	if knc := cluster.KubernetesNetworkConfig; knc != nil && api.ElasticLoadBalancingUpToDate(knc.ElasticLoadBalancing, elb) {
		logger.Info("the load balancing capability is already up-to-date (%s)", api.DescribeElasticLoadBalancing(knc.ElasticLoadBalancing))
		return nil
	}
	if err := ctl.UpdateClusterConfig(ctx, &eks.UpdateClusterConfigInput{
		Name: aws.String(cfg.Metadata.Name),
		KubernetesNetworkConfig: &ekstypes.KubernetesNetworkConfigRequest{
			ElasticLoadBalancing: elb,
		},
	}); err != nil {
		return fmt.Errorf("updating Kubernetes network config: %w", err)
	}
	logger.Info("the load balancing capability was updated successfully (%s)", api.DescribeElasticLoadBalancing(elb))
	return nil
}
