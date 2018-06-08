package utils

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/blang/semver"
	"github.com/pkg/errors"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/weaveworks/launcher/pkg/kubectl"
)

var kubectlPath string

func CheckKubectlVersion() error {
	ktl := &kubectl.LocalClient{}
	kubectlPath, err := ktl.LookPath()
	if err != nil {
		return fmt.Errorf("kubectl not found, v1.10.0 or newever is required")
	}
	logger.Debug("kubectl: %q", kubectlPath)

	clientVersion, _, err := kubectl.GetVersionInfo(ktl)
	logger.Debug("clientVersion=%#v err=%q", clientVersion, err)

	version, err := semver.Parse(strings.TrimLeft(clientVersion, "v"))
	if err != nil {
		return errors.Wrapf(err, "parsing kubectl version string %q", version)
	}
	if version.Major == 1 && version.Minor < 10 {
		return fmt.Errorf("kubectl version %s was found at %q, minimum required version to use EKS is v1.10.0", clientVersion, kubectlPath)
	}
	return nil
}

func CheckHeptioAuthenticatorAWS() error {
	path, err := exec.LookPath("heptio-authenticator-aws")
	if err == nil {
		logger.Debug("heptio-authenticator-aws: %q", path)
	} else {
		return fmt.Errorf("heptio-authenticator-aws not installed")
	}
	return nil
}

func CheckAllCommands(kubeconfigPath string) error {
	if err := CheckKubectlVersion(); err != nil {
		return err
	}

	if err := CheckHeptioAuthenticatorAWS(); err != nil {
		return err
	}

	if kubeconfigPath != "" {
		ktl := &kubectl.LocalClient{
			GlobalArgs: []string{"--kubeconfig", kubeconfigPath},
		}

		suggestion := fmt.Sprintf("(check '%s %s version')", kubectlPath, strings.Join(ktl.GlobalArgs, " "))

		_, serverVersion, err := kubectl.GetVersionInfo(ktl)
		if err != nil {
			return errors.Wrapf(err, "unable to use kubectl with the EKS cluster %s", suggestion)
		}
		version, err := semver.Parse(strings.TrimLeft(serverVersion, "v"))
		if err != nil {
			return errors.Wrapf(err, "parsing Kubernetes version string %q return by the EKS API server", version)
		}
		if version.Major == 1 && version.Minor < 10 {
			return fmt.Errorf("Kubernetes version %s found, v1.10.0 or newer is expected with EKS %s", serverVersion, suggestion)
		}

		logger.Info("all command should work, try '%s %s get nodes'", kubectlPath, strings.Join(ktl.GlobalArgs, " "))

	} else {
		logger.Debug("skipping kubectl integration ckecks, as writing kubeconfig file is disabled")
	}

	return nil
}
