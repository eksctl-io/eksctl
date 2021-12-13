package karpenter

import (
	"context"
	"fmt"

	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/karpenter/providers"
)

const (
	// DefaultNamespace default namespace for Karpenter
	DefaultNamespace = "karpenter"
	// DefaultServiceAccountName is the name of the service account which is needed for Karpenter
	DefaultServiceAccountName = "karpenter"

	clusterEndpoint    = "clusterEndpoint"
	clusterName        = "clusterName"
	controller         = "controller"
	create             = "create"
	defaultProvisioner = "defaultProvisioner"
	helmChartName      = "karpenter/karpenter"
	helmRepo           = "https://charts.karpenter.sh"
	releaseName        = "karpenter"
	serviceAccount     = "serviceAccount"
)

const (
	addDefaultProvisionerDefaultValue = false
)

// Options contains values which Karpenter uses to configure the installation.
type Options struct {
	HelmInstaller        providers.HelmInstaller
	Namespace            string
	ClusterName          string
	CreateServiceAccount bool
	ClusterEndpoint      string
	Version              string
	CreateNamespace      bool
}

// ChartInstaller defines a functionality to install Karpenter.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_chart_installer.go . ChartInstaller
type ChartInstaller interface {
	Install(ctx context.Context) error
}

// Installer implements the Karpenter installer functionality.
type Installer struct {
	Options
}

// NewKarpenterInstaller creates a new installer to configure and add Karpenter to a cluster.
func NewKarpenterInstaller(opts Options) *Installer {
	return &Installer{
		Options: opts,
	}
}

// Install adds Karpenter to a configured cluster in a separate CloudFormation stack.
func (k *Installer) Install(ctx context.Context) error {
	logger.Info("adding Karpenter to cluster %s", k.ClusterName)
	logger.Debug("cluster endpoint used by Karpenter: %s", k.ClusterEndpoint)
	if err := k.HelmInstaller.AddRepo(helmRepo, releaseName); err != nil {
		return fmt.Errorf("failed to add Karpenter repository: %w", err)
	}
	values := map[string]interface{}{
		controller: map[string]interface{}{
			clusterName:     k.ClusterName,
			clusterEndpoint: k.ClusterEndpoint,
		},
		serviceAccount: map[string]interface{}{
			create: k.CreateServiceAccount,
		},
		defaultProvisioner: map[string]interface{}{
			create: addDefaultProvisionerDefaultValue,
		},
	}
	logger.Debug("the following values will be applied to the install: %+v", values)
	if err := k.HelmInstaller.InstallChart(ctx, providers.InstallChartOpts{
		ChartName:       helmChartName,
		CreateNamespace: k.CreateNamespace,
		Namespace:       DefaultNamespace,
		ReleaseName:     releaseName,
		Values:          values,
		Version:         k.Version,
	}); err != nil {
		return fmt.Errorf("failed to install Karpenter chart: %w", err)
	}
	return nil
}
