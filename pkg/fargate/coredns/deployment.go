package coredns

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	kubeclient "k8s.io/client-go/kubernetes"
)

type betaDeployment struct {
	*v1beta1.Deployment
}

func (d betaDeployment) PodAnnotations() map[string]string { return d.Spec.Template.Annotations }

func (d betaDeployment) Replicas() *int32 { return d.Spec.Replicas }

func (d betaDeployment) ReadyReplicas() int32 { return d.Status.ReadyReplicas }

type appsDeployment struct {
	*appsv1.Deployment
}

func (d appsDeployment) PodAnnotations() map[string]string { return d.Spec.Template.Annotations }

func (d appsDeployment) Replicas() *int32 { return d.Spec.Replicas }

func (d appsDeployment) ReadyReplicas() int32 { return d.Status.ReadyReplicas }

type deployment interface {
	runtime.Object
	PodAnnotations() map[string]string
	Replicas() *int32
	ReadyReplicas() int32
}

func getDeployment(clientSet kubeclient.Interface, useBetaAPIGroup bool) (deployment, error) {
	if useBetaAPIGroup {
		deployment, err := clientSet.ExtensionsV1beta1().Deployments(Namespace).Get(Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return betaDeployment{deployment}, nil
	}

	deployment, err := clientSet.AppsV1().Deployments(Namespace).Get(Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return appsDeployment{deployment}, nil
}

func patchDeployment(clientSet kubeclient.Interface, useBetaAPIGroup bool, bytes []byte) (deployment, error) {
	if useBetaAPIGroup {
		patchedDeployment, err := clientSet.ExtensionsV1beta1().Deployments(Namespace).Patch(Name, types.MergePatchType, bytes)
		if err != nil {
			return nil, err
		}
		return betaDeployment{patchedDeployment}, nil
	}

	patchedDeployment, err := clientSet.AppsV1().Deployments(Namespace).Patch(Name, types.MergePatchType, bytes)
	if err != nil {
		return nil, err
	}
	return appsDeployment{patchedDeployment}, nil
}
