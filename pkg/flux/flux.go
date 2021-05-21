package flux

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/executor"
)

const fluxBin = "flux"

type Client struct {
	executor executor.Executor
	opts     *api.Flux
}

func NewClient(opts *api.Flux) (*Client, error) {
	env, err := setTokenEnv(opts.AuthTokenPath, opts.GitProvider)
	if err != nil {
		return nil, err
	}

	return &Client{
		executor: executor.NewShellExecutor(env),
		opts:     opts,
	}, nil
}

func (c *Client) PreFlight() error {
	if _, err := exec.LookPath(fluxBin); err != nil {
		logger.Warning(err.Error())
		return errors.New("flux not found, required")
	}

	args := []string{"check", "--pre"}

	if c.opts.Kubeconfig != "" {
		args = append(args, "--kubeconfig", c.opts.Kubeconfig)
	}

	return c.runFluxCmd(args...)
}

func (c *Client) Bootstrap() error {
	args := []string{"bootstrap", c.opts.GitProvider, "--repository", c.opts.Repository, "--owner", c.opts.Owner}

	if c.opts.Personal {
		args = append(args, "--personal")
	}
	if c.opts.Path != "" {
		args = append(args, "--path", c.opts.Path)
	}
	if c.opts.Branch != "" {
		args = append(args, "--branch", c.opts.Branch)
	}
	if c.opts.Namespace != "" {
		args = append(args, "--namespace", c.opts.Namespace)
	}
	if c.opts.Kubeconfig != "" {
		args = append(args, "--kubeconfig", c.opts.Kubeconfig)
	}

	return c.runFluxCmd(args...)
}

func (c *Client) runFluxCmd(args ...string) error {
	logger.Debug(fmt.Sprintf("running flux %v ", args))
	return c.executor.Exec(fluxBin, args...)
}

func setTokenEnv(tokenPath, gitProvider string) (executor.EnvVars, error) {
	if tokenPath == "" {
		return executor.EnvVars{}, nil
	}

	var token string
	data, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return executor.EnvVars{}, fmt.Errorf("reading auth token file %w", err)
	}

	token = strings.Replace(string(data), "\n", "", -1)
	envVars := executor.EnvVars{}
	switch gitProvider {
	case "github":
		envVars["GITHUB_TOKEN"] = token
	case "gitlab":
		envVars["GITLAB_TOKEN"] = token
	}

	return envVars, nil
}
