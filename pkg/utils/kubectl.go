package utils

import (
	"fmt"
	"strings"

	"github.com/blang/semver"
	shellquote "github.com/kballard/go-shellquote"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
	"github.com/weaveworks/launcher/pkg/kubectl"

	"k8s.io/client-go/tools/clientcmd"
)

func fmtKubectlCmd(ktl *kubectl.LocalClient, cmds ...string) string {
	args := []string{kubectl.Command}
	args = append(args, ktl.GlobalArgs...)
	args = append(args, cmds...)
	return shellquote.Join(args...)
}

// CheckKubectlVersion checks version of kubectl
func CheckKubectlVersion(env []string) error {
	ktl := &kubectl.LocalClient{Env: env}
	kubectlPath, err := ktl.LookPath()
	if err != nil {
		return fmt.Errorf("kubectl not found, v1.10.0 or newever is required")
	}
	logger.Debug("kubectl: %q", kubectlPath)

	clientVersion, _, ignoredErr := kubectl.GetVersionInfo(ktl)
	logger.Debug("kubectl version: %s", clientVersion)
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
	if version.Major == 1 && version.Minor < 10 {
		return fmt.Errorf("kubectl version %s was found at %q, minimum required version to use EKS is v1.10.0", clientVersion, kubectlPath)
	}
	return nil
}

// CheckAllCommands check version of kubectl, and if it can be used with either
// of the authenticator commands; most importantly it validates if kubectl can
// use kubeconfig we've created for it
func CheckAllCommands(kubeconfigPath string, isContextSet bool, contextName string, env []string) error {
	if err := CheckKubectlVersion(env); err != nil {
		return err
	}

	if authenticator, found := kubeconfig.LookupAuthenticator(); !found {
		return fmt.Errorf("could not find any of the authenticator commands: %s", strings.Join(kubeconfig.AuthenticatorCommands(), ", "))
	} else {
		logger.Debug("found authenticator: %s", authenticator)
	}

	if kubeconfigPath != "" {
		ktl := &kubectl.LocalClient{
			GlobalArgs: []string{},
			Env:        env,
		}
		if kubeconfigPath != clientcmd.RecommendedHomeFile {
			ktl.GlobalArgs = append(ktl.GlobalArgs, fmt.Sprintf("--kubeconfig=%s", kubeconfigPath))
		}
		if !isContextSet {
			ktl.GlobalArgs = append(ktl.GlobalArgs, fmt.Sprintf("--context=%s", contextName))
		}

		suggestion := fmt.Sprintf("(check '%s')", fmtKubectlCmd(ktl, "version"))

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

		logger.Info("kubectl command should work with %q, try '%s'", kubeconfigPath, fmtKubectlCmd(ktl, "get", "nodes"))
	} else {
		logger.Debug("skipping kubectl integration checks, as writing kubeconfig file is disabled")
	}

	return nil
}
