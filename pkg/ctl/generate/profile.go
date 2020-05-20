package generate

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/gitops"
	"github.com/weaveworks/eksctl/pkg/gitops/fileprocessor"
)

func generateProfile(cmd *cmdutils.Cmd) {
	generateProfileWithRunFunc(cmd, doGenerateProfile)
}

func generateProfileWithRunFunc(cmd *cmdutils.Cmd, runFunc func(*cmdutils.Cmd) error) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription("profile", "Generate a gitops profile", "")

	opts := configureGenerateProfileCmd(cmd)

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)

		if err := cmdutils.NewGitOpsConfigLoader(cmd, opts).
			WithProfileValidation().Load(); err != nil {
			return err
		}

		return runFunc(cmd)
	}

}

func configureGenerateProfileCmd(cmd *cmdutils.Cmd) *api.Git {
	opts := api.NewGit()

	cmd.FlagSetGroup.InFlagSet("Generate profile", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForProfile(fs, opts.BootstrapProfile)
		fs.StringVarP(&opts.BootstrapProfile.OutputPath, "profile-path", "", "", "path to generate the profile in. Defaults to ./<quickstart-repo-name>")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cmd.ClusterConfig.Metadata.Name, "cluster", "", "name of the EKS cluster to enable gitops on")
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
	return opts
}

func doGenerateProfile(cmd *cmdutils.Cmd) error {
	// TODO move the load of the region outside of the creation of the EKS client
	// currently that is done inside cmd.NewCtl() but we don't need EKS here
	cmd.ClusterConfig.Metadata.Region = cmd.ProviderConfig.Region

	processor := &fileprocessor.GoTemplateProcessor{
		Params: fileprocessor.NewTemplateParameters(cmd.ClusterConfig),
	}
	profile := &gitops.Profile{
		Processor:     processor,
		ProfileCloner: git.NewGitClient(git.ClientParams{}),
		FS:            afero.NewOsFs(),
		IO:            afero.Afero{Fs: afero.NewOsFs()},
	}

	bootstrapProfile := cmd.ClusterConfig.Git.BootstrapProfile
	err := profile.Generate(*bootstrapProfile)
	if err != nil {
		return errors.Wrap(err, "error generating profile")
	}

	return nil
}
