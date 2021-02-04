package update

import (
	"fmt"

	awseks "github.com/aws/aws-sdk-go/service/eks"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/actions/addon"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/logger"
)

func updateAddonCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"addon",
		"Upgrade an Addon",
		"",
	)

	var force bool
	cmd.ClusterConfig.Addons = []*api.Addon{{}}
	cmd.FlagSetGroup.InFlagSet("Addon", func(fs *pflag.FlagSet) {
		fs.StringVar(&cmd.ClusterConfig.Addons[0].Name, "name", "", "Addon name")
		fs.StringVar(&cmd.ClusterConfig.Addons[0].Version, "version", "", "Addon version")
		fs.StringVar(&cmd.ClusterConfig.Addons[0].ServiceAccountRoleARN, "service-account-role-arn", "", "Addon serviceAccountRoleARN")
		fs.BoolVar(&force, "force", false, "Force applies the add-on to overwrite an existing add-on")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return updateAddon(cmd, force)
	}
}

func updateAddon(cmd *cmdutils.Cmd, force bool) error {
	if err := cmdutils.NewCreateOrUpgradeAddonLoader(cmd).Load(); err != nil {
		return err
	}
	clusterProvider, err := cmd.NewCtl()
	if err != nil {
		return err
	}

	oidc, err := clusterProvider.NewOpenIDConnectManager(cmd.ClusterConfig)
	if err != nil {
		return err
	}

	oidcProviderExists, err := oidc.CheckProviderExists()
	if err != nil {
		return err
	}

	if !oidcProviderExists {
		logger.Warning("no IAM OIDC provider associated with cluster, try 'eksctl utils associate-iam-oidc-provider --region=%s --cluster=%s'", cmd.ClusterConfig.Metadata.Region, cmd.ClusterConfig.Metadata.Name)
	}

	stackManager := clusterProvider.NewStackManager(cmd.ClusterConfig)

	output, err := clusterProvider.Provider.EKS().DescribeCluster(&awseks.DescribeClusterInput{
		Name: &cmd.ClusterConfig.Metadata.Name,
	})

	if err != nil {
		return fmt.Errorf("failed to fetch cluster %q version: %v", cmd.ClusterConfig.Metadata.Name, err)
	}

	logger.Info("Kubernetes version %q in use by cluster %q", *output.Cluster.Version, cmd.ClusterConfig.Metadata.Name)
	cmd.ClusterConfig.Metadata.Version = *output.Cluster.Version

	addonManager, err := addon.New(cmd.ClusterConfig, clusterProvider, stackManager, oidcProviderExists, oidc, nil)

	if err != nil {
		return err
	}

	for _, a := range cmd.ClusterConfig.Addons {
		if force { //force is specified at cmdline level
			a.Force = true
		}
		err := addonManager.Update(a)
		if err != nil {
			return err
		}
	}

	return nil
}
