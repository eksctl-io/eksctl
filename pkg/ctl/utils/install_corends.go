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
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()

	cmd := &cobra.Command{
		Use:   "install-coredns",
		Short: "Installs latest version of CoreDNS add-on into clusters, replacing kube-dns",
		Run: func(cmd *cobra.Command, args []string) {
			if err := doInstallCoreDNS(p, cfg, cmdutils.GetNameArg(args), cmd); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cmd)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, p)
		cmdutils.AddConfigFileFlag(&clusterConfigFile, fs)
		cmdutils.AddApproveFlag(&plan, cmd, fs)
	})

	cmdutils.AddCommonFlagsForAWS(group, p, false)

	group.AddTo(cmd)

	return cmd
}

func doInstallCoreDNS(p *api.ProviderConfig, cfg *api.ClusterConfig, nameArg string, cmd *cobra.Command) error {
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

	updateRequired, err := defaultaddons.InstallCoreDNS(rawClient, meta.Region, &waitTimeout, plan)
	if err != nil {
		return err
	}

	cmdutils.LogPlanModeWarning(plan && updateRequired)

	return nil
}
