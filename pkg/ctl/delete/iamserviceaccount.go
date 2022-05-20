package delete

import (
	"context"
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/irsa"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func deleteIAMServiceAccountCmd(cmd *cmdutils.Cmd) {
	deleteIAMServiceAccountCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd, serviceAccount *api.ClusterIAMServiceAccount, onlyMissing bool) error {
		return doDeleteIAMServiceAccount(cmd, serviceAccount, onlyMissing)
	})
}

func deleteIAMServiceAccountCmdWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd, serviceAccount *api.ClusterIAMServiceAccount, onlyMissing bool) error) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	serviceAccount := &api.ClusterIAMServiceAccount{}

	cfg.IAM.WithOIDC = api.Enabled()
	cfg.IAM.ServiceAccounts = append(cfg.IAM.ServiceAccounts, serviceAccount)

	var onlyMissing bool

	cmd.SetDescription("iamserviceaccount", "Delete an IAM service account", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return runFunc(cmd, serviceAccount, onlyMissing)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "name of the EKS cluster to delete the iamserviceaccount from")

		fs.StringVar(&serviceAccount.Name, "name", "", "name of the iamserviceaccount to delete")
		fs.StringVar(&serviceAccount.Namespace, "namespace", "default", "namespace where to delete the iamserviceaccount")

		cmdutils.AddIAMServiceAccountFilterFlags(fs, &cmd.Include, &cmd.Exclude)
		fs.BoolVar(&onlyMissing, "only-missing", false, "Only delete iamserviceaccounts that are not defined in the given config file")
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)

		cmd.Wait = false
		cmdutils.AddWaitFlag(fs, &cmd.Wait, "deletion of all resources")
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, true)
}

func doDeleteIAMServiceAccount(cmd *cmdutils.Cmd, serviceAccount *api.ClusterIAMServiceAccount, onlyMissing bool) error {
	saFilter := filter.NewIAMServiceAccountFilter()

	if err := cmdutils.NewDeleteIAMServiceAccountLoader(cmd, serviceAccount, saFilter).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig

	printer := printers.NewJSONPrinter()

	ctx := context.Background()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
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

	oidc, err := ctl.NewOpenIDConnectManager(ctx, cfg)
	if err != nil {
		return err
	}

	providerExists, err := oidc.CheckProviderExists(ctx)
	if err != nil {
		return err
	}

	if !providerExists {
		return fmt.Errorf("unable to delete iamserviceaccount(s) without IAM OIDC provider enabled")
	}

	stackManager := ctl.NewStackManager(cfg)

	if cmd.ClusterConfigFile != "" {
		logger.Info("comparing %d iamserviceaccounts defined in the given config (%q) against remote state", len(cfg.IAM.ServiceAccounts), cmd.ClusterConfigFile)
		if err := saFilter.SetDeleteFilter(ctx, stackManager, onlyMissing, cfg); err != nil {
			return err
		}
	}

	saFilter.LogInfo(cfg.IAM.ServiceAccounts)

	saSubset, _ := saFilter.MatchAll(cfg.IAM.ServiceAccounts)

	irsaManager := irsa.New(cfg.Metadata.Name, stackManager, oidc, clientSet)

	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}
	return irsaManager.Delete(ctx, saSubset.List(), cmd.Plan, cmd.Wait)
}
