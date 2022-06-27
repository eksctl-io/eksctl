package create

import (
	"context"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	mappingactions "github.com/weaveworks/eksctl/pkg/actions/iamidentitymapping"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func createIAMIdentityMappingCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("iamidentitymapping", "Create an IAM identity mapping",
		dedent.Dedent(`Creates a mapping from IAM role or user to Kubernetes user and groups.

			Note aws-iam-authenticator only considers the last entry for any given
			role. If you create a duplicate entry it will shadow all the previous
			username and groups mapping.
		`),
	)

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doCreateIAMIdentityMapping(cmd)
	}

	cfg.IAMIdentityMappings = []*api.IAMIdentityMapping{{}}
	cmd.FlagSetGroup.InFlagSet("IAMIdentityMapping", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.IAMIdentityMappings[0].Account, "account", "", "Account ID to automatically map to its username")
		fs.StringVar(&cfg.IAMIdentityMappings[0].Username, "username", "", "User name within Kubernetes to map to IAM role")
		fs.StringSliceVar(&cfg.IAMIdentityMappings[0].Groups, "group", []string{}, "Group within Kubernetes to which IAM role is mapped")
		fs.StringVar(&cfg.IAMIdentityMappings[0].ServiceName, "service-name", "", "Service name; valid value: emr-containers")
		fs.StringVar(&cfg.IAMIdentityMappings[0].Namespace, "namespace", "", "Namespace in which to create RBAC resources (only valid with --service-name)")
		fs.BoolVar(&cfg.IAMIdentityMappings[0].NoDuplicateArns, "no-duplicate-arns", false, "Throw error when an aws_auth record already exists with the given arn.")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddIAMIdentityMappingARNFlags(fs, cmd, &cfg.IAMIdentityMappings[0].ARN, "create")
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)

}

func doCreateIAMIdentityMapping(cmd *cmdutils.Cmd) error {
	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}
	cfg := cmd.ClusterConfig

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

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

	m, err := mappingactions.New(cfg, clientSet, ctl, cmd.ProviderConfig.Region)
	if err != nil {
		return err
	}

	for _, mapping := range cmd.ClusterConfig.IAMIdentityMappings {
		if err := mapping.Validate(); err != nil {
			return err
		}
		if err := m.Create(ctx, mapping); err != nil {
			return err
		}
	}
	return nil
}
