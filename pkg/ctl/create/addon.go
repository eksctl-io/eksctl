package create

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"

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

func createAddonCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"addon",
		"Create an Addon",
		"",
	)

	var force, wait bool
	cmd.ClusterConfig.Addons = []*api.Addon{{}}
	cmd.FlagSetGroup.InFlagSet("Addon", func(fs *pflag.FlagSet) {
		fs.StringVar(&cmd.ClusterConfig.Addons[0].Name, "name", "", "Add-on name")
		fs.StringVar(&cmd.ClusterConfig.Addons[0].Version, "version", "", "Add-on version. Use `eksctl utils describe-addon-versions` to discover a version or set to \"latest\"")
		fs.StringVar(&cmd.ClusterConfig.Addons[0].ServiceAccountRoleARN, "service-account-role-arn", "", "Add-on serviceAccountRoleARN")
		fs.BoolVar(&cmd.ClusterConfig.AddonsConfig.AutoApplyPodIdentityAssociations, "auto-apply-pod-identity-associations", false, "apply recommended pod identity associations for the addon(s), if supported")
		fs.BoolVar(&force, "force", false, "Force migrates an existing self-managed add-on to an EKS managed add-on")
		fs.BoolVar(&wait, "wait", false, "Wait for the addon creation to complete")

		fs.StringSliceVar(&cmd.ClusterConfig.Addons[0].AttachPolicyARNs, "attach-policy-arn", []string{}, "ARN of the policies to attach")
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
		if err := cmdutils.NewCreateOrUpgradeAddonLoader(cmd).Load(); err != nil {
			return err
		}

		ctx := context.TODO()
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

		addonManager, err := addon.New(cmd.ClusterConfig, clusterProvider.AWSProvider.EKS(), stackManager, oidcProviderExists, oidc, func() (kubernetes.Interface, error) {
			return clusterProvider.NewStdClientSet(cmd.ClusterConfig)
		})
		if err != nil {
			return err
		}

		iamRoleCreator := &podidentityassociation.IAMRoleCreator{
			ClusterName:  cmd.ClusterConfig.Metadata.Name,
			StackCreator: stackManager,
		}
		// always install EKS Pod Identity Agent Addon first, if present,
		// as other addons might require IAM permissions
		for _, a := range cmd.ClusterConfig.Addons {
			if a.CanonicalName() != api.PodIdentityAgentAddon {
				continue
			}
			if force { //force is specified at cmdline level
				a.Force = true
			}
			if err := addonManager.Create(ctx, a, iamRoleCreator, cmd.ProviderConfig.WaitTimeout); err != nil {
				return err
			}
		}

		for _, a := range cmd.ClusterConfig.Addons {
			if a.CanonicalName() == api.PodIdentityAgentAddon {
				continue
			}
			if force { //force is specified at cmdline level
				a.Force = true
			}
			if err := addonManager.Create(ctx, a, iamRoleCreator, cmd.ProviderConfig.WaitTimeout); err != nil {
				return err
			}
		}

		return nil
	}
}

func validatePodIdentityAgentAddon(ctx context.Context, eksAPI awsapi.EKS, cfg *api.ClusterConfig) error {
	isPodIdentityAgentInstalled, err := podidentityassociation.IsPodIdentityAgentInstalled(ctx, eksAPI, cfg.Metadata.Name)
	if err != nil {
		return err
	}

	shallCreatePodIdentityAssociations := cfg.AddonsConfig.AutoApplyPodIdentityAssociations
	podIdentityAgentFoundInConfig := false
	for _, a := range cfg.Addons {
		if a.CanonicalName() == api.PodIdentityAgentAddon {
			podIdentityAgentFoundInConfig = true
		}
		if a.HasPodIDsSet() || a.UseDefaultPodIdentityAssociations {
			shallCreatePodIdentityAssociations = true
		}
	}

	if shallCreatePodIdentityAssociations && !isPodIdentityAgentInstalled && !podIdentityAgentFoundInConfig {
		suggestion := fmt.Sprintf("please enable it using `eksctl create addon --cluster=%s --name=%s`, or by adding it to the config file", cfg.Metadata.Name, api.PodIdentityAgentAddon)
		return api.ErrPodIdentityAgentNotInstalled(suggestion)
	}

	return nil
}
