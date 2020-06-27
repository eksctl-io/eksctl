package flux

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	portforward "github.com/justinbarrick/go-k8s-portforward"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	portForwardingTimeout     = 120 * time.Second
	portForwardingRetryPeriod = 2 * time.Second
)

// PublicKey represents a public SSH key as it is returned by flux
type PublicKey struct {
	Key string `json:"key"`
}

func waitForHelmOpToStart(namespace string, timeout time.Duration, cs kubeclient.Interface) error {
	return waitForDeploymentToStart(cs, namespace, "helm-operator", timeout)
}

func waitForFluxToStart(namespace string, timeout time.Duration, cs kubeclient.Interface) error {
	return waitForDeploymentToStart(cs, namespace, "flux", timeout)
}

func getPublicKeyFromFlux(ctx context.Context, namespace string, timeout time.Duration, restConfig *rest.Config,
	cs kubeclient.Interface) (PublicKey, error) {
	var deployKey PublicKey
	try := func(rootURL string) error {
		fluxURL := rootURL + "api/flux/v6/identity.pub"
		req, reqErr := http.NewRequest("GET", fluxURL, nil)
		if reqErr != nil {
			return fmt.Errorf("failed to create request: %s", reqErr)
		}
		repoCtx, repoCtxCancel := context.WithTimeout(ctx, timeout)
		defer repoCtxCancel()
		req = req.WithContext(repoCtx)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to query Flux API: %s", err)
		}
		if resp.Body == nil {
			return fmt.Errorf("failed to fetch Flux deploy key from: %s", fluxURL)
		}
		defer resp.Body.Close()

		jsonErr := json.NewDecoder(resp.Body).Decode(&deployKey)
		if jsonErr != nil {
			return fmt.Errorf("failed to decode Flux API response: %s", jsonErr)
		}

		if deployKey.Key == "" {
			return fmt.Errorf("failed to fetch Flux deploy key from: %s", fluxURL)
		}
		return nil
	}
	err := portForward(namespace, "flux", 3030, "Flux", restConfig, cs, try)
	return deployKey, err
}

type tryFunc func(rootURL string) error

func waitForDeploymentToStart(k8sClientSet kubeclient.Interface, namespace string, name string, timeout time.Duration) error {
	watcher, err := k8sClientSet.AppsV1().Deployments(namespace).Watch(metav1.ListOptions{
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

func portForward(namespace string, nameLabelValue string, port int, name string,
	restConfig *rest.Config, cs kubeclient.Interface, try tryFunc) error {
	fluxSelector := metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      "name",
				Operator: metav1.LabelSelectorOpIn,
				Values:   []string{nameLabelValue},
			},
		},
	}
	portforwarder := portforward.PortForward{
		Labels:          fluxSelector,
		Config:          restConfig,
		Clientset:       cs,
		DestinationPort: port,
		Namespace:       namespace,
	}
	podDeadline := time.Now().Add(portForwardingTimeout)
	for ; time.Now().Before(podDeadline); time.Sleep(portForwardingRetryPeriod) {
		err := portforwarder.Start()
		if err == nil {
			defer portforwarder.Stop()
			break
		}
		if !strings.Contains(err.Error(), "Could not find running pod for selector") {
			logger.Warning("%s is not ready yet (%s), retrying ...", name, err)
		}
	}
	if time.Now().After(podDeadline) {
		return fmt.Errorf("timed out waiting for %s's pod to be created", name)
	}
	baseURL := fmt.Sprintf("http://127.0.0.1:%d/", portforwarder.ListenPort)
	// Make sure it's alive
	retryDeadline := time.Now().Add(30 * time.Second)
	for ; time.Now().Before(retryDeadline); time.Sleep(2 * time.Second) {
		err := try(baseURL)
		if err == nil {
			break
		}
		logger.Warning("%s is not ready yet (%s), retrying ...", name, err)
	}
	if time.Now().After(retryDeadline) {
		return fmt.Errorf("timed out waiting for %s to be operative", name)
	}
	return nil
}
