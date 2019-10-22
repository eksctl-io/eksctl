package generate

import (
	"context"

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

type options struct {
	GitOptions        git.Options
	ProfilePath       string
	PrivateSSHKeyPath string
}

func generateProfileCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("profile", "Generate a gitops profile", "")

	var o options

	cmd.SetRunFuncWithNameArg(func() error {
		return doGenerateProfile(cmd, o)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&o.GitOptions.URL, "git-url", "", "", "URL for the quickstart base repository")
		fs.StringVarP(&o.GitOptions.Branch, "git-branch", "", "master", "Git branch")
		fs.StringVarP(&o.ProfilePath, "profile-path", "", "./", "Path to generate the profile in")
		_ = cobra.MarkFlagRequired(fs, "git-url")

		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func doGenerateProfile(cmd *cmdutils.Cmd, o options) error {
	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}

	// TODO move the load of the region outside of the creation of the EKS client
	// currently that is done inside cmd.NewCtl() but we don't need EKS here
	cmd.ClusterConfig.Metadata.Region = cmd.ProviderConfig.Region

	processor := &fileprocessor.GoTemplateProcessor{
		Params: fileprocessor.NewTemplateParameters(cmd.ClusterConfig),
	}
	profile := &gitops.Profile{
		Processor: processor,
		Path:      o.ProfilePath,
		GitOpts:   o.GitOptions,
		GitCloner: git.NewGitClient(git.ClientParams{
			PrivateSSHKeyPath: o.PrivateSSHKeyPath,
		}),
		FS: afero.NewOsFs(),
		IO: afero.Afero{Fs: afero.NewOsFs()},
	}

	err := profile.Generate(context.Background())
	if err != nil {
		return errors.Wrap(err, "error generating profile")
	}

	profile.DeleteClonedDirectory()
	return nil
}
