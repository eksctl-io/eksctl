// FLUX V1 DEPRECATION NOTICE. https://github.com/weaveworks/eksctl/issues/2963
package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	kubeclient "k8s.io/client-go/kubernetes"
)

func waitForHelmOpToStart(namespace string, timeout time.Duration, cs kubeclient.Interface) error {
	return waitForDeploymentToStart(cs, namespace, "helm-operator", timeout)
}

func waitForFluxToStart(namespace string, timeout time.Duration, cs kubeclient.Interface) error {
	return waitForDeploymentToStart(cs, namespace, "flux", timeout)
}

func waitForDeploymentToStart(k8sClientSet kubeclient.Interface, namespace string, name string, timeout time.Duration) error {
	watcher, err := k8sClientSet.AppsV1().Deployments(namespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return err
	}

	defer watcher.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return errors.Wrapf(err, "failed waiting for pod %q", name)
			}
			switch event.Type {
			case watch.Added, watch.Modified:
				deployment, ok := event.Object.(*v1.Deployment)
				if !ok {
					return errors.Errorf("expected event type to be %T; got %T", &v1.Deployment{}, event.Object)
				}
				if deployment.Status.ReadyReplicas == deployment.Status.Replicas {
					return nil
				}
			}
		case <-timer.C:
			return fmt.Errorf("timed out (after %v) waiting for deployment %q", timeout, name)
		}
	}
}
