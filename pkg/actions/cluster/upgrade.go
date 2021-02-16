package cluster

import (
	"fmt"

	"github.com/weaveworks/eksctl/pkg/printers"

	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/logger"
)

func upgrade(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, dryRun bool) (bool, error) {
	currentVersion := ctl.ControlPlaneVersion()
	versionUpdateRequired, err := requiresVersionUpgrade(cfg.Metadata, currentVersion)
	if err != nil {
		return false, err
	}

	printer := printers.NewJSONPrinter()
	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return false, err
	}

	if versionUpdateRequired {
		msgNodeGroupsAndAddons := "you will need to follow the upgrade procedure for all of nodegroups and add-ons"
		cmdutils.LogIntendedAction(dryRun, "upgrade cluster %q control plane from current version %q to %q", cfg.Metadata.Name, currentVersion, cfg.Metadata.Version)
		if !dryRun {
			if err := ctl.UpdateClusterVersionBlocking(cfg); err != nil {
				return false, err
			}
			logger.Success("cluster %q control plane has been upgraded to version %q", cfg.Metadata.Name, cfg.Metadata.Version)
			logger.Info(msgNodeGroupsAndAddons)
		}
	} else {
		logger.Info("no cluster version update required")
	}
	return versionUpdateRequired, nil
}

func requiresVersionUpgrade(clusterMeta *api.ClusterMeta, currentEKSVersion string) (bool, error) {
	nextVersion, err := getNextVersion(currentEKSVersion)
	if err != nil {
		return false, err
	}

	// If the version was not specified default to the next Kubernetes version and assume the user intended to upgrade if possible
	// also support "auto" as version (see #2461)
	if clusterMeta.Version == "" || clusterMeta.Version == "auto" {
		if api.IsSupportedVersion(nextVersion) {
			clusterMeta.Version = nextVersion
			return true, nil
		}

		// There is no new version, stay in the current one
		clusterMeta.Version = currentEKSVersion
		return false, nil
	}

	if c, err := utils.CompareVersions(clusterMeta.Version, currentEKSVersion); err != nil {
		return false, errors.Wrap(err, "couldn't compare versions for upgrade")
	} else if c < 0 {
		return false, fmt.Errorf("cannot upgrade to a lower version. Found given target version %q, current cluster version %q", clusterMeta.Version, currentEKSVersion)
	}

	if api.IsDeprecatedVersion(clusterMeta.Version) {
		return false, fmt.Errorf("control plane version %q has been deprecated", clusterMeta.Version)
	}

	if !api.IsSupportedVersion(clusterMeta.Version) {
		return false, fmt.Errorf("control plane version %q is not known to this version of eksctl, try to upgrade eksctl first", clusterMeta.Version)
	}

	if clusterMeta.Version == currentEKSVersion {
		return false, nil
	}

	if clusterMeta.Version == nextVersion {
		return true, nil
	}

	return false, fmt.Errorf(
		"upgrading more than one version at a time is not supported. Found upgrade from %q to %q. Please upgrade to %q first",
		currentEKSVersion,
		clusterMeta.Version,
		nextVersion)
}

func getNextVersion(currentVersion string) (string, error) {
	switch currentVersion {
	case "":
		return "", errors.New("unable to get control plane version")
	case api.Version1_12:
		return api.Version1_13, nil
	case api.Version1_13:
		return api.Version1_14, nil
	case api.Version1_14:
		return api.Version1_15, nil
	case api.Version1_15:
		return api.Version1_16, nil
	case api.Version1_16:
		return api.Version1_17, nil
	case api.Version1_17:
		return api.Version1_18, nil
	case api.Version1_18:
		return api.Version1_19, nil
	case api.Version1_19:
		return api.Version1_20, nil
	default:
		// version of control plane is not known to us, maybe we are just too old...
		return "", fmt.Errorf("control plane version %q is not known to this version of eksctl, try to upgrade eksctl first", currentVersion)
	}
}
