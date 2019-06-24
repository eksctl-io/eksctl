package create

import (
	"github.com/kris-nova/logger"
	"github.com/lithammer/dedent"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func createIAMIdentityMappingCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	id := &authconfigmap.MapIdentity{}

	rc.SetDescription("iamidentitymapping", "Create an IAM identity mapping",
		dedent.Dedent(`Creates a mapping from IAM role or user to Kubernetes user and groups.

			Note aws-iam-authenticator only considers the last entry for any given
			role. If you create a duplicate entry it will shadow all the previous
			username and groups mapping.
		`),
	)

	cmd.SetRunFunc(func() error {
		return doCreateIAMIdentityMapping(cmd, id)
	})

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&id.ARN, "arn", "", "ARN of the IAM role or user to create")
		fs.StringVar(&id.Username, "username", "", "User name within Kubernetes to map to IAM role")
		fs.StringArrayVar(&id.Groups, "group", []string{}, "Group within Kubernetes to which IAM role is mapped")
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func doCreateIAMIdentityMapping(rc *cmdutils.Cmd, id *authconfigmap.MapIdentity) error {
	if err := cmdutils.NewMetadataLoader(rc).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	logger.Info("using region %s", cfg.Metadata.Region)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}
	if id.ARN == "" {
		return cmdutils.ErrMustBeSet("--arn")
	}
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet("--name")
	}
	if err := id.Valid(); err != nil {
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

	// Check whether role already exists.
	identities, err := acm.Identities()
	if err != nil {
		return err
	}
	filtered := identities.Get(id.ARN)
	if len(filtered) > 0 {
		logger.Warning("found %d mappings with same arn %q (which will be shadowed by your new mapping)", len(filtered), id.ARN)
	}

	if err := acm.AddIdentity(id.ARN, id.Username, id.Groups); err != nil {
		return err
	}
	return acm.Save()
}
