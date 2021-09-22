package cmdutils

import (
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/git"
)

const (
	gitURL               = "git-url"
	gitBranch            = "git-branch"
	gitUser              = "git-user"
	gitEmail             = "git-email"
	gitPrivateSSHKeyPath = "git-private-ssh-key-path"

	gitPaths                   = "git-paths"
	gitFluxPath                = "git-flux-subdir"
	gitLabel                   = "git-label"
	namespace                  = "namespace"
	readOnly                   = "read-only"
	withHelm                   = "with-helm"
	additionalFluxArgs         = "additional-flux-args"
	additionalHelmOperatorArgs = "additional-helm-operator-args"

	commitOperatorManifests = "commit-operator-manifests"

	profileName     = "profile-source"
	profileRevision = "profile-revision"
)

// FLUX V1 DEPRECATION NOTICE. https://github.com/weaveworks/eksctl/issues/2963
// AddCommonFlagsForFlux configures the flags required to install Flux on an
// EKS cluster and have it point to the specified Git repository.
func AddCommonFlagsForFlux(fs *pflag.FlagSet, opts *api.Git) {
	AddCommonFlagsForGitRepo(fs, opts.Repo)

	fs.StringSliceVar(&opts.Repo.Paths, gitPaths, []string{},
		"Relative paths within the Git repo for Flux to locate Kubernetes manifests")
	fs.StringVar(&opts.Operator.Label, gitLabel, "flux",
		"Git label to keep track of Flux's sync progress; this is equivalent to overriding --git-sync-tag and --git-notes-ref in Flux")
	fs.StringVar(&opts.Repo.FluxPath, gitFluxPath, "flux/",
		"Directory within the Git repository where to commit the Flux manifests")
	fs.StringVar(&opts.Operator.Namespace, namespace, "flux",
		"Cluster namespace where to install Flux and the Helm Operator")
	fs.BoolVar(&opts.Operator.ReadOnly, readOnly, false,
		"Configure Flux in read-only mode and create the deploy key as read-only (Github only)")
	opts.Operator.CommitOperatorManifests = fs.Bool(commitOperatorManifests, true,
		"Commit and push Flux manifests to the Git repo on install")
	opts.Operator.WithHelm = fs.Bool(withHelm, true, "Install the Helm Operator")
	fs.StringSliceVar(&opts.Operator.AdditionalFluxArgs, additionalFluxArgs, []string{},
		"Additional command line arguments for the Flux daemon")
	fs.StringSliceVar(&opts.Operator.AdditionalHelmOperatorArgs, additionalHelmOperatorArgs, []string{},
		"Additional command line arguments for the Helm Operator")
}

// FLUX V1 DEPRECATION NOTICE. https://github.com/weaveworks/eksctl/issues/2963
// AddCommonFlagsForGitRepo configures the flags required to interact with a Git
// repository.
func AddCommonFlagsForGitRepo(fs *pflag.FlagSet, repo *api.Repo) {
	fs.StringVar(&repo.URL, gitURL, "",
		"SSH URL of the Git repository to be used for GitOps, e.g. git@github.com:<github_org>/<repo_name>")
	fs.StringVar(&repo.Branch, gitBranch, "master",
		"Git branch to be used for GitOps")
	fs.StringVar(&repo.User, gitUser, "Flux",
		"Username to use as Git committer")
	fs.StringVar(&repo.Email, gitEmail, "",
		"Email to use as Git committer")
	fs.StringVar(&repo.PrivateSSHKeyPath,
		gitPrivateSSHKeyPath, "",
		"Optional path to the private SSH key to use with Git, e.g. ~/.ssh/id_rsa")
}

// AddCommonFlagsForProfile configures the flags required to enable a Quick
// Start profile.
func AddCommonFlagsForProfile(fs *pflag.FlagSet, opts *api.Profile) {
	fs.StringVarP(&opts.Source, profileName, "", "", "name or URL of the Quick Start profile. For example, app-dev.")
	fs.StringVarP(&opts.Revision, profileRevision, "", "master", "revision of the Quick Start profile.")
}

// FLUX V1 DEPRECATION NOTICE. https://github.com/weaveworks/eksctl/issues/2963
// GitConfigLoader handles loading of ClusterConfigFile v.s. using CLI
// flags for Git-related commands.
type GitConfigLoader struct {
	cmd                                *Cmd
	flagsIncompatibleWithConfigFile    sets.String
	flagsIncompatibleWithoutConfigFile sets.String
	validateWithConfigFile             func() error
	validateWithoutConfigFile          func() error
	gitConfig                          *api.Git
}

