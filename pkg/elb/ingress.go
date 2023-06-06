package elb

import (
	"context"

	"github.com/kris-nova/logger"

	v1 "k8s.io/api/networking/v1"
	"k8s.io/api/networking/v1beta1"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const ingressClassAnnotation = "kubernetes.io/ingress.class"

type Ingress interface {
	Delete(kubernetesCS kubernetes.Interface) error
	GetIngressClass() string
	GetMetadata() metav1.ObjectMeta
	GetLoadBalancersHosts() []string
}

type v1BetaIngress struct {
	ingress v1beta1.Ingress
}

func (i *v1BetaIngress) Delete(kubernetesCS kubernetes.Interface) error {
	return kubernetesCS.NetworkingV1beta1().Ingresses(i.ingress.Namespace).Delete(context.TODO(), i.ingress.Name, metav1.DeleteOptions{})
}

func (i *v1BetaIngress) GetIngressClass() string {
	if i.ingress.Spec.IngressClassName != nil {
		return *i.ingress.Spec.IngressClassName
	}
	return i.ingress.ObjectMeta.Annotations[ingressClassAnnotation]
}

func (i *v1BetaIngress) GetMetadata() metav1.ObjectMeta {
	return i.ingress.ObjectMeta
}

func (i *v1BetaIngress) GetLoadBalancersHosts() []string {
	var hostNames []string
	for _, hostName := range i.ingress.Status.LoadBalancer.Ingress {
		hostNames = append(hostNames, hostName.Hostname)
	}

	return hostNames
}

type v1Ingress struct {
	ingress v1.Ingress
}

func (i *v1Ingress) Delete(kubernetesCS kubernetes.Interface) error {
	return kubernetesCS.NetworkingV1().Ingresses(i.ingress.Namespace).Delete(context.TODO(), i.ingress.Name, metav1.DeleteOptions{})
}

func (i *v1Ingress) GetIngressClass() string {
	if i.ingress.Spec.IngressClassName != nil {
		return *i.ingress.Spec.IngressClassName
	}
	return i.ingress.ObjectMeta.Annotations[ingressClassAnnotation]
}

func (i *v1Ingress) GetMetadata() metav1.ObjectMeta {
	return i.ingress.ObjectMeta
}

func (i *v1Ingress) GetLoadBalancersHosts() []string {
	var hostNames []string
	for _, hostName := range i.ingress.Status.LoadBalancer.Ingress {
		hostNames = append(hostNames, hostName.Hostname)
	}
	return hostNames
}

func listIngress(kubernetesCS kubernetes.Interface, clusterConfig *api.ClusterConfig) ([]Ingress, error) {
	useV1API, err := utils.IsMinVersion(api.Version1_19, clusterConfig.Metadata.Version)
	if err != nil {
		return nil, err
	}

	var ingressList []Ingress
	if useV1API {
		logger.Debug("using v1 networking API")
		ingresses, err := kubernetesCS.NetworkingV1().Ingresses(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		for _, item := range ingresses.Items {
			ingressList = append(ingressList, &v1Ingress{ingress: item})
		}
		return ingressList, nil
	}

	logger.Debug("using v1beta networking API")
	ingresses, err := kubernetesCS.NetworkingV1beta1().Ingresses(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, item := range ingresses.Items {
		ingressList = append(ingressList, &v1BetaIngress{ingress: item})
	}
	return ingressList, nil
}
