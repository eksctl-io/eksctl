package utils

import (
	"os"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	defaultaddons "github.com/weaveworks/eksctl/pkg/addons/default"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func updateKubeProxyCmd(g *cmdutils.Grouping) *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()

	cmd := &cobra.Command{
		Use:   "update-kube-proxy",
		Short: "Update kube-proxy add-on to ensure image matches Kubernetes control plane version",
		Run: func(cmd *cobra.Command, args []string) {
			if err := doUpdateKubeProxy(p, cfg, cmdutils.GetNameArg(args), cmd); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cmd)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, p)
		cmdutils.AddConfigFileFlag(&clusterConfigFile, fs)
	})

	cmdutils.AddCommonFlagsForAWS(group, p, false)

	group.AddTo(cmd)

	return cmd
}

func doUpdateKubeProxy(p *api.ProviderConfig, cfg *api.ClusterConfig, nameArg string, cmd *cobra.Command) error {
	if err := cmdutils.NewMetadataLoader(p, cfg, clusterConfigFile, nameArg, cmd).Load(); err != nil {
		return err
	}

	ctl := eks.New(p, cfg)
	meta := cfg.Metadata

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(p)
	}
	logger.Info("using region %s", meta.Region)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if err := ctl.GetCredentials(cfg); err != nil {
		return errors.Wrapf(err, "getting credentials for cluster %q", meta.Name)
	}

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	cprv, err := ctl.ControlPlaneReleaseVersion(clientSet)
	if err != nil {
		return err
	}

	return defaultaddons.UpdateKubeProxyImageTag(clientSet, cprv, false)
}
