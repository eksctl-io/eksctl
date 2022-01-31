package defaultaddons

import (
	"context"
	"embed"
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
	"github.com/weaveworks/eksctl/pkg/fargate/coredns"

	// For go:embed
	_ "embed"
)

const (
	// CoreDNS is the name of the coredns addon
	CoreDNS = "coredns"
	// KubeDNS is the name of the kube-dns addon
	KubeDNS = "kube-dns"
)

//go:embed assets/coredns*.json
var coreDNSDir embed.FS

func IsCoreDNSUpToDate(input AddonInput) (bool, error) {
	kubeDNSDeployment, err := input.RawClient.ClientSet().AppsV1().Deployments(metav1.NamespaceSystem).Get(context.TODO(), CoreDNS, metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("%q was not found", CoreDNS)
			return true, nil
		}
		return false, errors.Wrapf(err, "getting %q", CoreDNS)
	}

	// if Deployment is present, go through our list of assets
	list, err := loadAssetCoreDNS(input.ControlPlaneVersion)
	if err != nil {
		return false, err
	}

	for _, rawObj := range list.Items {
		resource, err := input.RawClient.NewRawResource(rawObj.Object)
		if err != nil {
			return false, err
		}
		if resource.GVK.Kind != "Deployment" {
			continue
		}
		if resource.Info.Name != "coredns" {
			continue
		}
		deployment, ok := resource.Info.Object.(*appsv1.Deployment)
		if !ok {
			return false, fmt.Errorf("expected type %T; got %T", &appsv1.Deployment{}, resource.Info.Object)
		}
		if err := addons.UseRegionalImage(&deployment.Spec.Template, input.Region); err != nil {
			return false, err
		}
		if computeType, ok := kubeDNSDeployment.Spec.Template.Annotations[coredns.ComputeTypeAnnotationKey]; ok {
			deployment.Spec.Template.Annotations[coredns.ComputeTypeAnnotationKey] = computeType
		}
		tagMismatch, err := addons.ImageTagsDiffer(
			deployment.Spec.Template.Spec.Containers[0].Image,
			kubeDNSDeployment.Spec.Template.Spec.Containers[0].Image,
		)
		if err != nil {
			return false, err
		}
		return !tagMismatch, err
	}
	return true, nil
}

// UpdateCoreDNS will update the `coredns` add-on and returns true
// if an update is available
func UpdateCoreDNS(input AddonInput, plan bool) (bool, error) {
	kubeDNSSevice, err := input.RawClient.ClientSet().CoreV1().Services(metav1.NamespaceSystem).Get(context.TODO(), KubeDNS, metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("%q service was not found", KubeDNS)
			return false, nil
		}
		return false, errors.Wrapf(err, "getting %q service", KubeDNS)
	}

	kubeDNSDeployment, err := input.RawClient.ClientSet().AppsV1().Deployments(metav1.NamespaceSystem).Get(context.TODO(), CoreDNS, metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("%q was not found", CoreDNS)
			return false, nil
		}
		return false, errors.Wrapf(err, "getting %q", CoreDNS)
	}

	// if Deployment is present, go through our list of assets
	list, err := loadAssetCoreDNS(input.ControlPlaneVersion)
	if err != nil {
		return false, err
	}

	tagMismatch := true
	for _, rawObj := range list.Items {
		resource, err := input.RawClient.NewRawResource(rawObj.Object)
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
			template := &deployment.Spec.Template
			if err := addons.UseRegionalImage(template, input.Region); err != nil {
				return false, err
			}
			if computeType, ok := kubeDNSDeployment.Spec.Template.Annotations[coredns.ComputeTypeAnnotationKey]; ok {
				if template.Annotations == nil {
					template.Annotations = make(map[string]string)
				}
				template.Annotations[coredns.ComputeTypeAnnotationKey] = computeType
			}
			tagMismatch, err = addons.ImageTagsDiffer(
				template.Spec.Containers[0].Image,
				kubeDNSDeployment.Spec.Template.Spec.Containers[0].Image,
			)
			if err != nil {
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
		if tagMismatch {
			logger.Critical("(plan) %q is not up-to-date", CoreDNS)
			return true, nil
		}
		logger.Info("(plan) %q is already up-to-date", CoreDNS)
		return false, nil
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
			manifest, err := coreDNSDir.ReadFile(fmt.Sprintf("assets/%s-%s.json", CoreDNS, version))
			if err != nil {
				return nil, err
			}
			return newList(manifest)
		}
	}
	return nil, errors.New("unsupported Kubernetes version")
}
