package git

import (
	"context"
	"fmt"
	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/git/executor"
	"os"
	"os/exec"
	"time"
)

// Client can perform git operations on the given directory
type Client struct {
	executor executor.Executor
	dir      string
	user     string
	email    string
}

func NewGitClient(cloneDir string, user string, email string, ctx context.Context, timeout time.Duration) *Client {
	return &Client{
		executor: executor.NewShellExecutor(ctx, timeout),
		dir:      cloneDir,
		user:     user,
		email:    email,
	}
}

func NewGitClientFromExecutor(cloneDir string, user string, email string, executor executor.Executor) *Client {
	return &Client{
		executor: executor,
		dir:      cloneDir,
		user:     user,
		email:    email,
	}
}

// CloneRepo clones a repo specified in the gitURL and checks out the specified branch
func (git *Client) CloneRepo(branch string, gitURL string) (string, error) {
	// TODO TODO TODO create directory if needed. Fail if given dir is not empty
	// FIXME
	//cloneDir, err := ioutil.TempDir(os.TempDir(), "eksctl-install-flux-clone")
	//if err != nil {
	//	return "", fmt.Errorf("cannot create temporary directory: %s", err)
	//}
	args := []string{"clone", "-b", branch, gitURL, git.dir}
	err := git.runGitCmd(args...)
	return git.dir, err
}

func (git Client) Add(files ...string) error {
	args := append([]string{"add", "--"}, files...)
	if err := git.runGitCmd(args...); err != nil {
		return err
	}
	return nil
}

// Commit  makes a commit if there are staged changes
func (git Client) Commit() error {
	// Note, this useed to do runGitCmd(diffCtx, git.dir, "diff", "--cached", "--quiet", "--", fi.opts.gitFluxPath); err == nil {
	if err := git.runGitCmd(git.dir, "diff", "--cached", "--quiet"); err == nil {
		logger.Info("Nothing to commit (the repository contained identical manifests), moving on")
		return nil
	} else if _, ok := err.(*exec.ExitError); !ok {
		return err
	}

	// Commit
	args := []string{"commit",
		"-m", "Add Initial Flux configuration",
		fmt.Sprintf("--author=%s <%s>", git.user, git.email),
	}
	if err := git.runGitCmd(args...); err != nil {
		return err
	}
	return nil
}

// Push pushes the changes to the origin remote
func (git Client) Push() error {
	return git.runGitCmd("push")
}

// DeleteLocalRepo deletes the local copy of a repository, including the directory
func (git Client) DeleteLocalRepo() error {
	return os.RemoveAll(git.dir)
}

func (git Client) runGitCmd(args ...string) error {
	return git.executor.Exec("git", git.dir, args...)
}
