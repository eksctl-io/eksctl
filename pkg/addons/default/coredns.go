package defaultaddons

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/addons"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

const (
	// CoreDNS is the name of the coredns addon
	CoreDNS = "coredns"
	// KubeDNS is the name of the kube-dns addon
	KubeDNS = "kube-dns"
)

// UpdateCoreDNS will update the `coredns` add-on
func UpdateCoreDNS(rawClient kubernetes.RawClientInterface, region, controlPlaneVersion string, plan bool) (bool, error) {
	kubeDNSSevice, err := rawClient.ClientSet().CoreV1().Services(metav1.NamespaceSystem).Get(KubeDNS, metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("%q service was not found", KubeDNS)
			return false, nil
		}
		return false, errors.Wrapf(err, "getting %q service", KubeDNS)
	}

	_, err = rawClient.ClientSet().AppsV1().Deployments(metav1.NamespaceSystem).Get(CoreDNS, metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("%q was not found", CoreDNS)
			return false, nil
		}
		return false, errors.Wrapf(err, "getting %q", CoreDNS)
	}

	// if Deployment is present, go through our list of assets
	list, err := loadAssetCoreDNS(controlPlaneVersion)
	if err != nil {
		return false, err
	}

	for _, rawObj := range list.Items {
		resource, err := rawClient.NewRawResource(rawObj.Object)
		if err != nil {
			return false, err
		}
		switch resource.GVK.Kind {
		case "Deployment":
			if resource.Info.Name != "coredns" {
				continue
			}
			deployment, ok := resource.Info.Object.(*appsv1.Deployment)
			if !ok {
				return false, fmt.Errorf("expected type %T; got %T", &appsv1.Deployment{}, resource.Info.Object)
			}
			if err := addons.UseRegionalImage(&deployment.Spec.Template, region); err != nil {
				return false, err
			}
		case "Service":
			resource.Info.Object.(*corev1.Service).SetResourceVersion(kubeDNSSevice.GetResourceVersion())
			resource.Info.Object.(*corev1.Service).Spec.ClusterIP = kubeDNSSevice.Spec.ClusterIP
		}

		status, err := resource.CreateOrReplace(plan)
		if err != nil {
			return false, err
		}
		logger.Info(status)
	}

	if plan {
		logger.Critical("(plan) %q is not up-to-date", CoreDNS)
		return true, nil
	}

	logger.Info("%q is now up-to-date", CoreDNS)
	return false, nil
}

func loadAssetCoreDNS(controlPlaneVersion string) (*metav1.List, error) {
	if strings.HasPrefix(controlPlaneVersion, "1.10.") {
		return nil, errors.New("CoreDNS is not supported on Kubernetes 1.10")
	}

	for _, version := range api.SupportedVersions() {
		if strings.HasPrefix(controlPlaneVersion, version+".") {
			return LoadAsset(fmt.Sprintf("%s-%s", CoreDNS, version), "json")
		}
	}

	return nil, errors.New("unsupported Kubernetes version")
}
