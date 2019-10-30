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
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/gitops"
	"github.com/weaveworks/eksctl/pkg/gitops/fileprocessor"
	"github.com/weaveworks/eksctl/pkg/gitops/flux"
	"github.com/weaveworks/eksctl/pkg/gitops/profile"
)

// Options groups input for the "create cluster" command.
type Options struct {
	fluxOptions    flux.InstallOpts
	profileOptions profile.Options
}

func (opts Options) validate() error {
	if err := opts.profileOptions.Validate(); err != nil {
		return err
	}
	return cmdutils.ValidateGitOptions(&opts.fluxOptions.GitOptions)
}

// AllCmd combines the functionality of
// - eksctl enable repo
// - eksctl enable profile
// TODO: eventually have eksctl create cluster run something equivalent to the
// below, once ClusterConfig supports GitOps settings.
func AllCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"all",
		"Set up Flux and deploy the components from the selected Quick Start profile.",
		"",
	)
	opts := ConfigureAllCmd(cmd)
	cmd.SetRunFuncWithNameArg(func() error {
		return All(cmd, opts)
	})
}

// ConfigureAllCmd configures the provided command object so that it can
// process CLI options and ClusterConfig file, to prepare for the "enablement"
// of the configured profile on the configured repository & cluster.
func ConfigureAllCmd(cmd *cmdutils.Cmd) *Options {
	var opts Options
	cmd.FlagSetGroup.InFlagSet("Enable profile", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForProfile(fs, &opts.profileOptions)
		cmdutils.AddCommonFlagsForFlux(fs, &opts.fluxOptions)
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

// All enables the configured profile on the configured repository.
func All(cmd *cmdutils.Cmd, opts *Options) error {
	if cmd.NameArg != "" && opts.profileOptions.Name != "" {
		return cmdutils.ErrClusterFlagAndArg(cmd, cmd.NameArg, opts.profileOptions.Name)
	}
	if cmd.NameArg != "" {
		opts.profileOptions.Name = cmd.NameArg
	}
	if err := opts.validate(); err != nil {
		return err
	}

	if err := cmdutils.NewGitOpsConfigLoader(cmd).Load(); err != nil {
		return err
	}

	gitOpsApplier, tempDir, err := newGitOpsApplier(cmd, opts)
	if err = gitOpsApplier.Run(context.Background()); err != nil {
		return err
	}
	os.RemoveAll(tempDir) // Only clean up if the command completely successfully, for more convenient debugging.
	return nil
}

func newGitOpsApplier(cmd *cmdutils.Cmd, opts *Options) (*gitops.Applier, string, error) {
	// Create the profile generator. It will output the processed templates into a new "/base" directory into the user's repo
	usersRepoName, err := git.RepoName(opts.fluxOptions.GitOptions.URL)
	if err != nil {
		return nil, "", err
	}
	dir, err := ioutil.TempDir("", usersRepoName)
	logger.Debug("Directory %s will be used to clone the configuration repository and install the profile", dir)
	usersRepoDir := filepath.Join(dir, usersRepoName)
	profileOutputPath := filepath.Join(usersRepoDir, "base")

	profileRepoURL, err := profile.RepositoryURL(opts.profileOptions.Name)
	if err != nil {
		return nil, "", errors.Wrap(err, "please supply a valid Quick Start profile name or URL")
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
		GitCloner: git.NewGitClient(git.ClientParams{}),
		FS:        afero.NewOsFs(),
		IO:        afero.Afero{Fs: afero.NewOsFs()},
	}

	k8sClientSet, k8sRestConfig, err := KubernetesClientAndConfigFrom(cmd)
	if err != nil {
		return nil, "", err
	}
	gitOps := gitops.Applier{
		UserRepoPath:  usersRepoDir,
		UsersRepoOpts: opts.fluxOptions.GitOptions,
		GitClient: git.NewGitClient(git.ClientParams{
			PrivateSSHKeyPath: opts.fluxOptions.GitOptions.PrivateSSHKeyPath,
		}),
		ProfileGenerator: profile,
		FluxInstaller:    flux.NewInstaller(k8sRestConfig, k8sClientSet, &opts.fluxOptions),
		ClusterConfig:    cmd.ClusterConfig,
		QuickstartName:   opts.profileOptions.Name,
	}
	return &gitOps, dir, nil
}
