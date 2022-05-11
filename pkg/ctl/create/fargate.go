package create

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	actionsfargate "github.com/weaveworks/eksctl/pkg/actions/fargate"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/fargate"
)

func createFargateProfileWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd) error) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"fargateprofile",
		"Create a Fargate profile",
		"",
	)
	options := configureCreateFargateProfileCmd(cmd)
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		if err := cmdutils.NewCreateFargateProfileLoader(cmd, options).Load(); err != nil {
			return err
		}
		return runFunc(cmd)
	}
}

func createFargateProfile(cmd *cmdutils.Cmd) {
	createFargateProfileWithRunFunc(cmd, doCreateFargateProfile)
}

func doCreateFargateProfile(cmd *cmdutils.Cmd) error {
	ctx := context.TODO()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return errors.Wrap(err, "couldn't create cluster provider from command line options")
	}

	manager := actionsfargate.New(cmd.ClusterConfig, ctl, ctl.NewStackManager(cmd.ClusterConfig))
	return manager.Create(ctx)
}

func configureCreateFargateProfileCmd(cmd *cmdutils.Cmd) *fargate.CreateOptions {
	var options fargate.CreateOptions
	cmd.FlagSetGroup.InFlagSet("Fargate", func(fs *pflag.FlagSet) {
		cmdutils.AddFlagsForFargateProfileCreation(fs, &options)
	})
	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
	return &options
}
