package apply

import (
	"fmt"

	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/weaveworks/eksctl/pkg/actions/apply"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func Command(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("apply", "EXPERIMENTAL: Reconcile cluster config", "EXPERIMENTAL: Reconcile the cluster config with"+
		"the deloyed cluster configuration. Supported resource:\n"+
		"  - IAMServiceAccounts")
	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)

	cmd.CobraCommand.RunE = func(_ *cobra.Command, _ []string) error {
		return doApply(cmd)
	}
}

func doApply(cmd *cmdutils.Cmd) error {
	if cmd.ClusterConfigFile == "" {
		return fmt.Errorf("please provide the cluster config file")
	}

	err := newConfigLoader(cmd).Load()
	if err != nil {
		return err
	}
	cfg := cmd.ClusterConfig

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	oidc, err := ctl.NewOpenIDConnectManager(cfg)
	if err != nil {
		return err
	}

	providerExists, err := oidc.CheckProviderExists()
	if err != nil {
		return err
	}
	if !providerExists {
		return fmt.Errorf("oidc must be enabled to create IAMServiceAccounts")
	}

	stackManager := manager.NewStackCollection(ctl.Provider, cfg)

	logger.Warning("EXPERIMENTAL: eksctl apply is an experimental command, the currently supported resources are: IAMServiceAccounts")

	err = apply.New(cfg, ctl, stackManager, oidc, clientSet, cmd.Plan).Reconcile()
	if err != nil {
		return err
	}

	if cmd.Plan {
		logger.Warning("no changes were applied, run again with '--approve' to apply the changes")

	}
	return nil
}

func newConfigLoader(cmd *cmdutils.Cmd) cmdutils.ClusterConfigLoader {
	l := cmdutils.NewConfigLoaderBuilder()
	return l.Build(cmd)
}
