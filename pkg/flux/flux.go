package flux

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/executor"
	"github.com/weaveworks/logger"
)

const fluxBin = "flux"

type Client struct {
	executor executor.Executor
}

func NewClient(tokenPath, gitProvider string) *Client {
	return &Client{executor.NewShellExecutor(setEnvVars(tokenPath, gitProvider))}
}

func (c *Client) PreFlight() error {
	if _, err := exec.LookPath(fluxBin); err != nil {
		logger.Warning(err.Error())
		return errors.New("flux not found, required")
	}

	return c.runFluxCmd("check", "--pre")
}

func (c *Client) Bootstrap(opts *api.Flux) error {
	args := []string{"bootstrap", opts.GitProvider, "--repository", opts.Repository, "--owner", opts.Owner}

	if opts.Personal {
		args = append(args, "--personal")
	}
	if opts.Path != "" {
		args = append(args, "--path", opts.Path)
	}
	if opts.Branch != "" {
		args = append(args, "--branch", opts.Branch)
	}
	if opts.Namespace != "" {
		args = append(args, "--namespace", opts.Namespace)
	}

	return c.runFluxCmd(args...)
}

func (c *Client) runFluxCmd(args ...string) error {
	logger.Debug(fmt.Sprintf("running flux %v ", args))
	return c.executor.Exec(fluxBin, args...)
}

func setEnvVars(tokenPath, gitProvider string) executor.EnvVars {
	envVars := executor.EnvVars{
		"PATH": os.Getenv("PATH"),
		"HOME": os.Getenv("HOME"),
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
