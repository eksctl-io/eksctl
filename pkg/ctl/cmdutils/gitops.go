package cmdutils

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/git"
)

const (
	gitURL               = "git-url"
	gitBranch            = "git-branch"
	gitUser              = "git-user"
	gitEmail             = "git-email"
	gitPrivateSSHKeyPath = "git-private-ssh-key-path"
)

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
