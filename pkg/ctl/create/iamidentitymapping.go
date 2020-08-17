package create

import (
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

	var arn string
	var username string
	var groups []string
	var account string

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doCreateIAMIdentityMapping(cmd, account, arn, username, groups)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&account, "account", "", "Account ID to automatically map to its username")
		fs.StringVar(&username, "username", "", "User name within Kubernetes to map to IAM role")
		fs.StringArrayVar(&groups, "group", []string{}, "Group within Kubernetes to which IAM role is mapped")
		cmdutils.AddIAMIdentityMappingARNFlags(fs, cmd, &arn)
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
}

func doCreateIAMIdentityMapping(cmd *cmdutils.Cmd, account string, arn string, username string, groups []string) error {
	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	cmdutils.LogRegionAndVersionInfo(cfg.Metadata)

	if err := ctl.CheckAuth(); err != nil {
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
	acm, err := authconfigmap.NewFromClientSet(clientSet)
	if err != nil {
		return err
	}

	if account == "" {
		id, err := iam.NewIdentity(arn, username, groups)
		if err != nil {
			return err
		}

		// Check whether role already exists.
		identities, err := acm.Identities()
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
	} else if arn == "" && username == "" && len(groups) == 0 {
		if err := acm.AddAccount(account); err != nil {
			return err
		}
	} else {
		return errors.New("account can only be set alone")
	}
	return acm.Save()
}
