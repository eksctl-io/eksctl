package delete

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

type deleteIAMIdentityMappingCmdParams struct {
	arn             string
	all             bool
	clusterEndpoint string
}

func deleteIAMIdentityMappingCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	params := &deleteIAMIdentityMappingCmdParams{}

	cmd.SetDescription("iamidentitymapping", "Delete a IAM identity mapping", "")

	cmd.SetRunFunc(func() error {
		return doDeleteIAMIdentityMapping(cmd, params)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
		cmdutils.AddClusterEndpointOverrideFlag(fs, &params.clusterEndpoint)
	})

	cmd.FlagSetGroup.InFlagSet("Delete IAM Identity Mapping", func(fs *pflag.FlagSet) {
		fs.BoolVar(&params.all, "all", false, "Delete all matching mappings instead of just one")
		cmdutils.AddIAMIdentityMappingARNFlags(fs, cmd, &params.arn)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func doDeleteIAMIdentityMapping(cmd *cmdutils.Cmd, params *deleteIAMIdentityMappingCmdParams) error {
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

	if params.arn == "" {
		return cmdutils.ErrMustBeSet("--arn")
	}
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}
	clientSet, err := ctl.NewStdClientSet(cfg, params.clusterEndpoint)
	if err != nil {
		return err
	}
	acm, err := authconfigmap.NewFromClientSet(clientSet)
	if err != nil {
		return err
	}

	if err := acm.RemoveIdentity(params.arn, params.all); err != nil {
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

		if params.arn == identity.ARN() {
			duplicates++
		}
	}

	if duplicates > 0 {
		logger.Warning("there are %d mappings left with same arn %q (use --all to delete them at once)", duplicates, params.arn)
	}
	return nil
}
