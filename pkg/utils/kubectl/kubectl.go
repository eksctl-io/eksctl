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

var (
	versionArgs  = []string{"version", "--output=json"}
	execCommand  = exec.Command
	execLookPath = exec.LookPath
)

// gitVersion holds git version info of kubectl client/server
type gitVersion struct {
	GitVersion string `json:"gitVersion"`
}

// kubectlInfo holds version info of kubectl client & server
type kubectlInfo struct {
	ClientVersion gitVersion `json:"clientVersion"`
	ServerVersion gitVersion `json:"serverVersion"`
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_kubectl_client.go . KubernetesClient
type KubernetesClient interface {
	GetClientVersion() (string, error)
	GetServerVersion() (string, error)
	CheckKubectlVersion() error
	FmtCmd(args []string) string
	SetEnv(env []string)
	AppendArgForNextCmd(arg string)
}

// Client implements Kubectl
type Client struct {
	args []string
	env  []string
}

// NewClient return a new kubectl client
func NewClient() KubernetesClient {
	return &Client{}
}

// SetEnv sets env
func (ktl *Client) SetEnv(env []string) {
	ktl.env = env
}

// AppendArgForNextCmd adds args for the next command to be run
func (ktl *Client) AppendArgForNextCmd(arg string) {
	ktl.args = append(ktl.args, arg)
}

// cleanupArgs removes all args that the current command used
func (ktl *Client) cleanupArgs() {
	ktl.args = []string{}
}

// GetClientVersion returns the kubectl client version
func (ktl *Client) GetClientVersion() (string, error) {
	ktl.AppendArgForNextCmd("--client")
	defer ktl.cleanupArgs()
	clientVersion, _, err := ktl.getVersion()
	if err != nil {
		return "", err
	}
	return clientVersion, nil
}

// GetServerVersion returns the kubernetes version on server
func (ktl *Client) GetServerVersion() (string, error) {
	if len(ktl.env) == 0 {
		return "", fmt.Errorf("client env should be set before trying to fetch server version")
	}
	defer ktl.cleanupArgs()
	_, serverVersion, err := ktl.getVersion()
	if err != nil {
		return "", err
	}
	return serverVersion, nil
}

// getVersion returns the kubectl version
func (ktl *Client) getVersion() (string, string, error) {
	cmd := execCommand(command, versionArgs...)
	cmd.Args = append(cmd.Args, ktl.args...)
	cmd.Env = append(os.Environ(), ktl.env...)

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
	kubectlPath, err := execLookPath(command)
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

func (ktl *Client) FmtCmd(args []string) string {
	cmd := []string{command}
	cmd = append(cmd, args...)
	return shellquote.Join(cmd...)
}
