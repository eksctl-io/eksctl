package defaultaddons

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/hashicorp/go-version"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	v1 "k8s.io/api/apps/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/weaveworks/eksctl/pkg/addons"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/printers"
	"github.com/weaveworks/eksctl/pkg/utils"
)

const (
	// KubeProxy is the name of the kube-proxy addon
	KubeProxy     = "kube-proxy"
	ArchBetaLabel = "beta.kubernetes.io/arch"
	ArchLabel     = "kubernetes.io/arch"
)

func IsKubeProxyUpToDate(input AddonInput) (bool, error) {
	d, err := input.RawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), KubeProxy, metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("%q was not found", KubeProxy)
			return true, nil
		}
		return false, errors.Wrapf(err, "getting %q", KubeProxy)
	}
	if numContainers := len(d.Spec.Template.Spec.Containers); !(numContainers >= 1) {
		return false, fmt.Errorf("%s has %d containers, expected at least 1", KubeProxy, numContainers)
	}

	greaterThanOrEqualTo1_18, err := utils.IsMinVersion(api.Version1_18, input.ControlPlaneVersion)
	if err != nil {
		return false, err
	}

	desiredTag, err := getLatestKubeProxyImage(input, greaterThanOrEqualTo1_18)
	if err != nil {
		return false, err
	}
	image := d.Spec.Template.Spec.Containers[0].Image
	imageTag, err := addons.ImageTag(image)
	if err != nil {
		return false, err
	}
	return desiredTag == imageTag, nil
}

// UpdateKubeProxy updates image tag for kube-system:daemonset/kube-proxy based to match ControlPlaneVersion
func UpdateKubeProxy(input AddonInput, plan bool) (bool, error) {
	printer := printers.NewJSONPrinter()

	d, err := input.RawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), KubeProxy, metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("%q was not found", KubeProxy)
			return false, nil
		}
		return false, errors.Wrapf(err, "getting %q", KubeProxy)
	}

	archLabel := ArchLabel
	greaterThanOrEqualTo1_18, err := utils.IsMinVersion(api.Version1_18, input.ControlPlaneVersion)
	if err != nil {
		return false, err
	}
	if !greaterThanOrEqualTo1_18 {
		archLabel = ArchBetaLabel
	}

	hasArm64NodeSelector := daemeonSetHasArm64NodeSelector(d, archLabel)
	if !hasArm64NodeSelector {
		logger.Info("missing arm64 nodeSelector value")
	}

	if numContainers := len(d.Spec.Template.Spec.Containers); !(numContainers >= 1) {
		return false, fmt.Errorf("%s has %d containers, expected at least 1", KubeProxy, numContainers)
	}

	if err := printer.LogObj(logger.Debug, KubeProxy+" [current] = \\\n%s\n", d); err != nil {
		return false, err
	}

	image := &d.Spec.Template.Spec.Containers[0].Image
	imageParts := strings.Split(*image, ":")

	if len(imageParts) != 2 {
		return false, fmt.Errorf("unexpected image format %q for %q", *image, KubeProxy)
	}

	desiredTag, err := getLatestKubeProxyImage(input, greaterThanOrEqualTo1_18)
	if err != nil {
		return false, err
	}
	if imageParts[1] == desiredTag && hasArm64NodeSelector {
		logger.Debug("imageParts = %v, desiredTag = %s", imageParts, desiredTag)
		logger.Info("%q is already up-to-date", KubeProxy)
		return false, nil
	}

	if plan {
		logger.Critical("(plan) %q is not up-to-date", KubeProxy)
		return true, nil
	}

	imageParts[1] = desiredTag
	*image = strings.Join(imageParts, ":")

	if err := printer.LogObj(logger.Debug, KubeProxy+" [updated] = \\\n%s\n", d); err != nil {
		return false, err
	}

	if !hasArm64NodeSelector {
		if err := addArm64NodeSelector(d, archLabel); err != nil {
			return false, err
		}
	}

	if _, err := input.RawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Update(context.TODO(), d, metav1.UpdateOptions{}); err != nil {
		return false, err
	}

	logger.Info("%q is now up-to-date", KubeProxy)
	return false, nil
}

