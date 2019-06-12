package create

import (
	"github.com/kris-nova/logger"
	"github.com/lithammer/dedent"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func createIAMIdentityMappingCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	rc.ClusterConfig = cfg

	id := &authconfigmap.MapRole{}

	rc.SetDescription("iamidentitymapping", "Create an IAM identity mapping",
		dedent.Dedent(`Creates a mapping from IAM role to Kubernetes user and groups.

			Note aws-iam-authenticator only considers the last entry for any given
			role. If you create a duplicate entry it will shadow all the previous
			username and groups mapping.
		`),
	)

	rc.SetRunFunc(func() error {
		return doCreateIAMIdentityMapping(rc, id)
	})

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&id.RoleARN, "role", "", "ARN of the IAM role to create")
		fs.StringVar(&id.Username, "username", "", "User name within Kubernetes to map to IAM role")
		fs.StringArrayVar(&id.Groups, "group", []string{}, "Group within Kubernetes to which IAM role is mapped")
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &rc.ClusterConfigFile)
	})

	cmdutils.AddCommonFlagsForAWS(rc.FlagSetGroup, rc.ProviderConfig, false)
}

func doCreateIAMIdentityMapping(rc *cmdutils.ResourceCmd, id *authconfigmap.MapRole) error {
	if err := cmdutils.NewMetadataLoader(rc).Load(); err != nil {
		return err
	}

	cfg := rc.ClusterConfig

	ctl := eks.New(rc.ProviderConfig, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}
	if id.RoleARN == "" {
		return cmdutils.ErrMustBeSet("--role")
	}
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet("--name")
	}
	if err := id.Valid(); err != nil {
		return err
	}

	if err := ctl.GetCredentials(cfg); err != nil {
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
	roles, err := acm.Roles()
	if err != nil {
		return err
	}
	filtered := roles.Get(id.RoleARN)
	if len(filtered) > 0 {
		logger.Warning("found %d mappings with same role %q (which will be shadowed by your new mapping)", len(filtered), id.RoleARN)
	}

	if err := acm.AddRole(id.RoleARN, id.Username, id.Groups); err != nil {
		return err
	}
	return acm.Save()
}
