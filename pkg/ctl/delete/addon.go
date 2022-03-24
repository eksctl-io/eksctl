package delete

import (
	"context"
	"fmt"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/actions/addon"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func deleteAddonCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"addon",
		"Delete an Addon",
		"",
	)

	cmd.ClusterConfig.Addons = []*api.Addon{{}}
	var preserve bool
	cmd.FlagSetGroup.InFlagSet("Addon", func(fs *pflag.FlagSet) {
		fs.StringVar(&cmd.ClusterConfig.Addons[0].Name, "name", "", "Addon name")
		fs.BoolVar(&preserve, "preserve", false, "Delete the addon from the API but preserve its Kubernetes resources")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return deleteAddon(cmd, preserve)
	}
}

func deleteAddon(cmd *cmdutils.Cmd, preserve bool) error {
	if err := cmdutils.NewDeleteAddonLoader(cmd).Load(); err != nil {
		return err
	}

	clusterProvider, err := cmd.NewProviderForExistingCluster()
	if err != nil {
		return err
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

	addonManager, err := addon.New(cmd.ClusterConfig, clusterProvider.Provider.EKS(), stackManager, *cmd.ClusterConfig.IAM.WithOIDC, nil, nil, cmd.ProviderConfig.WaitTimeout)

	if err != nil {
		return err
	}

	if preserve {
		return addonManager.DeleteWithPreserve(cmd.ClusterConfig.Addons[0])
	}

	return addonManager.Delete(context.TODO(), cmd.ClusterConfig.Addons[0])
}
