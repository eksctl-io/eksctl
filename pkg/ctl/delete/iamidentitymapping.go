package delete

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func deleteIAMIdentityMappingCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	rc.ClusterConfig = cfg

	var (
		role string
		all  bool
	)

	rc.SetDescription("iamidentitymapping", "Delete a IAM identity mapping", "")

	rc.SetRunFunc(func() error {
		return doDeleteIAMIdentityMapping(rc, role, all)
	})

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&role, "role", "", "ARN of the IAM role to delete")
		fs.BoolVar(&all, "all", false, "Delete all matching mappings instead of just one")
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &rc.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &rc.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(rc.FlagSetGroup, rc.ProviderConfig, false)
}

func doDeleteIAMIdentityMapping(rc *cmdutils.ResourceCmd, role string, all bool) error {
	if err := cmdutils.NewMetadataLoader(rc).Load(); err != nil {
		return err
	}

	cfg := rc.ClusterConfig

	ctl := eks.New(rc.ProviderConfig, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if role == "" {
		return cmdutils.ErrMustBeSet("--role")
	}
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet("--name")
	}

	if err := ctl.RefreshClusterConfig(cfg); err != nil {
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

	if err := acm.RemoveRole(role, all); err != nil {
		return err
	}
	if err := acm.Save(); err != nil {
		return err
	}

	// Check whether we have more roles that match
	roles, err := acm.Roles()
	if err != nil {
		return err
	}
	filtered := roles.Get(role)
	if len(filtered) > 0 {
		logger.Warning("there are %d mappings left with same role %q (use --all to delete them at once)", len(filtered), role)
	}
	return nil
}
