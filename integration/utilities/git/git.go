//go:build integration
// +build integration

package git

import (
	"github.com/weaveworks/eksctl/pkg/git"
)

const (
	// Repository is the default testing Git repository.
	Repository = "git@github.com:eksctl-bot/my-gitops-repo.git"
	// Email is the default testing Git email.
	Email = "eksctl-bot@weave.works"
)

// CreateBranch creates the provided branch.
func CreateBranch(branch string) (string, error) {
	cli := git.NewGitClient(git.ClientParams{})
	cloneDir, err := cli.CloneRepoInTmpDir(
		"eksctl-install-flux-test-branch-",
		git.CloneOptions{
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
	cli := git.NewGitClient(git.ClientParams{})
	cli.WithDir(cloneDir)
	if err := cli.DeleteRemoteBranch(branch); err != nil {
		return err
	}
	return cli.DeleteLocalRepo()
}

// GetBranch clones the testing Git repository and checks out the provided branch.
func GetBranch(branch string) (string, error) {
	cli := git.NewGitClient(git.ClientParams{})
	return cli.CloneRepoInTmpDir(
		"eksctl-install-flux-test-branch-",
		git.CloneOptions{
			URL:       Repository,
			Branch:    branch,
			Bootstrap: true,
		},
	)
}
