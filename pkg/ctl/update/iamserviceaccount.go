package update

import (
	"errors"

	"github.com/weaveworks/eksctl/pkg/actions/iam"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func updateIAMServiceAccountCmd(cmd *cmdutils.Cmd) {
	updateIAMServiceAccountCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd) error {
		return doUpdateIAMServiceAccount(cmd)
	})
}

func updateIAMServiceAccountCmdWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd) error) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	serviceAccount := &api.ClusterIAMServiceAccount{}

	cfg.IAM.WithOIDC = api.Enabled()
	cfg.IAM.ServiceAccounts = append(cfg.IAM.ServiceAccounts, serviceAccount)

	cmd.SetDescription("iamserviceaccount", "Update an iamserviceaccount", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return runFunc(cmd)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "name of the EKS cluster")

		fs.StringVar(&serviceAccount.Name, "name", "", "name of the iamserviceaccount to update")
		fs.StringVar(&serviceAccount.Namespace, "namespace", "default", "namespace where to update the iamserviceaccount")
		fs.StringSliceVar(&serviceAccount.AttachPolicyARNs, "attach-policy-arn", []string{}, "ARN of the policy where to update the iamserviceaccount")

		cmdutils.AddIAMServiceAccountFilterFlags(fs, &cmd.Include, &cmd.Exclude)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, true)
}

func doUpdateIAMServiceAccount(cmd *cmdutils.Cmd) error {
	saFilter := filter.NewIAMServiceAccountFilter()

	if err := cmdutils.NewCreateIAMServiceAccountLoader(cmd, saFilter).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	meta := cmd.ClusterConfig.Metadata

	printer := printers.NewJSONPrinter()

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	cmdutils.LogRegionAndVersionInfo(meta)

	if err := ctl.CheckAuth(); err != nil {
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
		logger.Warning("no IAM OIDC provider associated with cluster, try 'eksctl utils associate-iam-oidc-provider --region=%s --cluster=%s'", meta.Region, meta.Name)
		return errors.New("unable to update iamserviceaccount(s) without IAM OIDC provider enabled")
	}
	stackManager := ctl.NewStackManager(cfg)

	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	return iam.New(cfg.Metadata.Name, stackManager, oidc, clientSet).UpdateIAMServiceAccounts(cfg.IAM.ServiceAccounts, cmd.Plan)
}
