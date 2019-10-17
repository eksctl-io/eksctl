package delete

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func deleteIAMIdentityMappingCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	var (
		arn string
		all bool
	)

	cmd.SetDescription("iamidentitymapping", "Delete a IAM identity mapping", "")

	cmd.SetRunFunc(func() error {
		return doDeleteIAMIdentityMapping(cmd, arn, all)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.BoolVar(&all, "all", false, "Delete all matching mappings instead of just one")
		cmdutils.AddIAMIdentityMappingARNFlags(fs, cmd, &arn)
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func doDeleteIAMIdentityMapping(cmd *cmdutils.Cmd, arn string, all bool) error {
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

	if arn == "" {
		return cmdutils.ErrMustBeSet("--arn")
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

	if err := acm.RemoveIdentity(arn, all); err != nil {
		return err
	}
	if err := acm.Save(); err != nil {
		return err
	}

	// Check whether we have more roles that match
	identities, err := acm.Identities()
	if err != nil {
		return err
	}

	duplicates := 0
	for _, identity := range identities {

		if arn == identity.ARN() {
			duplicates++
		}
	}

	if duplicates > 0 {
		logger.Warning("there are %d mappings left with same arn %q (use --all to delete them at once)", duplicates, arn)
	}
	return nil
}
