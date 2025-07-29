package update

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
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
		fs.StringVar(&options.RoleARN, "role-arn", "", "ARN of the IAM role to be associated with the service account")
		var targetRoleArn string
		var disableSessionTags, noDisableSessionTags bool
		fs.StringVar(&targetRoleArn, "target-role-arn", "", "ARN of the target IAM role for cross-account access")
		fs.BoolVar(&disableSessionTags, "disable-session-tags", false, "Disable session tags added by EKS Pod Identity")
		fs.BoolVar(&noDisableSessionTags, "no-disable-session-tags", false, "Enable session tags added by EKS Pod Identity")
		cmdutils.AddPreRun(cmd.CobraCommand, func(cobraCmd *cobra.Command, args []string) {
			if fs.Changed("target-role-arn") {
				options.TargetRoleARN = &targetRoleArn
			}
			if fs.Changed("no-disable-session-tags") {
				options.DisableSessionTags = aws.Bool(false)
			} else if fs.Changed("disable-session-tags") {
				options.DisableSessionTags = aws.Bool(true)
			}
		})
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
				TargetRoleARN:      options.TargetRoleARN,
				DisableSessionTags: options.DisableSessionTags,
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
