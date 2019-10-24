package enable

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/gitops"
	"github.com/weaveworks/eksctl/pkg/gitops/fileprocessor"
	"github.com/weaveworks/eksctl/pkg/gitops/flux"
	"github.com/weaveworks/eksctl/pkg/gitops/profile"
)

// ProfileOptions groups input for the "enable profile" command.
type ProfileOptions struct {
	gitOptions     git.Options
	profileOptions profile.Options
}

func (opts ProfileOptions) validate() error {
	if err := opts.profileOptions.Validate(); err != nil {
		return err
	}
	return cmdutils.ValidateGitOptions(&opts.gitOptions)
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
		fs.StringVarP(&opts.profileOptions.Name, "name", "", "", "name or URL of the Quick Start profile. For example, app-dev.")
		fs.StringVarP(&opts.profileOptions.Revision, "revision", "", "master", "revision of the Quick Start profile.")
		cmdutils.AddCommonFlagsForGit(fs, &opts.gitOptions)

		requiredFlags := []string{"git-url", "git-email"}
		for _, f := range requiredFlags {
			if err := cobra.MarkFlagRequired(fs, f); err != nil {
				logger.Critical("unexpected error: %v", err)
				os.Exit(1)
			}
		}
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
	if err := opts.validate(); err != nil {
		return err
	}

	profileRepoURL, err := profile.RepositoryURL(opts.profileOptions.Name)
	if err != nil {
		return errors.Wrap(err, "please supply a valid Quick Start profile name or URL")
	}

	if err := cmdutils.NewEnableProfileLoader(cmd).Load(); err != nil {
		return err
	}

	k8sClientSet, k8sRestConfig, err := KubernetesClientAndConfigFrom(cmd)
	if err != nil {
		return err
	}

	// Create the flux installer. It will clone the user's repository in a temporary directory.
	fluxOpts := flux.InstallOpts{
		GitOptions:  opts.gitOptions,
		Namespace:   "flux",
		GitFluxPath: "flux/",
		WithHelm:    true,
		Timeout:     cmd.ProviderConfig.WaitTimeout,
	}
	fluxInstaller := flux.NewInstaller(k8sRestConfig, k8sClientSet, &fluxOpts)

	processor := &fileprocessor.GoTemplateProcessor{
		Params: fileprocessor.NewTemplateParameters(cmd.ClusterConfig),
	}

	// Create the profile generator. It will output the processed templates into a new "/base" directory into the user's repo
	usersRepoName, err := git.RepoName(opts.gitOptions.URL)
	if err != nil {
		return err
	}
	dir, err := ioutil.TempDir("", usersRepoName)
	logger.Debug("Directory %s will be used to clone the configuration repository and install the profile", dir)
	usersRepoDir := filepath.Join(dir, usersRepoName)
	profileOutputPath := filepath.Join(usersRepoDir, "base")

	profile := &gitops.Profile{
		Processor: processor,
		Path:      profileOutputPath,
		GitOpts: git.Options{
			URL:    profileRepoURL,
			Branch: opts.profileOptions.Revision,
		},
		GitCloner: git.NewGitClient(git.ClientParams{}),
		FS:        afero.NewOsFs(),
		IO:        afero.Afero{Fs: afero.NewOsFs()},
	}

	// A git client that operates in the user's repo
	gitClient := git.NewGitClient(git.ClientParams{
		PrivateSSHKeyPath: opts.gitOptions.PrivateSSHKeyPath,
	})

	gitOps := gitops.Applier{
		UserRepoPath:     usersRepoDir,
		UsersRepoOpts:    opts.gitOptions,
		GitClient:        gitClient,
		ProfileGenerator: profile,
		FluxInstaller:    fluxInstaller,
		ClusterConfig:    cmd.ClusterConfig,
		QuickstartName:   opts.profileOptions.Name,
	}

	if err = gitOps.Run(context.Background()); err != nil {
		return err
	}
	os.RemoveAll(dir) // Only clean up if the command completely successfully, for more convenient debugging.
	return nil
}
