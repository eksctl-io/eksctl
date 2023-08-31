package kubectl

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/blang/semver"
	"github.com/kballard/go-shellquote"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
)

const command = "kubectl"

var versionArgs = []string{"version", "--output=json"}

// gitVersion holds git version info of kubectl client/server
type gitVersion struct {
	GitVersion string `json:"gitVersion"`
}

// kubectlInfo holds version info of kubectl client & server
type kubectlInfo struct {
	ClientVersion gitVersion `json:"clientVersion"`
	ServerVersion gitVersion `json:"serverVersion"`
}

// Client implements Kubectl
type Client struct {
	GlobalArgs []string
	Env        []string
}

// NewClient return a new kubectl client
func NewClient() *Client {
	return &Client{}
}

// GetClientVersion returns the kubectl client version
func (ktl *Client) GetClientVersion() (string, error) {
	ktl.GlobalArgs = []string{"--client"}
	defer func() {
		ktl.GlobalArgs = []string{}
	}()
	clientVersion, _, err := ktl.getVersion()
	if err != nil {
		return "", err
	}
	return clientVersion, nil
}

// GetServerVersion returns the kubernetes version on server
func (ktl *Client) GetServerVersion() (string, error) {
	if len(ktl.Env) == 0 {
		return "", fmt.Errorf("client env should be set before trying to fetch server version")
	}
	_, serverVersion, err := ktl.getVersion()
	if err != nil {
		return "", err
	}
	return serverVersion, nil
}

// getVersion returns the kubectl version
func (ktl *Client) getVersion() (string, string, error) {
	cmd := exec.Command(command, versionArgs...)
	cmd.Args = append(cmd.Args, ktl.GlobalArgs...)
	cmd.Env = append(os.Environ(), ktl.Env...)

	out, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("error running `kubectl version`: %w", err)
	}

	var info kubectlInfo
	if err := json.Unmarshal(out, &info); err != nil {
		return "", "", fmt.Errorf("error parsing `kubectl version` output: %w", err)
	}

	return info.ClientVersion.GitVersion, info.ServerVersion.GitVersion, nil
}

// CheckKubectlVersion checks version of kubectl
func (ktl *Client) CheckKubectlVersion() error {
	kubectlPath, err := exec.LookPath(command)
	if err != nil {
		return fmt.Errorf("kubectl not found, v1.10.0 or newer is required")
	}
	logger.Debug("kubectl: %q", kubectlPath)

	clientVersion, ignoredErr := ktl.GetClientVersion()
	logger.Debug("kubectl client version: %s", clientVersion)
	if ignoredErr != nil {
		logger.Debug("ignored error: %s", ignoredErr)
	}

	version, err := semver.Parse(strings.TrimLeft(clientVersion, "v"))
	if err != nil {
		if ignoredErr != nil {
			return errors.Wrapf(err, "parsing kubectl version string %s (upstream error: %s) / %q", clientVersion, ignoredErr, version)
		}
		return errors.Wrapf(err, "parsing kubectl version string %s / %q", clientVersion, version)
	}

	if version.Compare(semver.Version{
		Major: 1,
		Minor: 10,
	}) < 0 {
		return fmt.Errorf("kubectl version %s was found at %q, minimum required version to use EKS is v1.10.0", clientVersion, kubectlPath)
	}

	return nil
}

func (ktl *Client) FmtCmd(cmds ...string) string {
	args := []string{command}
	args = append(args, ktl.GlobalArgs...)
	args = append(args, cmds...)
	return shellquote.Join(args...)
}
