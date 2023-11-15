package create

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/pkg/errors"
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

	if _, err := ctl.AWSProvider.EKS().DescribeAddon(ctx, &eks.DescribeAddonInput{
		AddonName:   aws.String(api.PodIdentityAgentAddon),
		ClusterName: &cfg.Metadata.Name,
	}); err != nil {
		var notFoundErr *ekstypes.ResourceNotFoundException
		if !errors.As(err, &notFoundErr) {
			return fmt.Errorf("error calling `EKS::DescribeAddon::%s`: %v", api.PodIdentityAgentAddon, err)
		}
		return api.ErrPodIdentityAgentNotInstalled
	}

	return podidentityassociation.NewCreator(cmd.ClusterConfig.Metadata.Name, ctl.NewStackManager(cfg), ctl.AWSProvider.EKS()).
		CreatePodIdentityAssociations(ctx, cmd.ClusterConfig.IAM.PodIdentityAssociations)
}

func configureCreatePodIdentityAssociationCmd(cmd *cmdutils.Cmd, podIdentityAssociation *api.PodIdentityAssociation) {
	cmd.FlagSetGroup.InFlagSet("PodIdentityAssociation", func(fs *pflag.FlagSet) {
		fs.StringVar(&podIdentityAssociation.Namespace, "namespace", "", "Namespace the service account belongs to")
		fs.StringVar(&podIdentityAssociation.ServiceAccountName, "service-account-name", "", "Name of the service account")
		fs.StringVar(&podIdentityAssociation.RoleARN, "role-arn", "", "ARN of the IAM role to be associated with the service account")
		fs.StringVar(&podIdentityAssociation.RoleName, "role-name", "", "Set a custom name for the created role")
		fs.StringVar(&podIdentityAssociation.PermissionsBoundaryARN, "permission-boundary-arn", "", "ARN of the policy that is used to set the permission boundary for the role")

		fs.StringSliceVar(&podIdentityAssociation.PermissionPolicyARNs, "permission-policy-arns", []string{}, "List of ARNs of the IAM permission policies to attach")

		fs.VarP(&podIdentityAssociation.WellKnownPolicies, "well-known-policies", "", "Used to attach common IAM policies")

		cmdutils.AddStringToStringVarPFlag(fs, &podIdentityAssociation.Tags, "tags", "", map[string]string{}, "AWS tags to attach to the PodIdentityAssosciation")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})
}
