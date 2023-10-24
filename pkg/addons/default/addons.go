package defaultaddons

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

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
// We know that AWS node requires 1.6.3+ to work, so we check for that.
// For kube-proxy and CoreDNS, we do not know what version adds support, so we just ensure they contain a node affinity
// that allows them to be scheduled on ARM64 nodes.
func DoAddonsSupportMultiArch(ctx context.Context, clientSet kubernetes.Interface) (bool, error) {
	kubeProxy, err := getKubeProxy(ctx, clientSet)
	if err != nil {
		return false, err
	}

	if kubeProxy != nil && !supportsMultiArch(kubeProxy.Spec.Template.Spec) {
		return false, nil
	}

	awsNodeSupportsMultiArch, err := DoesAWSNodeSupportMultiArch(ctx, clientSet)
	if err != nil {
		return false, err
	}
	if !awsNodeSupportsMultiArch {
		return false, nil
	}

	coreDNS, err := getCoreDNS(ctx, clientSet)
	if err != nil {
		return false, err
	}
	return coreDNS == nil || supportsMultiArch(coreDNS.Spec.Template.Spec), nil
}

// supportsMultiArch returns true if the PodSpec contains a node affinity that allows the pod to be scheduled on
// multiple architectures.
func supportsMultiArch(podSec corev1.PodSpec) bool {
	if podSec.Affinity == nil || podSec.Affinity.NodeAffinity == nil || podSec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution == nil {
		return false
	}
	for _, nodeSelectorTerm := range podSec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms {
		for _, me := range nodeSelectorTerm.MatchExpressions {
			if me.Key == corev1.LabelArchStable && me.Operator == corev1.NodeSelectorOpIn {
				for _, val := range me.Values {
					if val == "arm64" {
						return true
					}
				}
			}
		}
	}
	return false
}

func makeGetError[T any](resource *T, err error, resourceName string) (*T, error) {
	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.Warning("%q was not found", resourceName)
			return nil, nil
		}
		return nil, fmt.Errorf("getting %q: %w", resourceName, err)
	}
	return resource, nil
}

// LoadAsset return embedded manifest as a runtime.Object
func newList(data []byte) (*metav1.List, error) {
	list, err := kubernetes.NewList(data)
	if err != nil {
		return nil, errors.Wrapf(err, "loading individual resources from manifest")
	}
	return list, nil
}
