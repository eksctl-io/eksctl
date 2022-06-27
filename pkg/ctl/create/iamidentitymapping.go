package create

import (
	"context"
	"fmt"

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

	var options api.IAMIdentityMapping

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
		fs.BoolVar(&options.NoDuplicateArns, "no-duplicate-arns", false, "Throw error when an aws_auth record already exists with the given arn.")
		cmdutils.AddIAMIdentityMappingARNFlags(fs, cmd, &options.ARN, "create")
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
}

func doCreateIAMIdentityMapping(cmd *cmdutils.Cmd, options api.IAMIdentityMapping) error {
	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}

	if cmd.ClusterConfigFile != "" {
		for _, mapping := range cmd.ClusterConfig.IAMIdentityMappings {
			if err := mapping.Validate(); err != nil {
			    return err
			}
			if err := createIAMIdentityMapping(cmd, *mapping); err != nil {
				return err
			}
		}
	} else {
		err := options.Validate()
		if err != nil {
			return err
		}
		err = createIAMIdentityMapping(cmd, options)
		if err != nil {
			return err
		}
	}
	return nil
}

func createIAMIdentityMapping(cmd *cmdutils.Cmd, options api.IAMIdentityMapping) error {
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

	acm, err := authconfigmap.NewFromClientSet(clientSet)
	if err != nil {
		return err
	}

	if options.ServiceName != "" {
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
		logger.Info("checking arn %s against entries in the configmap", id.ARN())
		for _, identity := range identities {
			arn := identity.ARN()
			if options.NoDuplicateArns && iam.CompareIdentity(id, identity) {
				logger.Warning("found existing mapping that matches the one being created, skipping.")
				return nil
			}

			if createdArn == arn && options.NoDuplicateArns {
				return fmt.Errorf("found existing mapping with the same arn %q and shadowing is disabled", createdArn)
			}

			if createdArn == arn {
				logger.Warning("found existing mappings with same arn %q (which will be shadowed by your new mapping)", createdArn)
				break
			}
		}

		if err := acm.AddIdentity(id); err != nil {
			return err
		}
	} else {
		if err := acm.AddAccount(options.Account); err != nil {
			return err
		}
	}
	return acm.Save()
}
