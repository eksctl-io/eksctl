package delete

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/actions/identitymapping"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func deleteIAMIdentityMappingCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	var (
		all bool
	)

	cmd.SetDescription("iamidentitymapping", "Delete a IAM identity mapping", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		return doDeleteIAMIdentityMapping(cmd, all)
	}

	options := &api.IAMIdentityMapping{}

	cfg.IAM.IdentityMapping = append(cfg.IAM.IdentityMapping, options)

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.BoolVar(&all, "all", false, "Delete all matching mappings instead of just one")
		cmdutils.AddIAMIdentityMappingARNFlags(fs, cmd, &options.ARN)
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
}

func doDeleteIAMIdentityMapping(cmd *cmdutils.Cmd, all bool) error {
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

	if len(cfg.IAM.IdentityMapping) == 0 {
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
	return identitymapping.New(nil, acm).Delete(cfg.IAM.IdentityMapping, all)
}
