package update

import (
	"context"
	"fmt"

	awseks "github.com/aws/aws-sdk-go-v2/service/eks"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/addon"
	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func updateAddonCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"addon",
		"Upgrade an Addon",
		"",
	)

	var force, wait bool
	cmd.ClusterConfig.Addons = []*api.Addon{{}}
	cmd.FlagSetGroup.InFlagSet("Addon", func(fs *pflag.FlagSet) {
		fs.StringVar(&cmd.ClusterConfig.Addons[0].Name, "name", "", "Addon name")
		fs.StringVar(&cmd.ClusterConfig.Addons[0].Version, "version", "", "Add-on version. Use `eksctl utils describe-addon-versions` to discover a version or set to \"latest\"")
		fs.StringVar(&cmd.ClusterConfig.Addons[0].ServiceAccountRoleARN, "service-account-role-arn", "", "Addon serviceAccountRoleARN")
		fs.BoolVar(&force, "force", false, "Force migrates an existing self-managed add-on to an EKS managed add-on")
		fs.BoolVar(&wait, "wait", false, "Wait for the addon update to complete")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})
	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return updateAddon(cmd, force, wait)
	}
}

func updateAddon(cmd *cmdutils.Cmd, force, wait bool) error {
	if err := cmdutils.NewCreateOrUpgradeAddonLoader(cmd).Load(); err != nil {
		return err
	}

	ctx := context.Background()
	clusterProvider, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	oidc, err := clusterProvider.NewOpenIDConnectManager(ctx, cmd.ClusterConfig)
	if err != nil {
		return err
	}

	oidcProviderExists, err := oidc.CheckProviderExists(ctx)
	if err != nil {
		return err
	}

	if !oidcProviderExists {
		logger.Warning("no IAM OIDC provider associated with cluster, try 'eksctl utils associate-iam-oidc-provider --region=%s --cluster=%s'", cmd.ClusterConfig.Metadata.Region, cmd.ClusterConfig.Metadata.Name)
	}

	if err := validatePodIdentityAgentAddon(ctx, clusterProvider.AWSProvider.EKS(), cmd.ClusterConfig); err != nil {
		return err
	}

	stackManager := clusterProvider.NewStackManager(cmd.ClusterConfig)

	output, err := clusterProvider.AWSProvider.EKS().DescribeCluster(ctx, &awseks.DescribeClusterInput{
		Name: &cmd.ClusterConfig.Metadata.Name,
	})

	if err != nil {
		return fmt.Errorf("failed to fetch cluster %q version: %v", cmd.ClusterConfig.Metadata.Name, err)
	}

	logger.Info("Kubernetes version %q in use by cluster %q", *output.Cluster.Version, cmd.ClusterConfig.Metadata.Name)
	cmd.ClusterConfig.Metadata.Version = *output.Cluster.Version

	addonManager, err := addon.New(cmd.ClusterConfig, clusterProvider.AWSProvider.EKS(), stackManager, oidcProviderExists, oidc, nil)

	if err != nil {
		return err
	}

	piaUpdater := &addon.PodIdentityAssociationUpdater{
		ClusterName: cmd.ClusterConfig.Metadata.Name,
		IAMRoleCreator: &podidentityassociation.IAMRoleCreator{
			ClusterName:  cmd.ClusterConfig.Metadata.Name,
			StackCreator: stackManager,
		},
		IAMRoleUpdater: &podidentityassociation.IAMRoleUpdater{
			StackUpdater: stackManager,
		},
		EKSPodIdentityDescriber: clusterProvider.AWSProvider.EKS(),
		StackDeleter:            stackManager,
	}

	for _, a := range cmd.ClusterConfig.Addons {
		if force { //force is specified at cmdline level
			a.Force = true
		}
		if err := addonManager.Update(ctx, a, piaUpdater, cmd.ProviderConfig.WaitTimeout); err != nil {
			return err
		}
	}

	return nil
}

func validatePodIdentityAgentAddon(ctx context.Context, eksAPI awsapi.EKS, cfg *api.ClusterConfig) error {
	if isPodIdentityAgentInstalled, err := podidentityassociation.IsPodIdentityAgentInstalled(ctx, eksAPI, cfg.Metadata.Name); err != nil {
		return fmt.Errorf("checking if %q addon is installed on the cluster: %w", api.PodIdentityAgentAddon, err)
	} else if isPodIdentityAgentInstalled {
		return nil
	}

	for _, a := range cfg.Addons {
		if a.HasPodIDsSet() || a.UseDefaultPodIdentityAssociations {
			suggestion := fmt.Sprintf("please enable it using `eksctl create addon --cluster=%s --name=%s`, or by adding it to the config file", cfg.Metadata.Name, api.PodIdentityAgentAddon)
			return api.ErrPodIdentityAgentNotInstalled(suggestion)
		}
	}

	return nil
}
