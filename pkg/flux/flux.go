package flux

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/executor"
)

const (
	fluxBin             = "flux"
	minSupportedVersion = "0.32.0"
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
	logger.Debug(fmt.Sprintf("checking flux version is greater than minimum supported version %s", minSupportedVersion))
	out, err := c.executor.ExecWithOut(fluxBin, "--version")
	if err != nil {
		return err
	}

	trimmed := strings.TrimRight(string(out), "\n")
	parts := strings.Split(trimmed, " ")
	if len(parts) < 3 {
		return fmt.Errorf("unexpected format returned from 'flux --version': %s", parts)
	}

	v, err := version.NewVersion(parts[2])
	if err != nil {
		return fmt.Errorf("failed to parse Flux version %q: %w", parts[2], err)
	}

	supportedVersion, err := version.NewVersion(minSupportedVersion)
	if err != nil {
		return fmt.Errorf("failed to parse supported Flux version %s: %w", minSupportedVersion, err)
	}

	if v.LessThan(supportedVersion) {
		return fmt.Errorf("found flux version %s, eksctl requires >= %s", v.String(), minSupportedVersion)
	}

	return nil
}
