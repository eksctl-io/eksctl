package flux

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/executor"
)

const (
	fluxBin             = "flux"
	minSupportedVersion = ">= 0.13.3"
)

type Client struct {
	executor executor.Executor
	opts     *api.Flux
}

func NewClient(opts *api.Flux) (*Client, error) {
	return &Client{
		executor: executor.NewShellExecutor(executor.EnvVars{}),
		opts:     opts,
	}, nil
}

func (c *Client) PreFlight() error {
	if _, err := exec.LookPath(fluxBin); err != nil {
		logger.Warning(err.Error())
		return errors.New("flux not found, required")
	}

	if err := c.checkFluxVersion(); err != nil {
		return err
	}

	args := []string{"check", "--pre"}
	for k, v := range c.opts.Flags {
		if k == "kubeconfig" || k == "context" {
			args = append(args, fmt.Sprintf("--%s", k), v)
		}
	}

	return c.runFluxCmd(args...)
}

func (c *Client) Bootstrap() error {
	args := []string{"bootstrap", c.opts.GitProvider}

	for k, v := range c.opts.Flags {
		args = append(args, fmt.Sprintf("--%s", k), v)
	}

	return c.runFluxCmd(args...)
}

func (c *Client) runFluxCmd(args ...string) error {
	logger.Debug(fmt.Sprintf("running flux %v ", args))
	return c.executor.Exec(fluxBin, args...)
}

func (c *Client) checkFluxVersion() error {
	logger.Debug(fmt.Sprintf("checking flux version is %s", minSupportedVersion))
	out, err := c.executor.ExecWithOut(fluxBin, "--version")
	if err != nil {
		return err
	}

	trimmed := strings.TrimRight(string(out), "\n")
	parts := strings.Split(trimmed, " ")
	if len(parts) < 3 {
		return fmt.Errorf("unexpected format returned from 'flux --version': %s", parts)
	}

	v, err := semver.NewVersion(parts[2])
	if err != nil {
		return err
	}

	constraint, err := semver.NewConstraint(minSupportedVersion)
	if err != nil {
		return err
	}

	withinBounds := constraint.Check(v)
	if !withinBounds {
		return fmt.Errorf("found flux version 0.13.2, eksctl requires %s", minSupportedVersion)
	}

	return nil
}
