package defaultaddons

import (
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

// EnsureAllAddonsUpToDate checks the default addons (aws-node, kube-proxy and coredns) are up to date
// and if they are not it will update them to the latest version for the given control plane version
// TODO Delete this function when the multi architecture images are used by default when a new cluster is created. When
// that happens eksctl won't need to update them before creating ARM nodegroups anymore.
func EnsureAddonsUpToDate(clientSet kubernetes.Interface, rawClient kubernetes.RawClientInterface, controlPlaneVersion string, region string) error {
	_, err := UpdateKubeProxyImageTag(clientSet, controlPlaneVersion, false)
	if err != nil {
		return errors.Wrapf(err, "error updating kube-proxy")
	}

	_, err = UpdateAWSNode(rawClient, region, false)
	if err != nil {
		return errors.Wrapf(err, "error updating aws-node")
	}

	_, err = UpdateCoreDNS(rawClient, region, controlPlaneVersion, false)
	if err != nil {
		return errors.Wrapf(err, "error updating coredns")
	}

	return nil
}

func DoAddonsSupportMultiArch(clientSet kubernetes.Interface, rawClient kubernetes.RawClientInterface, controlPlaneVersion string, region string) (bool, error) {
	kubeProxyUpToDate, err := IsKubeProxyUpToDate(clientSet, controlPlaneVersion)
	if err != nil {
		return true, err
	}
	if !kubeProxyUpToDate {
		return false, nil
	}

	awsNodeUpToDate, err := DoesAWSNodeSupportMultiArch(rawClient, region)
	if err != nil {
		return true, err
	}
	if !awsNodeUpToDate {
		return false, nil
	}

	coreDNSUpToDate, err := IsCoreDNSUpToDate(rawClient, region, controlPlaneVersion)
	if err != nil {
		return true, err
	}
	return coreDNSUpToDate, nil
}
