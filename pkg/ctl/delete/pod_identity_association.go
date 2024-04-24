package delete

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func deletePodIdentityAssociation(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	var options cmdutils.PodIdentityAssociationOptions

	cmd.SetDescription("podidentityassociation", "Delete pod identity associations", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		return doDeletePodIdentityAssociation(cmd, options)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmd.FlagSetGroup.InFlagSet("Pod Identity Association", func(fs *pflag.FlagSet) {
		fs.StringVar(&options.Namespace, "namespace", "", "Namespace of the pod identity association")
		fs.StringVar(&options.ServiceAccountName, "service-account-name", "", "Service account name of the pod identity association")

	})

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)
}

func doDeletePodIdentityAssociation(cmd *cmdutils.Cmd, options cmdutils.PodIdentityAssociationOptions) error {
	if err := cmdutils.NewDeletePodIdentityAssociationLoader(cmd, options).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	ctx := context.Background()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	if cmd.ClusterConfigFile == "" {
		cmd.ClusterConfig.IAM.PodIdentityAssociations = []api.PodIdentityAssociation{
			{
				Namespace:          options.Namespace,
				ServiceAccountName: options.ServiceAccountName,
			},
		}
	}

	deleter := &podidentityassociation.Deleter{
		ClusterName:  cfg.Metadata.Name,
		StackDeleter: ctl.NewStackManager(cfg),
		APIDeleter:   ctl.AWSProvider.EKS(),
		ClientSet:    clientSet,
	}

	return deleter.Delete(ctx, podidentityassociation.ToIdentifiers(cfg.IAM.PodIdentityAssociations))
}
