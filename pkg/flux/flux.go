package flux

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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

func NewClient(opts *api.Flux) *Client {
	return &Client{
		executor: executor.NewShellExecutor(setEnvVars(opts.AuthTokenPath, opts.GitProvider)),
		opts:     opts,
	}
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

func setEnvVars(tokenPath, gitProvider string) executor.EnvVars {
	envVars := executor.EnvVars{
		"PATH":                  os.Getenv("PATH"),
		"HOME":                  os.Getenv("HOME"),
		"AWS_ACCESS_KEY_ID":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"AWS_SECRET_ACCESS_KEY": os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	var token string
	if tokenPath != "" {
		data, err := ioutil.ReadFile(tokenPath)
		if err != nil {
			logger.Warning("reading auth token file %s", err)
		}

		token = strings.Replace(string(data), "\n", "", -1)
	}

	switch gitProvider {
	case "github":
		if token == "" {
			if githubToken, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
				token = githubToken
			}
		}
		envVars["GITHUB_TOKEN"] = token
	case "gitlab":
		if token == "" {
			if gitlabToken, ok := os.LookupEnv("GITLAB_TOKEN"); ok {
				token = gitlabToken
			}
		}
		envVars["GITLAB_TOKEN"] = token
	}

	return envVars
}
