package flux

import (
	"context"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/flux"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"
)

const allNamespaces = ""

var DefaultFluxComponents = []string{"helm-controller", "kustomize-controller", "notification-controller", "source-controller"}

//go:generate counterfeiter -o fakes/fake_flux_client.go . InstallerClient
type InstallerClient interface {
	PreFlight() error
	Bootstrap() error
}

type Installer struct {
	opts       *api.Flux
	kubeClient kubeclient.Interface
	fluxClient InstallerClient
}

func New(k8sClientSet kubeclient.Interface, opts *api.GitOps) (*Installer, error) {
	if opts.Flux == nil {
		return nil, errors.New("expected gitops.flux in cluster configuration but found nil")
	}

	fluxClient, err := flux.NewClient(opts.Flux)
	if err != nil {
		return nil, err
	}

	installer := &Installer{
		opts:       opts.Flux,
		kubeClient: k8sClientSet,
		fluxClient: fluxClient,
	}

	return installer, nil
}

func (ti *Installer) Run() error {
	// TODO for now if we discover V1 components we just abort
	// in future we may want to handle some migration magic
	logger.Info("ensuring v1 repo components not installed")
	v1Namespace, err := ti.checkV1()
	if err != nil {
		return errors.Wrap(err, "checking for flux v1 components")
	}

	if v1Namespace != "" {
		logger.Warning("flux v1 components already installed in namespace %q. auto-migration not yet supported by eksctl. skipping installation", v1Namespace)
		return nil
	}

	logger.Info("checking whether Flux v2 components already installed")
	alreadyInstalled := ti.checkInstallation()
	if alreadyInstalled {
		logger.Warning("found existing Flux v2 components in namespace %q. skipping installation", ti.opts.Namespace)
		return nil
	}

	logger.Info("running pre-flight checks")
	if err := ti.fluxClient.PreFlight(); err != nil {
		return errors.Wrap(err, "running Flux pre-flight checks")
	}

	logger.Info("bootstrapping Flux v2 into cluster")
	if err := ti.fluxClient.Bootstrap(); err != nil {
		return errors.Wrap(err, "running Flux Bootstrap")
	}

	logger.Success("Flux v2 installed successfully")
	logger.Success("see https://toolkit.fluxcd.io/ for usage instructions")

	return nil
}

func (ti *Installer) checkV1() (string, error) {
	deployments, err := ti.kubeClient.AppsV1().Deployments(allNamespaces).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	for _, d := range deployments.Items {
		if d.ObjectMeta.Name == "flux" {
			return d.ObjectMeta.Namespace, nil
		}
	}

	return "", nil
}

func (ti *Installer) checkInstallation() bool {
	var count int
	// TODO: this is the default set, we should maybe have an option to add whatever
	for _, c := range DefaultFluxComponents {
		// TODO: this checks for components in the currently configured namespace,
		// but it is possible that a previous bootstrap installed them elsewhere
		if found := ti.checkComponent(ti.opts.Namespace, c); found {
			count++
		}
	}

	return count == len(DefaultFluxComponents)
}

func (ti *Installer) checkComponent(namespace, component string) bool {
	_, err := ti.kubeClient.AppsV1().Deployments(namespace).Get(context.Background(), component, metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("%s deployment was not found", component)
			return false
		}
		logger.Warning("error while looking for %s deployment: %s", component, err)
		return false
	}

	logger.Info("component %s found", component)
	return true
}
