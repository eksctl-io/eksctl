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
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/weaveworks/eksctl/pkg/karpenter/providers"
)

// Options defines options for the Helm Installer.
type Options struct {
	Namespace        string
	RESTClientGetter genericclioptions.RESTClientGetter
}

// Installer implement the HelmInstaller interface.
type Installer struct {
	Settings     *cli.EnvSettings
	Getters      getter.Providers
	ActionConfig *action.Configuration
}

// NewInstaller creates a new Helm backed Installer for repo resources.
func NewInstaller(opts Options) (*Installer, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(opts.RESTClientGetter, opts.Namespace, "", logger.Debug); err != nil {
		return nil, fmt.Errorf("failed to initialize action config: %w", err)
	}
	return &Installer{
		Settings:     settings,
		Getters:      getter.All(settings),
		ActionConfig: actionConfig,
	}, nil
}

var _ providers.HelmInstaller = &Installer{}

// AddRepo adds a repository to helm repositories.
func (i *Installer) AddRepo(repoURL, release string) error {
	if err := os.MkdirAll(filepath.Dir(i.Settings.RegistryConfig), os.ModePerm); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to make cache folder: %w", err)
	}
	b, err := ioutil.ReadFile(i.Settings.RegistryConfig)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read registry file: %w", err)
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return fmt.Errorf("failed to marshal repo file: %w", err)
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
func (i *Installer) InstallChart(ctx context.Context, opts providers.InstallChartOpts) error {
	client := action.NewInstall(i.ActionConfig)
	client.Wait = true
	client.Namespace = opts.Namespace
	client.ReleaseName = opts.ReleaseName
	client.Version = opts.Version
	client.CreateNamespace = opts.CreateNamespace
	client.Timeout = 10 * time.Minute

	chartPath, err := client.ChartPathOptions.LocateChart(opts.ChartName, i.Settings)
	if err != nil {
		return fmt.Errorf("failed to locate chart: %w", err)
	}

	// possibly deal with chart dependencies, but for now, maybe we don't care.
	ch, err := loader.Load(chartPath)
	if err != nil {
		return fmt.Errorf("failed to load chart: %w", err)
	}

	release, err := client.RunWithContext(ctx, ch, opts.Values)
	if err != nil {
		return fmt.Errorf("failed to install chart: %w", err)
	}
	logger.Debug("successfully installed helm chart: %s", release.Name)
	return nil
}
