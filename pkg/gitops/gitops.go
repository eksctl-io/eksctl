package gitops

import (
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/weaveworks/eksctl/pkg/actions/flux"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// DefaultPodReadyTimeout is the time it will wait for Flux and Helm Operator to become ready
const DefaultPodReadyTimeout = 5 * time.Minute

type FluxInstaller interface {
	Run() error
}

// Setup sets up gitops in a repository for a cluster.
func Setup(kubeconfigPath string, k8sRestConfig *rest.Config, k8sClientSet kubeclient.Interface, cfg *api.ClusterConfig, timeout time.Duration) error {
	installer, err := flux.New(k8sClientSet, cfg.GitOps)
	logger.Info("gitops configuration detected, setting installer to Flux v2")
	if err != nil {
		return errors.Wrapf(err, "could not initialise Flux installer")
	}

	return installer.Run()
}
