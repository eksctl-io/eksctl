package create

import (
	"errors"

	"github.com/weaveworks/eksctl/pkg/actions/irsa"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func createIAMServiceAccountCmd(cmd *cmdutils.Cmd) {
	createIAMServiceAccountCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd, overrideExistingServiceAccounts bool) error {
		return doCreateIAMServiceAccount(cmd, overrideExistingServiceAccounts)
	})
}

func createIAMServiceAccountCmdWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd, overrideExistingServiceAccounts bool) error) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	serviceAccount := &api.ClusterIAMServiceAccount{
		RoleOnly: api.Disabled(),
	}

	cfg.IAM.WithOIDC = api.Enabled()
	cfg.IAM.ServiceAccounts = append(cfg.IAM.ServiceAccounts, serviceAccount)

	var overrideExistingServiceAccounts bool

	cmd.SetDescription("iamserviceaccount", "Create an iamserviceaccount - AWS IAM role bound to a Kubernetes service account", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return runFunc(cmd, overrideExistingServiceAccounts)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "name of the EKS cluster to add the iamserviceaccount to")

		fs.StringVar(&serviceAccount.Name, "name", "", "name of the iamserviceaccount to create")
		fs.StringVar(&serviceAccount.Namespace, "namespace", "default", "namespace where to create the iamserviceaccount")
		fs.StringSliceVar(&serviceAccount.AttachPolicyARNs, "attach-policy-arn", []string{}, "ARN of the policy where to create the iamserviceaccount")
		fs.StringVar(&serviceAccount.AttachRoleARN, "attach-role-arn", "", "ARN of the role to attach to the iamserviceaccount")
		fs.StringVar(&serviceAccount.RoleName, "role-name", "", "Set a custom name for the created role")
		fs.BoolVar(serviceAccount.RoleOnly, "role-only", false, "disable service account creation, only the role will be created")

		cmdutils.AddStringToStringVarPFlag(fs, &serviceAccount.Tags, "tags", "", map[string]string{}, "Used to tag the IAM role")

		fs.BoolVar(&overrideExistingServiceAccounts, "override-existing-serviceaccounts", false, "create IAM roles for existing serviceaccounts and update the serviceaccount")

		cmdutils.AddIAMServiceAccountFilterFlags(fs, &cmd.Include, &cmd.Exclude)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, true)
}

func doCreateIAMServiceAccount(cmd *cmdutils.Cmd, overrideExistingServiceAccounts bool) error {
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
		return errors.New("unable to create iamserviceaccount(s) without IAM OIDC provider enabled")
	}
	stackManager := ctl.NewStackManager(cfg)

	if err := saFilter.SetExcludeExistingFilter(stackManager, clientSet, cfg.IAM.ServiceAccounts, overrideExistingServiceAccounts); err != nil {
		return err
	}

	filteredServiceAccounts := saFilter.FilterMatching(cfg.IAM.ServiceAccounts)
	saFilter.LogInfo(cfg.IAM.ServiceAccounts)
	if !overrideExistingServiceAccounts {
		logger.Warning("serviceaccounts that exists in Kubernetes will be excluded, use --override-existing-serviceaccounts to override")
	} else {
		logger.Warning("metadata of serviceaccounts that exist in Kubernetes will be updated, as --override-existing-serviceaccounts was set")
	}

	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	return irsa.New(cfg.Metadata.Name, stackManager, oidc, clientSet).CreateIAMServiceAccount(filteredServiceAccounts, cmd.Plan)
}