// FLUX V1 DEPRECATION NOTICE. https://github.com/weaveworks/eksctl/issues/2963
// NewGitConfigLoader creates a new ClusterConfigLoader which handles
// loading of ClusterConfigFile v.s. using CLI flags for Git-related
// commands.
func NewGitConfigLoader(cmd *Cmd, cfg *api.Git) *GitConfigLoader {
	l := &GitConfigLoader{
		cmd: cmd,
		flagsIncompatibleWithConfigFile: sets.NewString(
			"region",
			"version",
			"cluster",
			gitURL,
			gitBranch,
			gitUser,
			gitEmail,
			gitPrivateSSHKeyPath,
			gitPaths,
			gitLabel,
			gitFluxPath,
			namespace,
			withHelm,
			"amend",
			profileName,
			profileRevision,
			additionalFluxArgs,
			additionalHelmOperatorArgs,
		),
		flagsIncompatibleWithoutConfigFile: sets.NewString(),
	}

	l.gitConfig = cfg

	l.validateWithoutConfigFile = func() error {
		meta := l.cmd.ClusterConfig.Metadata
		if meta.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}
		if meta.Region == "" {
			return ErrMustBeSet("--region")
		}

		return nil
	}

	l.validateWithConfigFile = func() error {
		meta := l.cmd.ClusterConfig.Metadata
		if meta.Name == "" {
			return ErrMustBeSet("metadata.name")
		}

		if meta.Region == "" {
			return ErrMustBeSet("metadata.region")
		}

		if l.cmd.ClusterConfig.Git == nil {
			return nil
		}

		if l.cmd.ClusterConfig.Git.Repo != nil {
			if err := git.ValidatePrivateSSHKeyPath(l.cmd.ClusterConfig.Git.Repo.PrivateSSHKeyPath); err != nil {
				return errors.Wrapf(err, "please supply a valid file for git.repo.privateSSHKeyPath")
			}
		}

		return nil
	}

	return l
}

// FLUX V1 DEPRECATION NOTICE. https://github.com/weaveworks/eksctl/issues/2963
// WithRepoValidation adds extra validation to make sure that the git url and the email are provided as they
// are required for the commands enable profile and enable repo (but not for generate profile)
func (l *GitConfigLoader) WithRepoValidation() *GitConfigLoader {
	newLoader := *l
	newLoader.validateWithoutConfigFile = func() error {
		if newLoader.cmd.ClusterConfig.Git.Repo.URL == "" {
			return ErrMustBeSet("--git-url")
		}

		if newLoader.cmd.ClusterConfig.Git.Repo.Email == "" {
			return ErrMustBeSet("--git-email")
		}
		if err := git.ValidateURL(newLoader.cmd.ClusterConfig.Git.Repo.URL); err != nil {
			return errors.Wrapf(err, "please supply a valid URL for --%s argument", gitURL)
		}

		return l.validateWithoutConfigFile()
	}

	newLoader.validateWithConfigFile = func() error {
		repo := newLoader.cmd.ClusterConfig.Git.Repo
		if repo == nil || repo.URL == "" {
			return ErrMustBeSet("git.repo.url")
		}

		if repo.Email == "" {
			return ErrMustBeSet("git.repo.email")
		}

		if l.cmd.ClusterConfig.GitOps != nil {
			return errors.New("config cannot be provided for gitops alongside git")
		}

		return l.validateWithConfigFile()
	}
	return &newLoader
}

// FLUX V1 DEPRECATION NOTICE. https://github.com/weaveworks/eksctl/issues/2963
// WithProfileValidation adds extra validation to make sure that the git url and the email are provided as they
// are required for the commands enable profile and enable repo (but not for generate profile)
func (l *GitConfigLoader) WithProfileValidation() *GitConfigLoader {
	newLoader := *l
	newLoader.validateWithoutConfigFile = func() error {
		if !newLoader.cmd.ClusterConfig.HasBootstrapProfile() {
			return ErrMustBeSet("--profile-source")
		}

		return l.validateWithoutConfigFile()
	}

	newLoader.validateWithConfigFile = func() error {
		if !newLoader.cmd.ClusterConfig.HasBootstrapProfile() {
			return ErrMustBeSet("git.bootstrapProfile.Source")
		}

		return l.validateWithConfigFile()
	}
	return &newLoader
}

