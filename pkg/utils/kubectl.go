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

func CheckKubectlVersion() error {
	ktl := &kubectl.LocalClient{}
	if ktl.IsPresent() {
		logger.Debug("kubectl: %q", ktl.CommandPath)
	} else {
		return fmt.Errorf("kubectl not installed, you should install kubectl v1.10.0 or newever")
	}
	clientVersion, _, err := kubectl.GetVersionInfo(ktl)
	logger.Debug("clientVersion=%#v err=%q", clientVersion, err.Error())

	version, err := semver.Parse(strings.TrimLeft(clientVersion, "v"))
	if err != nil {
		return errors.Wrapf(err, "parsing kubectl version string %q", version)
	}
	if version.Major == 1 && version.Minor < 10 {
		return fmt.Errorf("kubectl version %s was found at %q, minimum required version to use EKS is v1.10.0", clientVersion, ktl.CommandPath)
	}
	return nil
}

func CheckHeptioAuthenticatorAWS() error {
	path, err := exec.LookPath("heptio-authenticator-aws")
	if err == nil {
		logger.Debug("heptio-authenticator-aws: %q", path)
	} else {
		return fmt.Errorf("heptio-authenticator-aws not installed, you should install it")
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
		if !ktl.IsPresent() {
			return fmt.Errorf("re-checking kubectl command failed – kubectl binary gone missing unexpectedly")
		}
		_, serverVersion, err := kubectl.GetVersionInfo(ktl)
		if err != nil {
			suggestion := fmt.Sprintf("%s %s version", ktl.CommandPath, strings.Join(ktl.GlobalArgs, " "))
			return errors.Wrapf(err, "unable to use kubectl with the EKS cluster (check '%s')", suggestion)
		}
		version, err := semver.Parse(strings.TrimLeft(serverVersion, "v"))
		if err != nil {
			return errors.Wrapf(err, "parsing Kubernetes version string %q return by the EKS API server", version)
		}
		if version.Major == 1 && version.Minor < 10 {
			return fmt.Errorf("Kubernetes version %s is unexpected with EKS, it should be v1.10.0 or newever", serverVersion)
		}
		// TODO(p2): we can do a littl bit more here
	} else {
		logger.Debug("skipping kubectl integration ckecks, as wriging kubeconfig file is disabled")
	}

	return nil
}
