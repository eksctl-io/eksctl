package flux

import (
	"context"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/flux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"
)

const allNamespaces = ""

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_flux_client.go . InstallerClient
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

	logger.Info("running pre-flight checks")
	if err := ti.fluxClient.PreFlight(); err != nil {
		return errors.Wrap(err, "running Flux pre-flight checks")
	}

	logger.Info("bootstrapping Flux v2 into cluster")
	if err := ti.fluxClient.Bootstrap(); err != nil {
		logger.Info("Flux v2 failed to install successfully. check configuration and re-run `eksctl enable flux`")
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
