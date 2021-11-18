package helm

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/kris-nova/logger"
	"gopkg.in/yaml.v1"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/weaveworks/eksctl/pkg/karpenter/providers"
)

// Installer implement the HelmInstaller interface.
type Installer struct {
	Settings *cli.EnvSettings
	Getters  getter.Providers
}

// NewInstaller creates a new Helm backed Installer for repo resources.
func NewInstaller() *Installer {
	settings := cli.New()
	return &Installer{
		Settings: settings,
		Getters:  getter.All(settings),
	}
}

var _ providers.HelmInstaller = &Installer{}

// AddRepo adds a repository to helm repositories.
func (i *Installer) AddRepo(repoURL string, release string) error {
	if err := os.MkdirAll(filepath.Dir(i.Settings.RegistryConfig), os.ModePerm); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to make cache folder: %w", err)
	}
	b, err := ioutil.ReadFile(i.Settings.RegistryConfig)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read registry file: %w", err)
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return fmt.Errorf("failed to marshall repo file: %w", err)
	}

	c := repo.Entry{
		Name: release,
		URL:  repoURL,
	}
	r, err := repo.NewChartRepository(&c, i.Getters)
	if err != nil {
		return fmt.Errorf("failed to create new chart repository: %w", err)
	}
	if _, err := r.DownloadIndexFile(); err != nil {
		return fmt.Errorf("failed to download index file: %w", err)
	}
	f.Update(&c)
	if err := f.WriteFile(i.Settings.RepositoryConfig, 0644); err != nil {
		return fmt.Errorf("failed to write out repository config file: %w", err)
	}
	return nil
}

// InstallChart takes a repo's name and a chart name and installs it. If namespace is not empty
// it will install into that namespace and create the namespace. Version is required.
func (i *Installer) InstallChart(releaseName string, chartName string, namespace string, version string, values map[string]interface{}) error {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(i.Settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), logger.Debug); err != nil {
		return fmt.Errorf("failed to initialize action config: %w", err)
	}
	client := action.NewInstall(actionConfig)
	client.Wait = true
	client.Namespace = namespace
	client.ReleaseName = releaseName
	client.Version = version
	client.CreateNamespace = true
	client.Timeout = 30 * time.Second

	chartPath, err := client.ChartPathOptions.LocateChart(chartName, i.Settings)
	if err != nil {
		return fmt.Errorf("failed to locate chart: %w", err)
	}

	// Check chart dependencies to make sure all are present in /charts
	ch, err := loader.Load(chartPath)
	if err != nil {
		return fmt.Errorf("failed to load chart: %w", err)
	}

	release, err := client.RunWithContext(context.Background(), ch, values)
	if err != nil {
		return fmt.Errorf("failed to install chart: %w", err)
	}
	logger.Debug("successfully installed helm chart: ", release.Name)
	return nil
}

func (i *Installer) UninstallChart(chart string) error {
	panic("implement me")
}

// RemoveRepo do we even need this?
func (i *Installer) RemoveRepo(repoName string) error {
	r, err := repo.LoadFile(i.Settings.RepositoryConfig)
	if os.IsNotExist(err) || len(r.Repositories) == 0 {
		return fmt.Errorf("no repositories configured")
	}
	if !r.Remove(repoName) {
		return fmt.Errorf("repo %s not found", repoName)
	}
	if err := r.WriteFile(i.Settings.RepositoryConfig, 0644); err != nil {
		return fmt.Errorf("failed to write out the repository config file: %w", err)
	}
	// TODO: should remove the cache
	return nil
}
