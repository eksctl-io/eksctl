package utils

import (
	"fmt"
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

func installCoreDNSCmd(g *cmdutils.Grouping) *cobra.Command {
	cfg := api.NewClusterConfig()
	cp := cmdutils.NewCommonParams(cfg)

	cp.Command = &cobra.Command{
		Use:   "install-coredns",
		Short: "Installs latest version of CoreDNS add-on into clusters, replacing kube-dns",
		Run: func(_ *cobra.Command, args []string) {
			cp.NameArg = cmdutils.GetNameArg(args)
			if err := doInstallCoreDNS(cp); err != nil {
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

func doInstallCoreDNS(cp *cmdutils.CommonParams) error {
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

	switch ctl.ControlPlaneVersion() {
	case "":
		return fmt.Errorf("unable to get control plane version")
	case api.Version1_10:
		return fmt.Errorf("%q is not supported on 1.10 cluster, run 'eksctl update cluster' first", defaultaddons.CoreDNS)
	}

	rawClient, err := ctl.NewRawClient(cfg)
	if err != nil {
		return err
	}

	waitTimeout := ctl.Provider.WaitTimeout()

	updateRequired, err := defaultaddons.InstallCoreDNS(rawClient, meta.Region, &waitTimeout, cp.Plan)
	if err != nil {
		return err
	}

	cmdutils.LogPlanModeWarning(cp.Plan && updateRequired)

	return nil
}
