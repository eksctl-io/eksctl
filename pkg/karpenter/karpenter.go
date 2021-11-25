package karpenter

import (
	"context"
	"fmt"

	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/karpenter/providers"
)

const (
	// DefaultKarpenterNamespace default namespace for Karpenter
	DefaultKarpenterNamespace = "karpenter"
	// DefaultKarpenterServiceAccountName is the name of the service account which is needed for Karpenter
	DefaultKarpenterServiceAccountName = "karpenter"

	karpenterHelmRepo      = "https://charts.karpenter.sh"
	karpenterHelmChartName = "karpenter/karpenter"
	karpenterReleaseName   = "karpenter"
	controller             = "controller"
	clusterName            = "clusterName"
	clusterEndpoint        = "clusterEndpoint"
	serviceAccount         = "serviceAccount"
	defaultProvisioner     = "defaultProvisioner"
	create                 = "create"
)

// Options contains values which Karpenter uses to configure the installation.
type Options struct {
	HelmInstaller         providers.HelmInstaller
	Namespace             string
	ClusterName           string
	AddDefaultProvisioner bool
	CreateServiceAccount  bool
	ClusterEndpoint       string
	Version               string
}

// InstallKarpenter defines a functionality to install Karpenter.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_karpenter_installer.go . InstallKarpenter
type InstallKarpenter interface {
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
	logger.Info("adding Karpenter to cluster %s with cluster", k.ClusterName)
	logger.Debug("cluster endpoint used by Karpenter: %s", k.ClusterEndpoint)
	if err := k.HelmInstaller.AddRepo(karpenterHelmRepo, karpenterReleaseName); err != nil {
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
			create: k.AddDefaultProvisioner,
		},
	}
	logger.Debug("the following values will be applied to the install: %+v", values)
	if err := k.HelmInstaller.InstallChart(ctx, karpenterReleaseName, karpenterHelmChartName, DefaultKarpenterNamespace, k.Version, values); err != nil {
		return fmt.Errorf("failed to install Karpenter chart: %w", err)
	}
	return nil
}
