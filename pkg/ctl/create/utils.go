package create

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func checkSubnetsGivenAsFlags() bool {
	return len(*subnets[api.SubnetTopologyPrivate])+len(*subnets[api.SubnetTopologyPublic]) != 0
}

func checkVersion(ctl *eks.ClusterProvider, meta *api.ClusterMeta) error {
	switch meta.Version {
	case "auto":
		break
	case "":
		meta.Version = "auto"
	case "latest":
		meta.Version = api.LatestVersion
		logger.Info("will use version latest version (%s) for new nodegroup(s)", meta.Version)
	default:
		validVersion := false
		for _, v := range api.SupportedVersions() {
			if meta.Version == v {
				validVersion = true
			}
		}
		if !validVersion {
			return fmt.Errorf("invalid version %s, supported values: auto, latest, %s", meta.Version, strings.Join(api.SupportedVersions(), ", "))
		}
	}

	if v := ctl.ControlPlaneVersion(); v == "" {
		return fmt.Errorf("unable to get control plane version")
	} else if meta.Version == "auto" {
		meta.Version = v
		logger.Info("will use version %s for new nodegroup(s) based on control plane version", meta.Version)
	} else if meta.Version != v {
		hint := "--version=auto"
		if clusterConfigFile != "" {
			hint = "metadata.version: auto"
		}
		logger.Warning("will use version %s for new nodegroup(s), while control plane version is %s; to automatically inherit the version use %q", meta.Version, v, hint)
	}

	return nil
}
