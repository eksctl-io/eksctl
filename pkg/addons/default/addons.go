package defaultaddons

import (
	"context"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

type AddonInput struct {
	RawClient           kubernetes.RawClientInterface
	EKSAPI              awsapi.EKS
	ControlPlaneVersion string
	Region              string
}

// DoAddonsSupportMultiArch checks if the coredns/kubeproxy/awsnode support multi arch nodegroups
// We know that AWS node requires 1.6.3+ to work, so we check for that
// Kubeproxy/coredns we don't know what version adds support, so we just ensure its up-to-date before proceeding.
// TODO: we should know what versions of kubeproxy/coredns added support, rather than always erroring if they are out of date
func DoAddonsSupportMultiArch(ctx context.Context, eksAPI awsapi.EKS, rawClient kubernetes.RawClientInterface, controlPlaneVersion string, region string) (bool, error) {
	input := AddonInput{
		RawClient:           rawClient,
		ControlPlaneVersion: controlPlaneVersion,
		Region:              region,
		EKSAPI:              eksAPI,
	}
	kubeProxyUpToDate, err := IsKubeProxyUpToDate(ctx, input)
	if err != nil {
		return true, err
	}
	if !kubeProxyUpToDate {
		return false, nil
	}

	awsNodeUpToDate, err := DoesAWSNodeSupportMultiArch(ctx, input)
	if err != nil {
		return true, err
	}
	if !awsNodeUpToDate {
		return false, nil
	}

	coreDNSUpToDate, err := IsCoreDNSUpToDate(ctx, input)
	if err != nil {
		return true, err
	}
	return coreDNSUpToDate, nil
}

// LoadAsset return embedded manifest as a runtime.Object
func newList(data []byte) (*metav1.List, error) {
	list, err := kubernetes.NewList(data)
	if err != nil {
		return nil, errors.Wrapf(err, "loading individual resources from manifest")
	}
	return list, nil
}
