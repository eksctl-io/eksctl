package create

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
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

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
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

	return podidentityassociation.NewCreator(cmd.ClusterConfig.Metadata.Name, ctl.NewStackManager(cfg), ctl.AWSProvider.EKS(), clientSet).
		CreatePodIdentityAssociations(ctx, cmd.ClusterConfig.IAM.PodIdentityAssociations)
}

func configureCreatePodIdentityAssociationCmd(cmd *cmdutils.Cmd, pia *api.PodIdentityAssociation) {
	cmd.FlagSetGroup.InFlagSet("PodIdentityAssociation", func(fs *pflag.FlagSet) {
		fs.StringVar(&pia.Namespace, "namespace", "", "Namespace the service account belongs to")
		fs.StringVar(&pia.ServiceAccountName, "service-account-name", "", "Name of the service account")
		fs.StringVar(&pia.RoleARN, "role-arn", "", "ARN of the IAM role to be associated with the service account")
		fs.StringVar(&pia.RoleName, "role-name", "", "Set a custom name for the created role")
		fs.StringVar(&pia.PermissionsBoundaryARN, "permission-boundary-arn", "", "ARN of the policy that is used to set the permission boundary for the role")
		var targetRoleARN string
		var disableSessionTags bool
		fs.StringVar(&targetRoleARN, "target-role-arn", "", "ARN of the target IAM role for cross-account access (default to empty string for no cross-account access)")
		fs.BoolVar(&disableSessionTags, "disable-session-tags", false, "Disable session tags added by EKS Pod Identity (if not provided, session tags are enabled by default)")

		// Store the flag values in the struct
		cmdutils.AddPreRun(cmd.CobraCommand, func(cobraCmd *cobra.Command, args []string) {
			if fs.Changed("target-role-arn") {
				pia.TargetRoleARN = &targetRoleARN
			}
			if fs.Changed("disable-session-tags") {
				pia.DisableSessionTags = aws.Bool(true)
			}
		})

		fs.BoolVar(&pia.CreateServiceAccount, "create-service-account", false, "instructs eksctl to create the K8s service account")

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
