package utils

import (
	"context"
	"errors"
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

func updateControlPlaneComponentConfig(cmd *cmdutils.Cmd, handler func(*cmdutils.Cmd) error) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("update-control-plane-component-config", "update control plane component config",
		"update the kube-apiserver, kube-scheduler and kube-controller-manager configuration of a cluster")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		if err := cmdutils.NewControlPlaneComponentConfigLoader(cmd).Load(); err != nil {
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

func updateControlPlaneComponentConfigCmd(cmd *cmdutils.Cmd) {
	updateControlPlaneComponentConfig(cmd, doUpdateControlPlaneComponentConfig)
}

func doUpdateControlPlaneComponentConfig(cmd *cmdutils.Cmd) error {
	cfg := cmd.ClusterConfig
	ctx := context.Background()
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}
	kubeAPIServerConfig := makeKubeAPIServerConfigRequest(cfg.KubeAPIServerConfig)
	kubeSchedulerConfig := makeKubeSchedulerConfigRequest(cfg.KubeSchedulerConfig)
	kubeControllerManagerConfig := makeKubeControllerManagerConfigRequest(cfg.KubeControllerManagerConfig)
	if kubeAPIServerConfig == nil && kubeSchedulerConfig == nil && kubeControllerManagerConfig == nil {
		return errors.New("no control plane component config values are set; nothing to update")
	}

	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	input := &eks.UpdateClusterConfigInput{
		Name:                        aws.String(cfg.Metadata.Name),
		KubeApiServerConfig:         kubeAPIServerConfig,
		KubeSchedulerConfig:         kubeSchedulerConfig,
		KubeControllerManagerConfig: kubeControllerManagerConfig,
	}
	if err := ctl.UpdateClusterConfig(ctx, input); err != nil {
		return fmt.Errorf("updating control plane component config: %w", err)
	}
	logger.Info("control plane component config updated successfully")
	return nil
}

// makeKubeAPIServerConfigRequest converts the API kube-apiserver config to its update
// request representation. It returns nil if no values are set, so that the component is
// omitted from the request entirely and left unchanged.
// It mirrors makeKubeAPIServerConfig in pkg/cfn/builder/cluster.go, which does the same
// conversion for the CloudFormation template used on create.
func makeKubeAPIServerConfigRequest(config *api.KubeAPIServerConfig) *ekstypes.KubeApiServerConfigRequest {
	if config == nil {
		return nil
	}
	out := &ekstypes.KubeApiServerConfigRequest{}
	isSet := false
	if config.EventTTL != nil {
		out.EventTtl = config.EventTTL
		isSet = true
	}
	if portRange := config.ServiceNodePortRange; portRange != nil {
		nodePortRange := &ekstypes.ServiceNodePortRange{}
		if portRange.MinPort != nil {
			nodePortRange.MinPort = int32(*portRange.MinPort)
		}
		if portRange.MaxPort != nil {
			nodePortRange.MaxPort = int32(*portRange.MaxPort)
		}
		if portRange.MinPort != nil || portRange.MaxPort != nil {
			out.ServiceNodePortRange = nodePortRange
			isSet = true
		}
	}
	if !isSet {
		return nil
	}
	return out
}

// makeKubeSchedulerConfigRequest converts the API kube-scheduler config to its update
// request representation. It returns nil if no values are set, so that the component is
// omitted from the request entirely and left unchanged.
// It mirrors makeKubeSchedulerConfig in pkg/cfn/builder/cluster.go.
func makeKubeSchedulerConfigRequest(config *api.KubeSchedulerConfig) *ekstypes.KubeSchedulerConfigRequest {
	if config == nil || config.NodeResourcesFit == nil || config.NodeResourcesFit.ScoringStrategy == nil {
		return nil
	}
	strategy := config.NodeResourcesFit.ScoringStrategy
	scoringStrategy := &ekstypes.ScoringStrategy{}
	isSet := false
	if strategy.Type != nil {
		scoringStrategy.Type = ekstypes.ScoringStrategyType(*strategy.Type)
		isSet = true
	}
	for _, resource := range strategy.Resources {
		resourceWeight := ekstypes.ResourceWeight{}
		if resource.Name != nil {
			resourceWeight.Name = resource.Name
		}
		if resource.Weight != nil {
			resourceWeight.Weight = aws.Int32(int32(*resource.Weight))
		}
		scoringStrategy.Resources = append(scoringStrategy.Resources, resourceWeight)
		isSet = true
	}
	if !isSet {
		return nil
	}
	return &ekstypes.KubeSchedulerConfigRequest{
		NodeResourcesFit: &ekstypes.NodeResourcesFitConfig{
			ScoringStrategy: scoringStrategy,
		},
	}
}

// makeKubeControllerManagerConfigRequest converts the API kube-controller-manager config to
// its update request representation. It returns nil if no values are set, so that the
// component is omitted from the request entirely and left unchanged.
// It mirrors makeKubeControllerManagerConfig in pkg/cfn/builder/cluster.go.
func makeKubeControllerManagerConfigRequest(config *api.KubeControllerManagerConfig) *ekstypes.KubeControllerManagerConfigRequest {
	if config == nil {
		return nil
	}
	result := &ekstypes.KubeControllerManagerConfigRequest{}
	set := false
	if hpaConfig := config.HorizontalPodAutoscalerControllerConfig; hpaConfig != nil && hpaConfig.HorizontalPodAutoscalerSyncPeriod != nil {
		result.HorizontalPodAutoscalerControllerConfig = &ekstypes.HorizontalPodAutoscalerControllerConfigRequest{
			HorizontalPodAutoscalerSyncPeriod: hpaConfig.HorizontalPodAutoscalerSyncPeriod,
		}
		set = true
	}
	if podGCConfig := config.PodGCControllerConfig; podGCConfig != nil && podGCConfig.TerminatedPodGCThreshold != nil {
		result.PodGcControllerConfig = &ekstypes.PodGcControllerConfigRequest{
			TerminatedPodGcThreshold: aws.Int32(int32(*podGCConfig.TerminatedPodGCThreshold)),
		}
		set = true
	}
	if !set {
		return nil
	}
	return result
}
