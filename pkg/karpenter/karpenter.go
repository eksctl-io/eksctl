package karpenter

import (
	"context"
	"fmt"

	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/karpenter/providers"
)

const (
	karpenterHelmRepo         = "https://charts.karpenter.sh"
	karpenterHelmChartName    = "karpenter/karpenter"
	karpenterReleaseName      = "karpenter"
	karpenterNamespace        = "karpenter"
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

// KarpenterInstaller defines an installer for Karpenter.
type KarpenterInstaller interface {
	InstallKarpenter(ctx context.Context) error
	UninstallKarpenter(ctx context.Context) error
}

// Installer implements the Karpenter installer using a HelmInstaller.
type Installer struct {
	Options
}

// NewKarpenterInstaller creates a new installer to configure and add Karpenter to a cluster.
func NewKarpenterInstaller(opts Options) *Installer {
	return &Installer{
		Options: opts,
	}
}

// InstallKarpenter adds Karpenter to a configured cluster in a separate CloudFormation stack.
func (k *Installer) InstallKarpenter(ctx context.Context) error {
	logger.Info("adding Karpenter to cluster %s with cluster endpoint", k.ClusterName, k.ClusterEndpoint)
	// Add the cloudformation stack and template creation here and the ask handling and all that jazz.
	// And lastly, when the CF stack returned, we add Karpenter on top using Helm.
	if err := k.HelmInstaller.AddRepo(karpenterHelmRepo, karpenterReleaseName); err != nil {
		return fmt.Errorf("failed to karpenter repo: %w", err)
	}
	values := map[string]interface{}{
		createServiceAccount:      k.CreateServiceAccount,
		controllerClusterName:     k.ClusterName,
		controllerClusterEndpoint: k.ClusterEndpoint,
		addDefaultProvisioner:     k.AddDefaultProvisioner,
	}
	if err := k.HelmInstaller.InstallChart(ctx, karpenterReleaseName, karpenterHelmChartName, karpenterNamespace, k.Version, values); err != nil {
		return fmt.Errorf("failed to install karpenter chart: %w", err)
	}
	return nil
}

func (k *Installer) UninstallKarpenter(ctx context.Context) error {
	return nil
}
