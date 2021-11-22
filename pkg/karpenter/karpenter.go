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

	karpenterHelmRepo         = "https://charts.karpenter.sh"
	karpenterHelmChartName    = "karpenter/karpenter"
	karpenterReleaseName      = "karpenter"
	controllerClusterName     = "controller.clusterName"
	controllerClusterEndpoint = "controller.clusterEndpoint"
	createServiceAccount      = "serviceAccount.create"
	addDefaultProvisioner     = "defaultProvisioner.create"
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
	logger.Info("adding Karpenter to cluster %s with cluster endpoint", k.ClusterName, k.ClusterEndpoint)
	if err := k.HelmInstaller.AddRepo(karpenterHelmRepo, karpenterReleaseName); err != nil {
		return fmt.Errorf("failed to add Karpenter repository: %w", err)
	}
	values := map[string]interface{}{
		createServiceAccount:      k.CreateServiceAccount,
		controllerClusterName:     k.ClusterName,
		controllerClusterEndpoint: k.ClusterEndpoint,
		addDefaultProvisioner:     k.AddDefaultProvisioner,
	}
	if err := k.HelmInstaller.InstallChart(ctx, karpenterReleaseName, karpenterHelmChartName, DefaultKarpenterNamespace, k.Version, values); err != nil {
		return fmt.Errorf("failed to install Karpenter chart: %w", err)
	}
	return nil
}
