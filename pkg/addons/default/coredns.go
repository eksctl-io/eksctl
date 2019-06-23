package defaultaddons

import (
	"fmt"
	"strings"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

const (
	// CoreDNS is the name of the coredns addon
	CoreDNS = "coredns"
	// KubeDNS is the name of the kube-dns addon
	KubeDNS = "kube-dns"

	componentLabel = "eks.amazonaws.com/component"

	coreDNSImagePrefix = "602401143452.dkr.ecr."
	coreDNSImageSuffix = ".amazonaws.com/eks/coredns"
)

// InstallCoreDNS will install the `coredns` add-on in place of `kube-dns`
func InstallCoreDNS(rawClient kubernetes.RawClientInterface, region, controlPlaneVersion string, waitTimeout *time.Duration, plan bool) (bool, error) {
	kubeDNSSevice, err := rawClient.ClientSet().CoreV1().Services(metav1.NamespaceSystem).Get(KubeDNS, metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("%q service was not found", KubeDNS)
			return false, nil
		}
		return false, errors.Wrapf(err, "getting %q service", KubeDNS)
	}

	kubeDNSDeployment, err := rawClient.ClientSet().AppsV1().Deployments(metav1.NamespaceSystem).Get(KubeDNS, metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("%q deployment was not found", KubeDNS)
			return false, nil
		}
		return false, errors.Wrapf(err, "getting %q deployment", KubeDNS)
	}

	if v, ok := kubeDNSDeployment.Spec.Selector.MatchLabels[componentLabel]; !ok || v != KubeDNS {
		logger.Debug("adding %q label to %q", componentLabel, KubeDNS)
		kubeDNSDeployment.Spec.Selector.MatchLabels[componentLabel] = KubeDNS
		if !plan {
			if _, err := rawClient.ClientSet().AppsV1().Deployments(metav1.NamespaceSystem).Update(kubeDNSDeployment); err != nil {
				return false, errors.Wrapf(err, "patching %q", KubeDNS)
			}
		}
	}

	// if kube-dns is present, go ahead and try to replace it with coredns
	list, err := loadAssetCoreDNS(controlPlaneVersion)
	if err != nil {
		return false, err
	}

	listPodsOptions := metav1.ListOptions{}
	replicas := 0
	for _, rawObj := range list.Items {
		resource, err := rawClient.NewRawResource(rawObj)
		if err != nil {
			return false, err
		}
		switch resource.GVK.Kind {
		case "Deployment":
			coreDNSDeployemnt := resource.Info.Object.(*appsv1.Deployment)
			listPodsOptions.LabelSelector = labels.FormatLabels(coreDNSDeployemnt.Spec.Selector.MatchLabels)
			replicas = int(*coreDNSDeployemnt.Spec.Replicas)
			image := &coreDNSDeployemnt.Spec.Template.Spec.Containers[0].Image
			imageParts := strings.Split(*image, ":")

			if len(imageParts) != 2 {
				return false, fmt.Errorf("unexpected image format %q for %q", *image, KubeProxy)
			}

			if strings.HasPrefix(imageParts[0], coreDNSImagePrefix) &&
				strings.HasSuffix(imageParts[0], coreDNSImageSuffix) {
				*image = coreDNSImagePrefix + region + coreDNSImageSuffix + ":" + imageParts[1]
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

	if waitTimeout != nil && !plan {
		timer := time.After(*waitTimeout)
		timeout := false
		readyPods := sets.NewString()
		watcher, err := rawClient.ClientSet().CoreV1().Pods(metav1.NamespaceSystem).Watch(listPodsOptions)
		if err != nil {
			return false, errors.Wrapf(err, "creating %q pod watcher", CoreDNS)
		}

		podList, err := rawClient.ClientSet().CoreV1().Pods(metav1.NamespaceSystem).List(listPodsOptions)
		if err != nil {
			return false, errors.Wrapf(err, "listing %q pods", CoreDNS)
		}
		for _, pod := range podList.Items {
			if pod.Status.Phase == corev1.PodRunning {
				readyPods.Insert(pod.Name)
			}
		}

		logger.Info("waiting for %d of %q pods to become ready", replicas, CoreDNS)
		for !timeout && readyPods.Len() < replicas {
			select {
			case event := <-watcher.ResultChan():
				logger.Debug("event = %#v", event)
				if event.Object != nil && event.Type != watch.Deleted {
					if pod, ok := event.Object.(*corev1.Pod); ok {
						if pod.Status.Phase == corev1.PodRunning {
							readyPods.Insert(pod.Name)
							logger.Debug("pod %q is ready", pod.Name)
						} else {
							logger.Debug("pod %q seen, but not ready yet", pod.Name)
							logger.Debug("node = %#v", *pod)
						}
					}
				}
			case <-timer:
				timeout = true
			}
		}
		watcher.Stop()
		if timeout {
			return false, fmt.Errorf("timed out (after %s) waiting for %q pods to become ready", waitTimeout, CoreDNS)
		}
	}

	if plan {
		logger.Info("(plan) would have waited for %q pods to become ready and then delete %q", CoreDNS, KubeDNS)
		return true, nil
	}

	if err := rawClient.ClientSet().AppsV1().Deployments(metav1.NamespaceSystem).Delete(KubeDNS, &metav1.DeleteOptions{}); err != nil {
		return false, errors.Wrapf(err, "deleting %q", KubeDNS)
	}

	logger.Info("deleted %q", KubeDNS)

	logger.Info("%q is now up-to-date", CoreDNS)
	return false, nil
}

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
		resource, err := rawClient.NewRawResource(rawObj)
		if err != nil {
			return false, err
		}
		switch resource.GVK.Kind {
		case "Deployment":
			image := &resource.Info.Object.(*appsv1.Deployment).Spec.Template.Spec.Containers[0].Image
			imageParts := strings.Split(*image, ":")

			if len(imageParts) != 2 {
				return false, fmt.Errorf("unexpected image format %q for %q", *image, KubeProxy)
			}

			if strings.HasPrefix(imageParts[0], coreDNSImagePrefix) &&
				strings.HasSuffix(imageParts[0], coreDNSImageSuffix) {
				*image = coreDNSImagePrefix + region + coreDNSImageSuffix + ":" + imageParts[1]
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
	assetName := CoreDNS
	if strings.HasPrefix(controlPlaneVersion, "1.10.") {
		return nil, fmt.Errorf("CoreDNS is not supported on Kubernetes 1.10")
	}
	if strings.HasPrefix(controlPlaneVersion, "1.11.") {
		assetName += "-1.11"
	}
	if strings.HasPrefix(controlPlaneVersion, "1.12.") {
		assetName += "-1.12"
	}
	if strings.HasPrefix(controlPlaneVersion, "1.13.") {
		assetName += "-1.13"
	}
	return LoadAsset(assetName, "json")
}
