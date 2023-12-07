package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func createPodIdentityAssociationCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"podidentityassociation",
		"Create a pod identity association",
		"",
	)

	podIdentityAssociation := &api.PodIdentityAssociation{}
	configureCreatePodIdentityAssociationCmd(cmd, podIdentityAssociation)

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		if err := cmdutils.NewCreatePodIdentityAssociationLoader(cmd, podIdentityAssociation).Load(); err != nil {
			return err
		}
		return doCreatePodIdentityAssociation(cmd)
	}
}

func doCreatePodIdentityAssociation(cmd *cmdutils.Cmd) error {
	cfg := cmd.ClusterConfig
	ctx := context.Background()

	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	isInstalled, err := podidentityassociation.IsPodIdentityAgentInstalled(ctx, ctl.AWSProvider.EKS(), cfg.Metadata.Name)
	if err != nil {
		return err
	}

	if !isInstalled {
		suggestion := fmt.Sprintf("please enable it using `eksctl create addon --cluster=%s --name=%s`", cmd.ClusterConfig.Metadata.Name, api.PodIdentityAgentAddon)
		return api.ErrPodIdentityAgentNotInstalled(suggestion)
	}

	return podidentityassociation.NewCreator(cmd.ClusterConfig.Metadata.Name, ctl.NewStackManager(cfg), ctl.AWSProvider.EKS()).
		CreatePodIdentityAssociations(ctx, cmd.ClusterConfig.IAM.PodIdentityAssociations)
}

func configureCreatePodIdentityAssociationCmd(cmd *cmdutils.Cmd, pia *api.PodIdentityAssociation) {
	cmd.FlagSetGroup.InFlagSet("PodIdentityAssociation", func(fs *pflag.FlagSet) {
		fs.StringVar(&pia.Namespace, "namespace", "", "Namespace the service account belongs to")
		fs.StringVar(&pia.ServiceAccountName, "service-account-name", "", "Name of the service account")
		fs.StringVar(&pia.RoleARN, "role-arn", "", "ARN of the IAM role to be associated with the service account")
		fs.StringVar(&pia.RoleName, "role-name", "", "Set a custom name for the created role")
		fs.StringVar(&pia.PermissionsBoundaryARN, "permission-boundary-arn", "", "ARN of the policy that is used to set the permission boundary for the role")

		fs.StringSliceVar(&pia.PermissionPolicyARNs, "permission-policy-arns", []string{}, "List of ARNs of the IAM permission policies to attach")

		fs.VarP(&pia.WellKnownPolicies, "well-known-policies", "", "Used to attach common IAM policies")

		cmdutils.AddStringToStringVarPFlag(fs, &pia.Tags, "tags", "", map[string]string{}, "AWS tags to attach to the PodIdentityAssosciation")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})
}
