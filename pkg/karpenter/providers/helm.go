package providers

import (
	"bytes"
	"context"

	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/registry"
)

// URLGetter is an interface to support GET to the specified URL.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_helm_getter.go . URLGetter
type URLGetter interface {
	// Get file content by url string
	Get(url string, options ...getter.Option) (*bytes.Buffer, error)
}

// InstallChartOpts defines parameters for InstallChart.
type InstallChartOpts struct {
	ChartName       string
	CreateNamespace bool
	Namespace       string
	ReleaseName     string
	Values          map[string]interface{}
	Version         string
	RegistryClient  *registry.Client
}

// HelmInstaller deals with setting up Helm related resources.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_helm_installer.go . HelmInstaller
type HelmInstaller interface {
	// InstallChart takes a releaseName's name and a chart name and installs it. If namespace is not empty
	// it will install into that namespace and create the namespace. Version is required.
	InstallChart(ctx context.Context, opts InstallChartOpts) error
}
