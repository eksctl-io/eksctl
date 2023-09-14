package kubectl

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/blang/semver"
	"github.com/pkg/errors"
)

const Command = "kubectl"

var (
	versionArgs = []string{"version", "--output=json"}
	execCommand = exec.Command
)

type VersionType string

var (
	Client           VersionType = "client"
	Server           VersionType = "server"
	minClientVersion             = semver.Version{
		Major: 1,
		Minor: 10,
	}
	minServerVersion = semver.Version{
		Major: 1,
		Minor: 10,
	}
)

func getMinVersionForType(vType VersionType) semver.Version {
	switch vType {
	case Client:
		return minClientVersion
	case Server:
		return minServerVersion
	default:
		return semver.Version{
			Major: 1,
			Minor: 10,
		}
	}
}

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
//counterfeiter:generate -o fakes/fake_kubernetes_version_getter.go . KubernetesVersionManager
type KubernetesVersionManager interface {
	ClientVersion() (string, error)
	ServerVersion(env []string, args []string) (string, error)
	ValidateVersion(version string, vType VersionType) error
}

// VersionManager implements KubernetesVersionManager
type VersionManager struct{}

func NewVersionManager() KubernetesVersionManager {
	return &VersionManager{}
}

// ClientVersion returns the kubectl client version
func (vm *VersionManager) ClientVersion() (string, error) {
	clientVersion, _, err := getVersion([]string{}, []string{"--client"})
	if err != nil {
		return "", err
	}
	return clientVersion, nil
}

// ServerVersion returns the kubernetes version on server
func (vm *VersionManager) ServerVersion(env []string, args []string) (string, error) {
	_, serverVersion, err := getVersion(append(os.Environ(), env...), args)
	if err != nil {
		return "", err
	}
	return serverVersion, nil
}

// ValidateVersion checks that the client / server version is valid and supported
func (vm *VersionManager) ValidateVersion(version string, vType VersionType) error {
	parsedVersion, err := semver.Parse(strings.TrimLeft(version, "v"))
	if err != nil {
		return errors.Wrapf(err, "parsing kubernetes %s version string %s / %q", vType, version, parsedVersion)
	}
	minVersion := getMinVersionForType(vType)
	if parsedVersion.Compare(getMinVersionForType(vType)) < 0 {
		return fmt.Errorf("kubernetes %s version %s was found, minimum required version is v%s", vType, version, minVersion)
	}
	return nil
}

// getVersion returns the kubernetes client / server version
func getVersion(env []string, args []string) (string, string, error) {
	cmd := execCommand(Command, versionArgs...)
	cmd.Args = append(cmd.Args, args...)
	cmd.Env = env

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
