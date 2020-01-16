package cmdutils

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/gitops/flux"
	"github.com/weaveworks/eksctl/pkg/gitops/profile"
	"k8s.io/apimachinery/pkg/util/sets"
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

	profileName     = "name"
	profileRevision = "revision"
)

// AddCommonFlagsForFlux configures the flags required to install Flux on an
// EKS cluster and have it point to the specified Git repository.
func AddCommonFlagsForFlux(fs *pflag.FlagSet, opts *flux.InstallOpts) {
	AddCommonFlagsForGit(fs, &opts.GitOptions)

	fs.StringSliceVar(&opts.GitPaths, gitPaths, []string{},
		"Relative paths within the Git repo for Flux to locate Kubernetes manifests")
	fs.StringVar(&opts.GitLabel, gitLabel, "flux",
		"Git label to keep track of Flux's sync progress; this is equivalent to overriding --git-sync-tag and --git-notes-ref in Flux")
	fs.StringVar(&opts.GitFluxPath, gitFluxPath, "flux/",
		"Directory within the Git repository where to commit the Flux manifests")
	fs.StringVar(&opts.Namespace, namespace, "flux",
		"Cluster namespace where to install Flux, the Helm Operator and Tiller")
	fs.BoolVar(&opts.WithHelm, withHelm, true,
		"Install the Helm Operator and Tiller")
	fs.BoolVar(&opts.Amend, "amend", false,
		"Stop to manually tweak the Flux manifests before pushing them to the Git repository")
}

// AddCommonFlagsForGit configures the flags required to interact with a Git
// repository.
func AddCommonFlagsForGit(fs *pflag.FlagSet, opts *git.Options) {
	fs.StringVar(&opts.URL, gitURL, "",
		"SSH URL of the Git repository to be used for GitOps, e.g. git@github.com:<github_org>/<repo_name>")
	fs.StringVar(&opts.Branch, gitBranch, "master",
		"Git branch to be used for GitOps")
	fs.StringVar(&opts.User, gitUser, "Flux",
		"Username to use as Git committer")
	fs.StringVar(&opts.Email, gitEmail, "",
		"Email to use as Git committer")
	fs.StringVar(&opts.PrivateSSHKeyPath, gitPrivateSSHKeyPath, "",
		"Optional path to the private SSH key to use with Git, e.g. ~/.ssh/id_rsa")
	_ = cobra.MarkFlagRequired(fs, gitURL)
	_ = cobra.MarkFlagRequired(fs, gitEmail)
}

var validateErrs *multierror.Error

// ValidateGitOptions validates the provided Git options.
func ValidateGitOptions(opts *git.Options) error {
	if err := opts.ValidateURL(); err != nil {
		return errors.Wrapf(err, "please supply a valid --%s argument", gitURL)
	}
	if err := opts.ValidateEmail(); err != nil {
		return fmt.Errorf("please supply a valid --%s argument", gitEmail)
	}
	if err := opts.ValidatePrivateSSHKeyPath(); err != nil {
		return errors.Wrapf(err, "please supply a valid --%s argument", gitPrivateSSHKeyPath)
	}
	return nil
}

// AddCommonFlagsForProfile configures the flags required to enable a Quick
// Start profile.
func AddCommonFlagsForProfile(fs *pflag.FlagSet, opts *profile.Options) {
	fs.StringVarP(&opts.Name, profileName, "", "", "name or URL of the Quick Start profile. For example, app-dev.")
	fs.StringVarP(&opts.Revision, profileRevision, "", "master", "revision of the Quick Start profile.")
}

// gitOpsConfigLoader handles loading of ClusterConfigFile v.s. using CLI
// flags for GitOps-related commands.
type gitOpsConfigLoader struct {
	cmd                                *Cmd
	flagsIncompatibleWithConfigFile    sets.String
	flagsIncompatibleWithoutConfigFile sets.String
	validateWithConfigFile             func() error
	validateWithoutConfigFile          func() error
}

// NewGitOpsConfigLoader creates a new ClusterConfigLoader which handles
// loading of ClusterConfigFile v.s. using CLI flags for GitOps-related
// commands.
func NewGitOpsConfigLoader(cmd *Cmd) ClusterConfigLoader {
	l := &gitOpsConfigLoader{
		cmd: cmd,
		flagsIncompatibleWithConfigFile: sets.NewString(
			"region",
			"version",
			"cluster",
		),
		flagsIncompatibleWithoutConfigFile: sets.NewString(),
	}

	l.validateWithoutConfigFile = func() error {
		meta := l.cmd.ClusterConfig.Metadata
		if meta.Name == "" {
			validateErrs = multierror.Append(validateErrs, ErrMustBeSet(ClusterNameFlag(cmd)))
		}
		if meta.Region == "" {
			validateErrs = multierror.Append(validateErrs, ErrMustBeSet("--region"))
		}
		return validateErrs.ErrorOrNil()
	}

	l.validateWithConfigFile = func() error {
		meta := l.cmd.ClusterConfig.Metadata
		if meta.Name == "" {
			validateErrs = multierror.Append(validateErrs, ErrMustBeSet("metadata.name"))
		}
		if meta.Region == "" {
			validateErrs = multierror.Append(validateErrs, ErrMustBeSet("metadata.region"))
		}
		return validateErrs.ErrorOrNil()
	}

	return l
}

// Load ClusterConfig or use CLI flags.
func (l *gitOpsConfigLoader) Load() error {
	if err := api.Register(); err != nil {
		return err
	}

	if l.cmd.ClusterConfigFile == "" {
		l.cmd.ClusterConfig.Metadata.Region = l.cmd.ProviderConfig.Region
		for f := range l.flagsIncompatibleWithoutConfigFile {
			if flag := l.cmd.CobraCommand.Flag(f); flag != nil && flag.Changed {
				multiErr = multierror.Append(fmt.Errorf("cannot use --%s unless a config file is specified via --config-file/-f", f))
				return multiErr.ErrorOrNil()
			}
		}
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
		multiErr = multierror.Append(ErrMustBeSet("metadata"))
	}

	for f := range l.flagsIncompatibleWithConfigFile {
		if flag := l.cmd.CobraCommand.Flag(f); flag != nil && flag.Changed {
			multiErr = multierror.Append(ErrCannotUseWithConfigFile(fmt.Sprintf("--%s", f)))
		}
	}

	if meta.Region != "" {
		l.cmd.ProviderConfig.Region = meta.Region
	}

	if multiErr != nil {
		return multiErr
	}

	return l.validateWithConfigFile()
}
