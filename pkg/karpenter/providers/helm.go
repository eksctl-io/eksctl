package providers

import (
	"bytes"
	"context"

	"helm.sh/helm/v3/pkg/getter"
)

// Getter is an interface to support GET to the specified URL.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_helm_getter.go . Getter
type Getter interface {
	// Get file content by url string
	Get(url string, options ...getter.Option) (*bytes.Buffer, error)
}

// HelmInstaller deals with setting up Helm related resources.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_helm_installer.go . HelmInstaller
type HelmInstaller interface {
	// AddRepo adds a repository to helm repositories.
	AddRepo(repoURL string, release string) error
	// InstallChart takes a releaseName's name and a chart name and installs it. If namespace is not empty
	// it will install into that namespace and create the namespace. Version is required.
	InstallChart(ctx context.Context, releaseName, chartName, namespace, version string, values map[string]interface{}) error
}
