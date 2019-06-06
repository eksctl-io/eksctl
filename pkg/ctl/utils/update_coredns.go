package utils

import (
	"os"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	defaultaddons "github.com/weaveworks/eksctl/pkg/addons/default"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func updateCoreDNSCmd(g *cmdutils.Grouping) *cobra.Command {
	cfg := api.NewClusterConfig()
	cp := cmdutils.NewCommonParams(cfg)

	cp.Command = &cobra.Command{
		Use:   "update-coredns",
		Short: "Update coredns add-on to ensure image the standard Amazon EKS version",
		Run: func(_ *cobra.Command, args []string) {
			cp.NameArg = cmdutils.GetNameArg(args)
			if err := doUpdateCoreDNS(cp); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cp.Command)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cp.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cp.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cp)
	})

	cmdutils.AddCommonFlagsForAWS(group, cp.ProviderConfig, false)

	group.AddTo(cp.Command)
	return cp.Command
}

func doUpdateCoreDNS(cp *cmdutils.CommonParams) error {
	if err := cmdutils.NewMetadataLoader(cp).Load(); err != nil {
		return err
	}

	cfg := cp.ClusterConfig
	meta := cp.ClusterConfig.Metadata

	ctl := eks.New(cp.ProviderConfig, cfg)

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(cp.ProviderConfig)
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

	updateRequired, err := defaultaddons.UpdateCoreDNSImageTag(clientSet, cp.Plan)
	if err != nil {
		return err
	}

	cmdutils.LogPlanModeWarning(cp.Plan && updateRequired)

	return nil
}
