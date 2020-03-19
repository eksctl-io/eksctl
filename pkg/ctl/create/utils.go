package create

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func checkSubnetsGivenAsFlags(params *cmdutils.CreateClusterCmdParams) bool {
	return len(*params.Subnets[api.SubnetTopologyPrivate])+len(*params.Subnets[api.SubnetTopologyPublic]) != 0
}

func checkVersion(cmd *cmdutils.Cmd, ctl *eks.ClusterProvider, meta *api.ClusterMeta) error {
	switch meta.Version {
	case "auto":
		break
	case "":
		meta.Version = "auto"
	case "default":
		meta.Version = api.DefaultVersion
		logger.Info("will use default version (%s) for new nodegroup(s)", meta.Version)
	case "latest":
		meta.Version = api.LatestVersion
		logger.Info("will use latest version (%s) for new nodegroup(s)", meta.Version)
	default:
		if !isValidVersion(meta.Version) {
			if isDeprecatedVersion(meta.Version) {
				return fmt.Errorf("invalid version, %s is no longer supported, supported values: auto, default, latest, %s\nsee also: https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html", meta.Version, strings.Join(api.SupportedVersions(), ", "))
			}
			return fmt.Errorf("invalid version %s, supported values: auto, default, latest, %s", meta.Version, strings.Join(api.SupportedVersions(), ", "))
		}
	}

	if v := ctl.ControlPlaneVersion(); v == "" {
		return fmt.Errorf("unable to get control plane version")
	} else if meta.Version == "auto" {
		meta.Version = v
		logger.Info("will use version %s for new nodegroup(s) based on control plane version", meta.Version)
	} else if meta.Version != v {
		hint := "--version=auto"
		if cmd.ClusterConfigFile != "" {
			hint = "metadata.version: auto"
		}
		logger.Warning("will use version %s for new nodegroup(s), while control plane version is %s; to automatically inherit the version use %q", meta.Version, v, hint)
	}

	return nil
}

func isValidVersion(version string) bool {
	for _, v := range api.SupportedVersions() {
		if version == v {
			return true
		}
	}
	return false
}

func isDeprecatedVersion(version string) bool {
	for _, v := range api.DeprecatedVersions() {
		if version == v {
			return true
		}
	}
	return false
}
