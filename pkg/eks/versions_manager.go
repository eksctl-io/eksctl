package eks

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/utils"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_cluster_versions_manager.go . ClusterVersionsManagerInterface
type ClusterVersionsManagerInterface interface {
	DefaultVersion() string
	LatestVersion() string
	SupportedVersions() []string
	IsSupportedVersion(string) bool
	IsDeprecatedVersion(string) bool
	ValidateVersion(string) error
	ResolveClusterVersion(string) (string, error)
	ResolveUpgradeVersion(string, string) (string, error)
}

type ClusterVersionsManager struct {
	eksAPI awsapi.EKS
	versionsInfo
}

type versionsInfo struct {
	supportedVersions                            []string
	deprecatedVersions                           []string
	defaultVersion, latestVersion, oldestVersion string
}

func NewClusterVersionsManager(eksAPI awsapi.EKS) (ClusterVersionsManagerInterface, error) {
	cvm := &ClusterVersionsManager{
		eksAPI: eksAPI,
	}
	output, err := cvm.eksAPI.DescribeClusterVersions(context.TODO(), &awseks.DescribeClusterVersionsInput{})
	if err != nil {
		return nil, fmt.Errorf("describing cluster versions: %w", err)
	}

	// build a list of supported versions,
	for _, version := range output.ClusterVersions {
		if version.Status != "UNSUPPORTED" {
			cvm.supportedVersions = append(cvm.supportedVersions, *version.ClusterVersion)
		}
		if version.DefaultVersion {
			cvm.defaultVersion = *version.ClusterVersion
		}
	}

	// resolve oldest and latest versions
	sortVersions(cvm.supportedVersions)
	cvm.oldestVersion = cvm.supportedVersions[0]
	cvm.latestVersion = cvm.supportedVersions[len(cvm.supportedVersions)-1]

	// build a list of deprecated versions
	cvm.deprecatedVersions, err = resolveDeprecatedVersions(cvm.oldestVersion)
	if err != nil {
		return nil, fmt.Errorf("resolving deprecated EKS versions: %w", err)
	}

	return cvm, nil
}

func (cvm *ClusterVersionsManager) SupportedVersions() []string {
	return cvm.supportedVersions
}

func (cvm *ClusterVersionsManager) IsSupportedVersion(version string) bool {
	for _, v := range cvm.SupportedVersions() {
		if version == v {
			return true
		}
	}
	return false
}

func (cvm *ClusterVersionsManager) DeprecatedVersions() []string {
	return cvm.deprecatedVersions
}

func (cvm *ClusterVersionsManager) IsDeprecatedVersion(version string) bool {
	for _, v := range cvm.DeprecatedVersions() {
		if version == v {
			return true
		}
	}
	return false
}

func (cvm *ClusterVersionsManager) DefaultVersion() string {
	//TODO: replace with output from DescribeClusterVersions endpoint
	return api.DefaultVersion
}

func (cvm *ClusterVersionsManager) LatestVersion() string {
	return cvm.latestVersion
}

func (cvm *ClusterVersionsManager) ValidateVersion(version string) error {
	if !cvm.IsSupportedVersion(version) {
		if cvm.IsDeprecatedVersion(version) {
			return fmt.Errorf("invalid version, %s is no longer supported, supported values: %s\nsee also: https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html", version, strings.Join(cvm.supportedVersions, ", "))
		}
		return fmt.Errorf("invalid version, supported values: %s", strings.Join(cvm.supportedVersions, ", "))
	}
	return nil
}

func (cvm *ClusterVersionsManager) ResolveClusterVersion(version string) (string, error) {
	switch version {
	case "auto", "":
		return cvm.DefaultVersion(), nil
	case "latest":
		return cvm.LatestVersion(), nil
	default:
		if err := cvm.ValidateVersion(version); err != nil {
			return "", err
		}
		return version, nil
	}
}

func (cvm *ClusterVersionsManager) ResolveUpgradeVersion(desiredVersion string, currentVersion string) (string, error) {
	// Resolve next version
	var nextVersion string
	switch {
	case currentVersion == "":
		return "", fmt.Errorf("couldn't resolve control plane version")
	case cvm.IsDeprecatedVersion(currentVersion):
		return "", fmt.Errorf("control plane version %q has been deprecated", currentVersion)
	case !cvm.IsSupportedVersion(currentVersion):
		return "", fmt.Errorf("control plane version %q is not supported", currentVersion)
	default:
		i := slices.Index(cvm.supportedVersions, currentVersion)
		if i == len(cvm.supportedVersions)-1 {
			logger.Info("control plane is already on latest version %q", currentVersion)
			return "", nil
		}
		nextVersion = cvm.supportedVersions[i+1]
	}

	// If the version was not specified, default to the next Kubernetes version, and assume the user intended to upgrade if possible.
	// Also support "auto" as version (see #2461)
	if desiredVersion == "" || desiredVersion == "auto" {
		if cvm.IsSupportedVersion(nextVersion) {
			return nextVersion, nil
		}
		// There is no new version, stay in the current one
		return "", nil
	}

	if c, err := utils.CompareVersions(desiredVersion, currentVersion); err != nil {
		return "", fmt.Errorf("couldn't compare versions for upgrade: %w", err)
	} else if c < 0 {
		return "", fmt.Errorf("cannot upgrade to a lower version. Found given target version %q, current cluster version %q", desiredVersion, currentVersion)
	}

	if cvm.IsDeprecatedVersion(desiredVersion) {
		return "", fmt.Errorf("control plane version %q has been deprecated", desiredVersion)
	}

	if !cvm.IsSupportedVersion(desiredVersion) {
		return "", fmt.Errorf("control plane version %q is not supported", desiredVersion)
	}

	if desiredVersion == currentVersion {
		return "", nil
	}

	if desiredVersion == nextVersion {
		return desiredVersion, nil
	}

	return "", fmt.Errorf(
		"upgrading more than one version at a time is not supported. Found upgrade from %q to %q. Please upgrade to %q first",
		currentVersion,
		desiredVersion,
		nextVersion)
}

func resolveDeprecatedVersions(currentVersion string) ([]string, error) {
	parts := strings.Split(currentVersion, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid version format: %s", currentVersion)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version in: %s", currentVersion)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version in: %s", currentVersion)
	}

	var versions []string
	for m := major; m >= 1; m-- {
		startMinor := minor
		if m < major {
			startMinor = 99
		}

		for n := startMinor - 1; n >= 0; n-- {
			versions = append(versions, fmt.Sprintf("%d.%d", m, n))
		}
	}
	return versions, nil
}

func sortVersions(versions []string) []string {
	sort.Slice(versions, func(i, j int) bool {
		v1Parts := strings.Split(versions[i], ".")
		v2Parts := strings.Split(versions[j], ".")

		v1Major, _ := strconv.Atoi(v1Parts[0])
		v1Minor, _ := strconv.Atoi(v1Parts[1])

		v2Major, _ := strconv.Atoi(v2Parts[0])
		v2Minor, _ := strconv.Atoi(v2Parts[1])

		// Compare major versions first, then minor versions
		if v1Major != v2Major {
			return v1Major < v2Major
		}
		return v1Minor < v2Minor
	})
	return versions
}
