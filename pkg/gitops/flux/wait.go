package flux

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	portforward "github.com/justinbarrick/go-k8s-portforward"
	"github.com/kris-nova/logger"
	fluxapi "github.com/fluxcd/flux/pkg/api/v6"
	transport "github.com/fluxcd/flux/pkg/http"
	"github.com/fluxcd/flux/pkg/http/client"
	"github.com/fluxcd/flux/pkg/ssh"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func waitForFluxToStart(ctx context.Context, namespace string, timeout time.Duration, restConfig *rest.Config,
	cs kubeclient.Interface) (ssh.PublicKey, error) {
	var fluxGitConfig fluxapi.GitConfig
	try := func(rootURL string) error {
		fluxURL := rootURL + "api/flux"
		fluxClient := client.New(http.DefaultClient, transport.NewAPIRouter(), fluxURL, client.Token(""))
		repoCtx, repoCtxCancel := context.WithTimeout(ctx, timeout)
		defer repoCtxCancel()
		var err error
		fluxGitConfig, err = fluxClient.GitRepoConfig(repoCtx, false)
		return err
	}
	err := waitForPodToStart(namespace, "flux", 3030, "Flux", restConfig, cs, try)
	return fluxGitConfig.PublicSSHKey, err
}

func waitForHelmOpToStart(ctx context.Context, namespace string, timeout time.Duration, restConfig *rest.Config,
	cs kubeclient.Interface) error {
	try := func(rootURL string) error {
		helmOpURL := rootURL + "healthz"
		req, err := http.NewRequest("GET", helmOpURL, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %s", err)
		}
		healthzCtx, healtzhCtxCancel := context.WithTimeout(ctx, timeout)
		defer healtzhCtxCancel()
		req = req.WithContext(healthzCtx)
		_, err = http.DefaultClient.Do(req)
		return err
	}
	return waitForPodToStart(namespace, "flux-helm-operator", 3030, "Helm Operator", restConfig, cs, try)
}

type tryFunc func(rootURL string) error

func waitForPodToStart(namespace string, nameLabelValue string, port int, name string,
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
	podDeadline := time.Now().Add(time.Second * 30)
	for ; time.Now().Before(podDeadline); time.Sleep(2 * time.Second) {
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
