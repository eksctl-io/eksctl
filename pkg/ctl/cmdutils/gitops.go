package cmdutils

import (
	"fmt"

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

	gitPaths    = "git-paths"
	gitFluxPath = "git-flux-subdir"
	gitLabel    = "git-label"
	namespace   = "namespace"
	withHelm    = "with-helm"

	profileName     = "profile-source"
	profileRevision = "profile-revision"
)

// Options holds options to interact with a Git repository.
type Options struct {
	URL               string
	Branch            string
	User              string
	Email             string
	PrivateSSHKeyPath string
}

// AddCommonFlagsForFlux configures the flags required to install Flux on an
// EKS cluster and have it point to the specified Git repository.
func AddCommonFlagsForFlux(fs *pflag.FlagSet, opts *api.Git) {
	AddCommonFlagsForGit(fs, opts.Repo)

	fs.StringSliceVar(&opts.Repo.Paths, gitPaths, []string{},
		"Relative paths within the Git repo for Flux to locate Kubernetes manifests")
	fs.StringVar(&opts.Operator.Label, gitLabel, "flux",
		"Git label to keep track of Flux's sync progress; this is equivalent to overriding --git-sync-tag and --git-notes-ref in Flux")
	fs.StringVar(&opts.Repo.FluxPath, gitFluxPath, "flux/",
		"Directory within the Git repository where to commit the Flux manifests")
	fs.StringVar(&opts.Operator.Namespace, namespace, "flux",
		"Cluster namespace where to install Flux, the Helm Operator and Tiller")
	opts.Operator.WithHelm = fs.Bool(withHelm, true, "Install the Helm Operator and Tiller")
}

// AddCommonFlagsForGit configures the flags required to interact with a Git
// repository.
func AddCommonFlagsForGit(fs *pflag.FlagSet, repo *api.Repo) {
	fs.StringVar(&repo.URL, gitURL, "",
		"SSH URL of the Git repository to be used for GitOps, e.g. git@github.com:<github_org>/<repo_name>")
	fs.StringVar(&repo.Branch, gitBranch, "master",
		"Git branch to be used for GitOps")
	fs.StringVar(&repo.User, gitUser, "Flux",
		"Username to use as Git committer")
	fs.StringVar(&repo.Email, gitEmail, "",
		"Email to use as Git committer")
	fs.StringVar(&repo.PrivateSSHKeyPath, gitPrivateSSHKeyPath, "",
		"Optional path to the private SSH key to use with Git, e.g. ~/.ssh/id_rsa")
}

// AddCommonFlagsForProfile configures the flags required to enable a Quick
// Start profile.
func AddCommonFlagsForProfile(fs *pflag.FlagSet, opts *api.Profile) {
	fs.StringVarP(&opts.Source, profileName, "", "", "name or URL of the Quick Start profile. For example, app-dev.")
	fs.StringVarP(&opts.Revision, profileRevision, "", "master", "revision of the Quick Start profile.")
}

// GitOpsConfigLoader handles loading of ClusterConfigFile v.s. using CLI
// flags for GitOps-related commands.
type GitOpsConfigLoader struct {
	cmd                                *Cmd
	flagsIncompatibleWithConfigFile    sets.String
	flagsIncompatibleWithoutConfigFile sets.String
	validateWithConfigFile             func() error
	validateWithoutConfigFile          func() error
	gitConfig                          *api.Git
}

// NewGitOpsConfigLoader creates a new ClusterConfigLoader which handles
// loading of ClusterConfigFile v.s. using CLI flags for GitOps-related
// commands.
func NewGitOpsConfigLoader(cmd *Cmd, cfg *api.Git) *GitOpsConfigLoader {
	l := &GitOpsConfigLoader{
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

// WithRepoValidation adds extra validation to make sure that the git url and the email are provided as they
// are required for the commands enable profile and enable repo (but not for generate profile)
func (l *GitOpsConfigLoader) WithRepoValidation() *GitOpsConfigLoader {
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
			return ErrMustBeSet("git.repo.URL")
		}

		if repo.Email == "" {
			return ErrMustBeSet("git.repo.email")
		}
		return l.validateWithConfigFile()
	}
	return &newLoader
}

// WithProfileValidation adds extra validation to make sure that the git url and the email are provided as they
// are required for the commands enable profile and enable repo (but not for generate profile)
func (l *GitOpsConfigLoader) WithProfileValidation() *GitOpsConfigLoader {
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

// Load ClusterConfig or use CLI flags.
func (l *GitOpsConfigLoader) Load() error {
	if err := api.Register(); err != nil {
		return err
	}

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

	var err error

	// The reference to ClusterConfig should only be reassigned if ClusterConfigFile is specified
	// because other parts of the code store the pointer locally and access it directly instead of via
	// the Cmd reference
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
