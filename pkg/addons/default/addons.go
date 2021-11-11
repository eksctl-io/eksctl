package defaultaddons

import (
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

type AddonInput struct {
	RawClient           kubernetes.RawClientInterface
	EKSAPI              eksiface.EKSAPI
	ControlPlaneVersion string
	Region              string
}

func DoAddonsSupportMultiArch(eksAPI eksiface.EKSAPI, rawClient kubernetes.RawClientInterface, controlPlaneVersion string, region string) (bool, error) {
	input := AddonInput{
		RawClient:           rawClient,
		ControlPlaneVersion: controlPlaneVersion,
		Region:              region,
		EKSAPI:              eksAPI,
	}
	kubeProxyUpToDate, err := IsKubeProxyUpToDate(input)
	if err != nil {
		return true, err
	}
	if !kubeProxyUpToDate {
		return false, nil
	}

	awsNodeUpToDate, err := DoesAWSNodeSupportMultiArch(input)
	if err != nil {
		return true, err
	}
	if !awsNodeUpToDate {
		return false, nil
	}

	coreDNSUpToDate, err := IsCoreDNSUpToDate(input)
	if err != nil {
		return true, err
	}
	return coreDNSUpToDate, nil
}
