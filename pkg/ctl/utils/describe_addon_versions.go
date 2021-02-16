package utils

import (
	"fmt"

	"github.com/weaveworks/logger"

	awseks "github.com/aws/aws-sdk-go/service/eks"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/actions/addon"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func describeAddonVersionsCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"describe-addon-versions",
		"describe addon versions supported",
		"",
	)

	var addonName, k8sVersion, clusterName string
	cmd.ClusterConfig.Addons = []*api.Addon{{}}
	cmd.FlagSetGroup.InFlagSet("Addon", func(fs *pflag.FlagSet) {
		fs.StringVar(&addonName, "name", "", "Addon name")
		fs.StringVar(&k8sVersion, "kubernetes-version", "", "Kubernetes version")
		fs.StringVarP(&clusterName, "cluster", "c", "", "EKS cluster name")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return describeAddonVersions(cmd, addonName, k8sVersion, clusterName)
	}
}

func describeAddonVersions(cmd *cmdutils.Cmd, addonName, k8sVersion, clusterName string) error {
	clusterProvider, err := cmd.NewCtl()
	if err != nil {
		return err
	}

	//you can provide kubernetes version or cluster name
	//if cluster name we lookup its version
	if k8sVersion != "" {
		cmd.ClusterConfig.Metadata.Version = k8sVersion
	} else if clusterName != "" {
		output, err := clusterProvider.Provider.EKS().DescribeCluster(&awseks.DescribeClusterInput{
			Name: &clusterName,
		})

		if err != nil {
			return fmt.Errorf("failed to fetch cluster %q version: %v", clusterName, err)
		}

		logger.Info("Kubernetes version %q in use by cluster %q", *output.Cluster.Version, clusterName)
		cmd.ClusterConfig.Metadata.Version = *output.Cluster.Version
	} else {
		return fmt.Errorf("cluster name or kubernetes version must be set")
	}

	stackManager := clusterProvider.NewStackManager(cmd.ClusterConfig)

	addonManager, err := addon.New(cmd.ClusterConfig, clusterProvider.Provider.EKS(), stackManager, *cmd.ClusterConfig.IAM.WithOIDC, nil, nil)

	if err != nil {
		return err
	}

	var summary string

	switch addonName {
	case "":
		summary, err = addonManager.DescribeAllVersions()
		if err != nil {
			return err
		}
	default:
		summary, err = addonManager.DescribeVersions(&api.Addon{Name: addonName})
		if err != nil {
			return err
		}
	}

	fmt.Println(summary)

	return nil
}
