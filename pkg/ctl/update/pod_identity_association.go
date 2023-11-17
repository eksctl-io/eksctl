package update

import (
	"context"

	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func updatePodIdentityAssociation(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	var options cmdutils.UpdatePodIdentityAssociationOptions

	cmd.SetDescription("podidentityassociation", "Update pod identity associations", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		return doUpdatePodIdentityAssociation(cmd, options)
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
		fs.StringVar(&options.RoleARN, "role-arn", "", "Service account name of the pod identity association")

	})

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)
}

func doUpdatePodIdentityAssociation(cmd *cmdutils.Cmd, options cmdutils.UpdatePodIdentityAssociationOptions) error {
	if err := cmdutils.NewUpdatePodIdentityAssociationLoader(cmd, options).Load(); err != nil {
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

	if cmd.ClusterConfigFile == "" {
		cmd.ClusterConfig.IAM.PodIdentityAssociations = []api.PodIdentityAssociation{
			{
				Namespace:          options.Namespace,
				ServiceAccountName: options.ServiceAccountName,
				RoleARN:            options.RoleARN,
			},
		}
	}

	stackManager := ctl.NewStackManager(cfg)
	updater := &podidentityassociation.Updater{
		ClusterName:  cfg.Metadata.Name,
		APIUpdater:   ctl.AWSProvider.EKS(),
		StackUpdater: stackManager,
	}
	return updater.Update(ctx, cfg.IAM.PodIdentityAssociations)
}
