package providers

import (
	"bytes"

	"helm.sh/helm/v3/pkg/getter"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Getter is an interface to support GET to the specified URL.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_helm_getter.go . Getter
type Getter interface {
	// Get file content by url string
	Get(url string, options ...getter.Option) (*bytes.Buffer, error)
}

// RESTClientGetter is an interface that the ConfigFlags describe to provide an easier way to mock for commands
// and eliminate the direct coupling to a struct type.  Users may wish to duplicate this type in their own packages
// as per the golang type overlapping.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_helm_rest_client_getter.go . RESTClientGetter
type RESTClientGetter interface {
	// ToRESTConfig returns restconfig
	ToRESTConfig() (*rest.Config, error)
	// ToDiscoveryClient returns discovery client
	ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error)
	// ToRESTMapper returns a restmapper
	ToRESTMapper() (meta.RESTMapper, error)
	// ToRawKubeConfigLoader return kubeconfig loader as-is
	ToRawKubeConfigLoader() clientcmd.ClientConfig
}

// HelmInstaller deals with setting up Helm related resources.
type HelmInstaller interface {
	// AddRepo adds a repository to helm repositories.
	AddRepo(repoURL string, release string) error
	// RemoveRepo removes a repository from the list of repositories.
	RemoveRepo(repoName string) error
	// InstallChart takes a releaseName's name and a chart name and installs it. If namespace is not empty
	// it will install into that namespace and create the namespace. Version is required.
	InstallChart(releaseName string, chartName string, namespace string, version string, values map[string]interface{}) error
	// UninstallChart removes an installed chart from the cluster.
	UninstallChart(chart string) error
}