func daemeonSetHasArm64NodeSelector(daemonSet *v1.DaemonSet, archLabel string) bool {
	if daemonSet.Spec.Template.Spec.Affinity != nil &&
		daemonSet.Spec.Template.Spec.Affinity.NodeAffinity != nil &&
		daemonSet.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil {
		for _, nodeSelectorTerms := range daemonSet.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms {
			for _, nodeSelector := range nodeSelectorTerms.MatchExpressions {
				if nodeSelector.Key == archLabel {
					for _, value := range nodeSelector.Values {
						if value == "arm64" {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func addArm64NodeSelector(daemonSet *v1.DaemonSet, archLabel string) error {
	if daemonSet.Spec.Template.Spec.Affinity != nil && daemonSet.Spec.Template.Spec.Affinity.NodeAffinity != nil {
		for nodeSelectorTermsIndex, nodeSelectorTerms := range daemonSet.Spec.Template.Spec.Affinity.NodeAffinity.
			RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms {
			for nodeSelectorIndex, nodeSelector := range nodeSelectorTerms.MatchExpressions {
				if nodeSelector.Key == archLabel {
					daemonSet.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.
						NodeSelectorTerms[nodeSelectorTermsIndex].MatchExpressions[nodeSelectorIndex].Values = append(nodeSelector.Values, "arm64")
				}
			}
		}
		return nil
	}
	return fmt.Errorf("NodeAffinity not configured on kube-proxy. Either manually update the proxy deployment, or switch to Managed Addons")
}

func getLatestKubeProxyImage(input AddonInput, greaterThanOrEqualTo1_18 bool) (string, error) {
	defaultClusterVersion := generateImageVersionFromClusterVersion(input.ControlPlaneVersion)
	// EKS Addons API only works for 1.18 and above
	if !greaterThanOrEqualTo1_18 {
		return defaultClusterVersion, nil
	}

	latestEKSReportedVersion, err := getLatestImageVersionFromEKS(input.EKSAPI, input.ControlPlaneVersion)
	if err != nil {
		return "", err
	}

	// Sometimes the EKS API is ahead, sometimes behind. Pick whichever is latest
	eksVersionIsGreaterThanDefaultVersion, err := versionGreaterThan(latestEKSReportedVersion, defaultClusterVersion)
	if err != nil {
		return "", err
	}

	if eksVersionIsGreaterThanDefaultVersion {
		return latestEKSReportedVersion, nil
	}

	return defaultClusterVersion, nil
}

func versionGreaterThan(v1, v2 string) (bool, error) {
	v1Version, err := parseVersion(v1)
	if err != nil {
		return false, err
	}
	v2Version, err := parseVersion(v2)
	if err != nil {
		return false, err
	}
	return v1Version.GreaterThan(v2Version), nil
}

func generateImageVersionFromClusterVersion(controlPlaneVersion string) string {
	return fmt.Sprintf("v%s-eksbuild.1", controlPlaneVersion)
}

func getLatestImageVersionFromEKS(eksAPI eksiface.EKSAPI, controlPlaneVersion string) (string, error) {
	controlPlaneMajorMinor, err := versionWithOnlyMajorAndMinor(controlPlaneVersion)
	if err != nil {
		return "", err
	}
	input := &eks.DescribeAddonVersionsInput{
		KubernetesVersion: &controlPlaneMajorMinor,
		AddonName:         aws.String(KubeProxy),
	}

	addonInfos, err := eksAPI.DescribeAddonVersions(input)
	if err != nil {
		return "", fmt.Errorf("failed to describe addon versions: %v", err)
	}

	if len(addonInfos.Addons) == 0 || len(addonInfos.Addons[0].AddonVersions) == 0 {
		return "", fmt.Errorf("no versions available for %q", KubeProxy)
	}

	var versions []*version.Version
	for _, addonVersionInfo := range addonInfos.Addons[0].AddonVersions {
		v, err := parseVersion(*addonVersionInfo.AddonVersion)
		if err != nil {
			return "", err
		}

		versions = append(versions, v)
	}

	sort.SliceStable(versions, func(i, j int) bool {
		return versions[j].LessThan(versions[i])
	})
	return versions[0].Original(), nil
}

func versionWithOnlyMajorAndMinor(v string) (string, error) {
	parsedVersion, err := parseVersion(v)
	if err != nil {
		return "", err
	}
	parsedVersionSegments := parsedVersion.Segments()
	return fmt.Sprintf("%d.%d", parsedVersionSegments[0], parsedVersionSegments[1]), nil
}

func parseVersion(v string) (*version.Version, error) {
	version, err := version.NewVersion(v)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version %q: %w", v, err)
	}
	return version, nil
}
