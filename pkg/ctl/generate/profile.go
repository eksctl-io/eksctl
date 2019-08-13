package generate

import (
	"context"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/gitops/fileprocessor"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/gitops"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

const (
	defaultGitTimeout = 20 * time.Second
)

// Command creates `generate` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("generate", "Generate GitOps manifests", "")
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, generateProfileCmd)
	return verbCmd
}

type options struct {
	gitops.GitOptions
	ProfilePath string
}

func generateProfileCmd(rc *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	rc.ClusterConfig = cfg

	rc.SetDescription("profile", "Generate a GitOps profile", "")

	var o options

	rc.SetRunFuncWithNameArg(func() error {
		return doGenerateProfile(rc, o)
	})

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&o.URL, "git-url", "", "", "Git repository URL")
		fs.StringVarP(&o.Branch, "git-branch", "", "master", "Git branch; defaults to master")
		fs.StringVarP(&o.ProfilePath, "profile-path", "", "", "Path to generate the profile in; defaults to CWD")
		_ = cobra.MarkFlagRequired(fs, "git-url")

		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &rc.ClusterConfigFile)
	})

	cmdutils.AddCommonFlagsForAWS(rc.FlagSetGroup, rc.ProviderConfig, false)
}

func doGenerateProfile(rc *cmdutils.Cmd, o options) error {
	if err := cmdutils.NewMetadataLoader(rc).Load(); err != nil {
		return err
	}

	processor := &fileprocessor.GoTemplateProcessor{
		Params: fileprocessor.NewTemplateParameters(rc.ClusterConfig),
	}
	profile := &gitops.Profile{
		Processor: processor,
		Path:      o.ProfilePath,
		GitOpts:   o.GitOptions,
		GitCloner: git.NewGitClient(context.Background(), defaultGitTimeout),
		Fs:        afero.NewOsFs(),
		IO:        afero.Afero{Fs: afero.NewOsFs()},
	}

	err := profile.Generate(context.Background())

	if err != nil {
		return errors.Wrap(err, "error generating profile")
	}

	return nil
}
