package enable

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/gitops"
	"github.com/weaveworks/eksctl/pkg/gitops/fileprocessor"
	"github.com/weaveworks/eksctl/pkg/gitops/profile"
)

// ProfileOptions groups input for the "enable profile" command.
type ProfileOptions struct {
	gitOptions     git.Options
	profileOptions profile.Options
}

// Validate validates this ProfileOptions object.
func (opts ProfileOptions) Validate() error {
	if err := opts.gitOptions.Validate(); err != nil {
		return err
	}
	return opts.profileOptions.Validate()
}

func enableProfileCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"profile",
		"Set up Flux and deploy the components from the selected Quick Start profile.",
		"",
	)
	opts := ConfigureProfileCmd(cmd)
	cmd.SetRunFuncWithNameArg(func() error {
		return Profile(cmd, opts)
	})
}

// ConfigureProfileCmd configures the provided command object so that it can
// process CLI options and ClusterConfig file, to prepare for the "enablement"
// of the configured profile on the configured repository & cluster.
func ConfigureProfileCmd(cmd *cmdutils.Cmd) *ProfileOptions {
	var opts ProfileOptions
	cmd.FlagSetGroup.InFlagSet("Enable profile", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForProfile(fs, &opts.profileOptions)
		cmdutils.AddCommonFlagsForGit(fs, &opts.gitOptions)
	})
	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cmd.ClusterConfig.Metadata.Name, "cluster", "", "name of the EKS cluster to enable this Quick Start profile on")
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlagWithValue(fs, &cmd.ProviderConfig.WaitTimeout, 20*time.Second)
	})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
	return &opts
}

// Profile enables the configured profile on the configured repository.
func Profile(cmd *cmdutils.Cmd, opts *ProfileOptions) error {
	if cmd.NameArg != "" && opts.profileOptions.Name != "" {
		return cmdutils.ErrClusterFlagAndArg(cmd, cmd.NameArg, opts.profileOptions.Name)
	}
	if cmd.NameArg != "" {
		opts.profileOptions.Name = cmd.NameArg
	}
	if err := opts.Validate(); err != nil {
		return err
	}
	profileRepoURL, err := profile.RepositoryURL(opts.profileOptions.Name)
	if err != nil {
		return errors.Wrap(err, "please supply a valid Quick Start name or URL")
	}
	if err := cmdutils.NewGitOpsConfigLoader(cmd).Load(); err != nil {
		return err
	}

	// Clone user's repo to apply Quick Start profile
	usersRepoName, err := git.RepoName(opts.gitOptions.URL)
	if err != nil {
		return err
	}
	usersRepoDir, err := ioutil.TempDir("", usersRepoName+"-")
	logger.Debug("Directory %s will be used to clone the configuration repository and install the profile", usersRepoDir)
	profileOutputPath := filepath.Join(usersRepoDir, "base")

	gitClient := git.NewGitClient(git.ClientParams{
		PrivateSSHKeyPath: opts.gitOptions.PrivateSSHKeyPath,
	})

	err = gitClient.CloneRepoInPath(
		usersRepoDir,
		git.CloneOptions{
			URL:       opts.gitOptions.URL,
			Branch:    opts.gitOptions.Branch,
			Bootstrap: true,
		},
	)
	if err != nil {
		return err
	}

	profile := &gitops.Profile{
		Processor: &fileprocessor.GoTemplateProcessor{
			Params: fileprocessor.NewTemplateParameters(cmd.ClusterConfig),
		},
		Path: profileOutputPath,
		GitOpts: git.Options{
			URL:    profileRepoURL,
			Branch: opts.profileOptions.Revision,
		},
		GitCloner: git.NewGitClient(git.ClientParams{
			PrivateSSHKeyPath: opts.gitOptions.PrivateSSHKeyPath,
		}),
		FS: afero.NewOsFs(),
		IO: afero.Afero{Fs: afero.NewOsFs()},
	}

	err = profile.Generate(context.Background())
	if err != nil {
		return errors.Wrap(err, "error generating profile")
	}

	// Git add, commit and push component files in the user's repo
	if err = gitClient.Add("."); err != nil {
		return err
	}
	commitMsg := fmt.Sprintf("Add %s quickstart components", opts.profileOptions.Name)
	if err = gitClient.Commit(commitMsg, opts.gitOptions.User, opts.gitOptions.Email); err != nil {
		return err
	}
	if err = gitClient.Push(); err != nil {
		return err
	}

	profile.DeleteClonedDirectory()
	return nil
}
