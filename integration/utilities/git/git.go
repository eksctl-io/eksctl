//go:build integration
// +build integration

package git

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/executor"
	"github.com/weaveworks/eksctl/pkg/utils/file"
)

const (
	// Repository is the default testing Git repository.
	Repository = "git@github.com:eksctl-bot/my-gitops-repo.git"
	// Email is the default testing Git email.
	Email = "eksctl-bot@weave.works"
)

// CreateBranch creates the provided branch.
func CreateBranch(branch string) (string, error) {
	cli := NewGitClient(ClientParams{})
	cloneDir, err := cli.CloneRepoInTmpDir(
		"eksctl-install-flux-test-branch-",
		CloneOptions{
			URL:       Repository,
			Branch:    branch,
			Bootstrap: true,
		},
	)
	if err != nil {
		return "", err
	}
	if err := cli.Push(); err != nil {
		return "", err
	}
	return cloneDir, nil
}

// CleanupBranchAndRepo deletes the local clone used for testing Git repository, and delete the branch from "origin" as well.
func CleanupBranchAndRepo(branch, cloneDir string) error {
	cli := NewGitClient(ClientParams{})
	cli.WithDir(cloneDir)
	if err := cli.DeleteRemoteBranch(branch); err != nil {
		return err
	}
	return cli.DeleteLocalRepo()
}

// GetBranch clones the testing Git repository and checks out the provided branch.
func GetBranch(branch string) (string, error) {
	cli := NewGitClient(ClientParams{})
	return cli.CloneRepoInTmpDir(
		"eksctl-install-flux-test-branch-",
		CloneOptions{
			URL:       Repository,
			Branch:    branch,
			Bootstrap: true,
		},
	)
}

// TmpCloner can clone git repositories in temporary directories
type TmpCloner interface {
	CloneRepoInTmpDir(cloneDirPrefix string, options CloneOptions) (string, error)
}

// Client can perform git operations on the given directory
type Client struct {
	executor executor.Executor
	dir      string
}

// ClientParams groups the arguments to provide to create a new Git client.
type ClientParams struct {
	PrivateSSHKeyPath string
}

// ValidatePrivateSSHKeyPath validates the path to the (optional) private SSH
// key used to interact with the Git repository configured in this object.
func ValidatePrivateSSHKeyPath(privateSSHKeyPath string) error {
	if privateSSHKeyPath != "" && !file.Exists(privateSSHKeyPath) {
		return fmt.Errorf("invalid path to private SSH key: %s", privateSSHKeyPath)
	}
	return nil
}

// NewGitClient returns a client that can perform git operations
func NewGitClient(params ClientParams) *Client {
	return &Client{
		executor: executor.NewShellExecutor(envVars(params)),
	}
}

func envVars(params ClientParams) executor.EnvVars {
	envVars := executor.EnvVars{
		"PATH": os.Getenv("PATH"),
	}
	if sshAuthSock, ok := os.LookupEnv("SSH_AUTH_SOCK"); ok {
		envVars["SSH_AUTH_SOCK"] = sshAuthSock
	}
	if params.PrivateSSHKeyPath != "" {
		envVars["GIT_SSH_COMMAND"] = "ssh -i " + params.PrivateSSHKeyPath
	}
	return envVars
}

// NewGitClientFromExecutor returns a client that can have an executor injected. Useful for testing
func NewGitClientFromExecutor(executor executor.Executor) *Client {
	return &Client{
		executor: executor,
	}
}

// CloneOptions are the options for cloning a Git repository
type CloneOptions struct {
	URL       string
	Branch    string
	Bootstrap bool // create the branch if the repository is empty
}

// WithDir directly sets the Client to use a directory, without havine to clone
// it
func (git *Client) WithDir(dir string) {
	git.dir = dir
}

// CloneRepoInTmpDir clones a repo specified in the gitURL in a temporary directory and checks out the specified branch
func (git *Client) CloneRepoInTmpDir(tmpDirPrefix string, options CloneOptions) (string, error) {
	cloneDir, err := os.MkdirTemp(os.TempDir(), tmpDirPrefix)
	if err != nil {
		return "", fmt.Errorf("cannot create temporary directory: %s", err)
	}
	return cloneDir, git.cloneRepoInPath(cloneDir, options)
}

// CloneRepoInPath behaves like CloneRepoInTmpDir but clones the repository in a specific directory
// which creates if needed
func (git *Client) CloneRepoInPath(clonePath string, options CloneOptions) error {
	if err := os.MkdirAll(clonePath, 0700); err != nil {
		return fmt.Errorf("unable to create directory for cloning: %w", err)
	}
	return git.cloneRepoInPath(clonePath, options)
}

func (git *Client) cloneRepoInPath(clonePath string, options CloneOptions) error {
	args := []string{"clone", options.URL, clonePath}
	if err := git.runGitCmd(args...); err != nil {
		return err
	}
	// Set the working directory to the cloned directory, but
	// only do it after the clone so that it doesn't create an
	// undesirable nested directory
	git.dir = clonePath

	if options.Branch != "" {
		// Switch to target branch
		args := []string{"checkout", options.Branch}
		if options.Bootstrap {
			if !git.isRemoteBranch(options.Branch) {
				args = []string{"checkout", "-b", options.Branch}
			}
		}
		if err := git.runGitCmd(args...); err != nil {
			return err
		}
	}

	return nil
}

func (git *Client) isRemoteBranch(branch string) bool {
	err := git.runGitCmd("ls-remote", "--heads", "--exit-code", "origin", branch)
	return err == nil
}

// Push pushes the changes to the origin remote
func (git Client) Push() error {
	if err := git.runGitCmd("config", "push.default", "current"); err != nil {
		return err
	}
	err := git.runGitCmd("push")
	return err
}

func (git Client) DeleteRemoteBranch(branch string) error {
	return git.runGitCmd("push", "origin", "--delete", branch)
}

// DeleteLocalRepo deletes the local copy of a repository, including the directory
func (git Client) DeleteLocalRepo() error {
	if git.dir != "" {
		return os.RemoveAll(git.dir)
	}
	return fmt.Errorf("no cloned directory to delete")
}

func (git Client) runGitCmd(args ...string) error {
	logger.Debug(fmt.Sprintf("running git %v in %s", args, git.dir))
	return git.executor.ExecInDir("git", git.dir, args...)
}
