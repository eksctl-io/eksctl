package defaultaddons

import (
	"fmt"
	"strings"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
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
func InstallCoreDNS(rawClient kubernetes.RawClientInterface, region string, waitTimeout *time.Duration) error {
	kubeDNSSevice, err := rawClient.ClientSet().CoreV1().Services(metav1.NamespaceSystem).Get(KubeDNS, metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("%q service was not found", KubeDNS)
			return nil
		}
		return errors.Wrapf(err, "getting %q service", KubeDNS)
	}

	kubeDNSDeployment, err := rawClient.ClientSet().AppsV1().Deployments(metav1.NamespaceSystem).Get(KubeDNS, metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("%q deployment was not found", KubeDNS)
			return nil
		}
		return errors.Wrapf(err, "getting %q deployment", KubeDNS)
	}

	if v, ok := kubeDNSDeployment.Spec.Selector.MatchLabels[componentLabel]; !ok || v != KubeDNS {
		logger.Debug("adding %q label to %q", componentLabel, KubeDNS)
		kubeDNSDeployment.Spec.Selector.MatchLabels[componentLabel] = KubeDNS
		if _, err := rawClient.ClientSet().AppsV1().Deployments(metav1.NamespaceSystem).Update(kubeDNSDeployment); err != nil {
			return errors.Wrapf(err, "patching %q", KubeDNS)
		}
	}

	// if kube-dns is present, go ahead and try to replace it with coredns
	list, err := LoadAsset(CoreDNS, "yaml")
	if err != nil {
		return err
	}

	listPodsOptions := metav1.ListOptions{}
	replicas := 0
	for _, rawObj := range list.Items {
		resource, err := rawClient.NewRawResource(rawObj)
		if err != nil {
			return err
		}
		switch resource.GVK.Kind {
		case "Deployment":
			coreDNSDeployemnt := resource.Info.Object.(*extensionsv1beta1.Deployment)
			listPodsOptions.LabelSelector = labels.FormatLabels(coreDNSDeployemnt.Spec.Selector.MatchLabels)
			replicas = int(*coreDNSDeployemnt.Spec.Replicas)
			image := &coreDNSDeployemnt.Spec.Template.Spec.Containers[0].Image
			imageParts := strings.Split(*image, ":")

			if len(imageParts) != 2 {
				return fmt.Errorf("unexpected image format %q for %q", *image, KubeProxy)
			}

			if strings.HasPrefix(imageParts[0], coreDNSImagePrefix) &&
				strings.HasSuffix(imageParts[0], coreDNSImageSuffix) {
				*image = coreDNSImagePrefix + region + coreDNSImageSuffix + ":" + imageParts[1]
			}
		case "Service":
			resource.Info.Object.(*corev1.Service).Spec.ClusterIP = kubeDNSSevice.Spec.ClusterIP
		}

		status, err := resource.CreateOrReplace()
		if err != nil {
			return err
		}
		logger.Info(status)
	}

	if waitTimeout != nil {
		timer := time.After(*waitTimeout)
		timeout := false
		readyPods := sets.NewString()
		watcher, err := rawClient.ClientSet().CoreV1().Pods(metav1.NamespaceSystem).Watch(listPodsOptions)
		if err != nil {
			return errors.Wrapf(err, "creating %q pod watcher", CoreDNS)
		}

		podList, err := rawClient.ClientSet().CoreV1().Pods(metav1.NamespaceSystem).List(listPodsOptions)
		if err != nil {
			return errors.Wrapf(err, "listing %q pods", CoreDNS)
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
			return fmt.Errorf("timed out (after %s) waiting for %q pods to become ready", waitTimeout, CoreDNS)
		}
	}

	if err := rawClient.ClientSet().AppsV1().Deployments(metav1.NamespaceSystem).Delete(KubeDNS, &metav1.DeleteOptions{}); err != nil {
		return errors.Wrapf(err, "deleting %q", KubeDNS)
	}

	logger.Info("deleted %q", KubeDNS)

	logger.Info("%q is now up-to-date", CoreDNS)
	return nil
}
