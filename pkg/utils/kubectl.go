package utils

import (
	"fmt"
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
	logger.Debug("clientVersion=%#v err=%v", clientVersion, err)

	version, err := semver.Parse(strings.TrimLeft(clientVersion, "v"))
	if err != nil {
		return errors.Wrapf(err, "parsing kubectl version string %q", clientVersion)
	}
	if version.Major == 1 && version.Minor < 10 {
		return fmt.Errorf("kubectl version %s was found at %q, minimum required version to use EKS is v1.10.0", clientVersion, ktl.CommandPath)
	}
	return nil
}
