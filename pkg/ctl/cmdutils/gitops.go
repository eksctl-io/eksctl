package cmdutils

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/gitops/flux"
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
)

// AddCommonFlagsForFlux configures the flags required to install Flux on an
// EKS cluster and have it point to the specified Git repository.
func AddCommonFlagsForFlux(fs *pflag.FlagSet, opts *flux.InstallOpts) {
	AddCommonFlagsForGit(fs, &opts.GitOptions)

	fs.StringSliceVar(&opts.GitPaths, gitPaths, []string{},
		"Relative paths within the Git repo for Flux to locate Kubernetes manifests")
	fs.StringVar(&opts.GitLabel, gitLabel, "flux",
		"Git label to keep track of Flux's sync progress; overrides both --git-sync-tag and --git-notes-ref")
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
}

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
