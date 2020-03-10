// +build integration

package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

const (
	// Repository is the default testing Git repository.
	Repository = "git@github.com:eksctl-bot/my-gitops-repo.git"
	// Email is the default testing Git email.
	Email = "eksctl-bot@weave.works"
	// Name is the default cluster name to test against.
	Name = "autoscaler"
	// Region is the default region to test against.
	Region = "ap-northeast-1"
)

// CreateBranch creates the provided branch.
func CreateBranch(branch, privateSSHKeyPath string) (string, error) {
	cloneDir, err := ioutil.TempDir(os.TempDir(), "eksctl-install-flux-test-clone-")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %s", err)
	}
	if err := gitWith(gitParams{Args: []string{"clone", "-b", "master", Repository, cloneDir}, Dir: cloneDir, Env: gitSSHCommand(privateSSHKeyPath)}); err != nil {
		return "", err
	}
	if err := gitWith(gitParams{Args: []string{"checkout", "-b", branch}, Dir: cloneDir, Env: gitSSHCommand(privateSSHKeyPath)}); err != nil {
		return "", err
	}
	if err := gitWith(gitParams{Args: []string{"push", "origin", branch}, Dir: cloneDir, Env: gitSSHCommand(privateSSHKeyPath)}); err != nil {
		return "", err
	}
	return cloneDir, nil
}

// DeleteBranch deletes the local clone used for testing Git repository, and delete the branch from "origin" as well.
func DeleteBranch(branch, cloneDir, privateSSHKeyPath string) error {
	defer os.RemoveAll(cloneDir)
	return gitWith(gitParams{Args: []string{"push", "origin", "--delete", branch}, Dir: cloneDir, Env: gitSSHCommand(privateSSHKeyPath)})
}

// GetBranch clones the testing Git repository and checks out the provided branch.
func GetBranch(branch, privateSSHKeyPath string) (string, error) {
	cloneDir, err := ioutil.TempDir(os.TempDir(), "eksctl-install-flux-test-branch-")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %s", err)
	}
	if err := gitWith(gitParams{Args: []string{"clone", "-b", branch, Repository, cloneDir}, Dir: cloneDir, Env: gitSSHCommand(privateSSHKeyPath)}); err != nil {
		return "", err
	}
	return cloneDir, nil
}

type gitParams struct {
	Args []string
	Env  []string
	Dir  string
}

func gitWith(params gitParams) error {
	gitCmd := exec.Command("git", params.Args...)
	if params.Env != nil {
		gitCmd.Env = params.Env
	}
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	if params.Dir != "" {
		gitCmd.Dir = params.Dir
	}
	return gitCmd.Run()
}

func gitSSHCommand(privateSSHKeyPath string) []string {
	return []string{
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
		fmt.Sprintf("GIT_SSH_COMMAND=ssh -i %s", privateSSHKeyPath),
	}
}
