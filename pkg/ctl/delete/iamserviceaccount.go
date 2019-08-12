package delete

import (
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func deleteIAMServiceAccountCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	serviceAccount := &api.ClusterIAMServiceAccount{}

	cfg.IAM.WithOIDC = api.Enabled()
	cfg.IAM.ServiceAccounts = append(cfg.IAM.ServiceAccounts, serviceAccount)

	var onlyMissing bool

	cmd.SetDescription("iamserviceaccount", "Create an iamserviceaccount - AWS IAM role bound to a Kubernetes service account", "")

	cmd.SetRunFunc(func() error {
		return doDeleteIAMServiceAccount(cmd, serviceAccount, onlyMissing)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "name of the EKS cluster to delete the iamserviceaccount from")

		fs.StringVar(&serviceAccount.Name, "name", "", "name of the iamserviceaccount to delete")
		fs.StringVar(&serviceAccount.Namespace, "namespace", "default", "namespace where to delete the iamserviceaccount")

		cmdutils.AddIAMServiceAccountFilterFlags(fs, &cmd.Include, &cmd.Exclude)
		fs.BoolVar(&onlyMissing, "only-missing", false, "Only delete nodegroups that are not defined in the given config file")
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)

		cmd.Wait = false
		cmdutils.AddWaitFlag(fs, &cmd.Wait, "deletion of all resources")
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, true)
}

func doDeleteIAMServiceAccount(cmd *cmdutils.Cmd, serviceAccount *api.ClusterIAMServiceAccount, onlyMissing bool) error {
	saFilter := cmdutils.NewIAMServiceAccountFilter()

	if err := cmdutils.NewDeleteIAMServiceAccountLoader(cmd, serviceAccount, saFilter).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	meta := cmd.ClusterConfig.Metadata

	printer := printers.NewJSONPrinter()

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	logger.Info("using region %s", meta.Region)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if err := ctl.RefreshClusterConfig(cfg); err != nil {
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
		return fmt.Errorf("unable to delete iamserviceaccount(s) without IAM OIDC provider enabled")
	}

	stackManager := ctl.NewStackManager(cfg)

	if cmd.ClusterConfigFile != "" {
		logger.Info("comparing %d iamserviceaccounts defined in the given config (%q) against remote state", len(cfg.IAM.ServiceAccounts), cmd.ClusterConfigFile)
		if err := saFilter.SetIncludeOrExcludeMissingFilter(stackManager, onlyMissing, &cfg.IAM.ServiceAccounts); err != nil {
			return err
		}
	}

	saSubset, _ := saFilter.MatchAll(cfg.IAM.ServiceAccounts)
	saFilter.LogInfo(cfg.IAM.ServiceAccounts)

	tasks, err := stackManager.NewTasksToDeleteIAMServiceAccounts(saSubset, oidc, kubernetes.NewCachedClientSet(clientSet), cmd.Wait)
	if err != nil {
		return err
	}
	tasks.PlanMode = cmd.Plan

	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	logger.Info(tasks.Describe())
	if errs := tasks.DoAllSync(); len(errs) > 0 {
		logger.Info("%d error(s) occurred and IAM Role stacks haven't been deleted properly, you may wish to check CloudFormation console", len(errs))
		for _, err := range errs {
			logger.Critical("%s\n", err.Error())
		}
		return fmt.Errorf("failed to delete iamserviceaccount(s)")
	}

	cmdutils.LogPlanModeWarning(cmd.Plan && saSubset.Len() > 0)

	return nil
}
