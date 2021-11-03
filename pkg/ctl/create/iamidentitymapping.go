package create

import (
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/kris-nova/logger"
	"github.com/lithammer/dedent"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/iam"
)

type iamIdentityMappingOptions struct {
	ARN         string
	Username    string
	Groups      []string
	Account     string
	ServiceName string
	Namespace   string
}

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

	var options iamIdentityMappingOptions

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doCreateIAMIdentityMapping(cmd, options)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&options.Account, "account", "", "Account ID to automatically map to its username")
		fs.StringVar(&options.Username, "username", "", "User name within Kubernetes to map to IAM role")
		fs.StringSliceVar(&options.Groups, "group", []string{}, "Group within Kubernetes to which IAM role is mapped")
		fs.StringVar(&options.ServiceName, "service-name", "", "Service name; valid value: emr-containers")
		fs.StringVar(&options.Namespace, "namespace", "", "Namespace in which to create RBAC resources (only valid with --service-name)")
		cmdutils.AddIAMIdentityMappingARNFlags(fs, cmd, &options.ARN, "create")
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
}

func doCreateIAMIdentityMapping(cmd *cmdutils.Cmd, options iamIdentityMappingOptions) error {
	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

	ctl, err := cmd.NewProviderForExistingCluster()
	if err != nil {
		return err
	}
	cmdutils.LogRegionAndVersionInfo(cfg.Metadata)

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}
	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	hasARNOptions := func() bool {
		return !(options.ARN == "" && options.Username == "" && len(options.Groups) == 0)
	}

	validateNonServiceOptions := func() error {
		if options.Namespace != "" {
			return errors.New("--namespace is only valid with --service-name")
		}
		return nil
	}

	acm, err := authconfigmap.NewFromClientSet(clientSet)
	if err != nil {
		return err
	}

	if options.ServiceName != "" {
		if hasARNOptions() {
			return errors.New("cannot use --arn, --username, and --groups with --service-name")
		}
		rawClient, err := ctl.NewRawClient(cfg)
		if err != nil {
			return err
		}
		parsedARN, err := arn.Parse(cfg.Status.ARN)
		if err != nil {
			return errors.Wrap(err, "error parsing cluster ARN")
		}
		sa := authconfigmap.NewServiceAccess(rawClient, acm, parsedARN.AccountID)
		return sa.Grant(options.ServiceName, options.Namespace, api.Partition(cmd.ProviderConfig.Region))
	}

	if options.Account == "" {
		if err := validateNonServiceOptions(); err != nil {
			return err
		}
		id, err := iam.NewIdentity(options.ARN, options.Username, options.Groups)
		if err != nil {
			return err
		}

		// Check whether role already exists.
		identities, err := acm.GetIdentities()
		if err != nil {
			return err
		}

		createdArn := id.ARN() // The call to Valid above makes sure this cannot error
		for _, identity := range identities {
			arn := identity.ARN()

			if createdArn == arn {
				logger.Warning("found existing mappings with same arn %q (which will be shadowed by your new mapping)", createdArn)
				break
			}
		}

		if err := acm.AddIdentity(id); err != nil {
			return err
		}
	} else if hasARNOptions() {
		if err := validateNonServiceOptions(); err != nil {
			return err
		}
		if err := acm.AddAccount(options.Account); err != nil {
			return err
		}
	} else {
		return errors.New("account can only be set alone")
	}
	return acm.Save()
}
