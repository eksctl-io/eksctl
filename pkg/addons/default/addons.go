package defaultaddons

import (
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

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