// FLUX V1 DEPRECATION NOTICE. https://github.com/weaveworks/eksctl/issues/2963
// Load ClusterConfig or use CLI flags.
func (l *GitConfigLoader) Load() error {
	if err := api.Register(); err != nil {
		return err
	}

	logger.Warning("the `enable repo`, `enable profile` and `generate profile` commands DEPRECATED: see https://github.com/weaveworks/eksctl/issues/2963")
	logger.Warning("the `git.repo`, `git.operator` and `git.bootstrapProfile` config options are DEPRECATED: see https://github.com/weaveworks/eksctl/issues/2963")

	if l.cmd.ClusterConfigFile == "" {
		l.cmd.ClusterConfig.Metadata.Region = l.cmd.ProviderConfig.Region
		for f := range l.flagsIncompatibleWithoutConfigFile {
			if flag := l.cmd.CobraCommand.Flag(f); flag != nil && flag.Changed {
				return fmt.Errorf("cannot use --%s unless a config file is specified via --config-file/-f", f)
			}
		}

		l.cmd.ClusterConfig.Git = l.gitConfig
		if l.cmd.NameArg != "" && l.cmd.ClusterConfig.Git.BootstrapProfile.Source != "" {
			return ErrFlagAndArg("--profile-source", l.cmd.ClusterConfig.Git.BootstrapProfile.Source, l.cmd.NameArg)
		}
		if l.cmd.NameArg != "" {
			l.cmd.ClusterConfig.Git.BootstrapProfile.Source = l.cmd.NameArg
		}
		api.SetDefaultGitSettings(l.cmd.ClusterConfig)
		return l.validateWithoutConfigFile()
	}

	// The reference to ClusterConfig should only be reassigned if ClusterConfigFile is specified
	// because other parts of the code store the pointer locally and access it directly instead of via
	// the Cmd reference
	var err error
	if l.cmd.ClusterConfig, err = eks.LoadConfigFromFile(l.cmd.ClusterConfigFile); err != nil {
		return err
	}

	meta := l.cmd.ClusterConfig.Metadata
	if meta == nil {
		return ErrMustBeSet("metadata")
	}

	for f := range l.flagsIncompatibleWithConfigFile {
		if flag := l.cmd.CobraCommand.Flag(f); flag != nil && flag.Changed {
			return ErrCannotUseWithConfigFile(fmt.Sprintf("--%s", f))
		}
	}

	if meta.Region != "" {
		l.cmd.ProviderConfig.Region = meta.Region
	}

	api.SetDefaultGitSettings(l.cmd.ClusterConfig)
	return l.validateWithConfigFile()
}

// GitOpsConfigLoader handles loading of ClusterConfigFile v.s. using CLI
// flags for GitOps-related commands.
type GitOpsConfigLoader struct {
	cmd                    *Cmd
	validateWithConfigFile func() error
}

// NewGitOpsConfigLoader creates a new ClusterConfigLoader which handles
// loading of ClusterConfigFile GitOps-related commands.
func NewGitOpsConfigLoader(cmd *Cmd) *GitOpsConfigLoader {
	l := &GitOpsConfigLoader{
		cmd: cmd,
	}

	l.validateWithConfigFile = func() error {
		meta := l.cmd.ClusterConfig.Metadata
		if meta.Name == "" {
			return ErrMustBeSet("metadata.name")
		}

		if meta.Region == "" {
			return ErrMustBeSet("metadata.region")
		}

		if l.cmd.ClusterConfig.Git != nil {
			return errors.New("config cannot be provided for git.repo, git.bootstrapProfile or git.operator alongside gitops.*")
		}

		if l.cmd.ClusterConfig.GitOps == nil || l.cmd.ClusterConfig.GitOps.Flux == nil {
			return ErrMustBeSet("gitops.flux")
		}

		fluxCfg := l.cmd.ClusterConfig.GitOps.Flux
		if fluxCfg.GitProvider == "" {
			return ErrMustBeSet("gitops.flux.gitProvider")
		}

		if len(fluxCfg.Flags) == 0 {
			return ErrMustBeSet("gitops.flux.flags")
		}

		return nil
	}

	return l
}

// Load ClusterConfig or use CLI flags.
func (l *GitOpsConfigLoader) Load() error {
	if err := api.Register(); err != nil {
		return err
	}

	if l.cmd.ClusterConfigFile == "" {
		return ErrMustBeSet("--config-file/-f <file>")
	}

	// The reference to ClusterConfig should only be reassigned if ClusterConfigFile is specified
	// because other parts of the code store the pointer locally and access it directly instead of via
	// the Cmd reference
	var err error
	if l.cmd.ClusterConfig, err = eks.LoadConfigFromFile(l.cmd.ClusterConfigFile); err != nil {
		return err
	}

	meta := l.cmd.ClusterConfig.Metadata
	if meta == nil {
		return ErrMustBeSet("metadata")
	}

	if meta.Region != "" {
		l.cmd.ProviderConfig.Region = meta.Region
	}

	return l.validateWithConfigFile()
}
