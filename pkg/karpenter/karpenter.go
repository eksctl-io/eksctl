package karpenter

import (
	"context"
	"fmt"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/karpenter/providers"
)

const (
	// DefaultNamespace default namespace for Karpenter
	DefaultNamespace = "karpenter"
	// DefaultServiceAccountName is the name of the service account which is needed for Karpenter
	DefaultServiceAccountName = "karpenter"

	aws                      = "aws"
	clusterEndpoint          = "clusterEndpoint"
	clusterName              = "clusterName"
	controller               = "controller"
	create                   = "create"
	defaultInstanceProfile   = "defaultInstanceProfile"
	helmChartName            = "karpenter/karpenter"
	helmRepo                 = "https://charts.karpenter.sh"
	releaseName              = "karpenter"
	serviceAccount           = "serviceAccount"
	serviceAccountAnnotation = "annotations"
)

// Options contains values which Karpenter uses to configure the installation.
type Options struct {
	HelmInstaller providers.HelmInstaller
	Namespace     string
	ClusterConfig *api.ClusterConfig
}

// ChartInstaller defines a functionality to install Karpenter.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_chart_installer.go . ChartInstaller
type ChartInstaller interface {
	Install(ctx context.Context, serviceAccountRoleARN string, instanceProfileName string) error
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
func (k *Installer) Install(ctx context.Context, serviceAccountRoleARN string, instanceProfileName string) error {
	logger.Info("adding Karpenter to cluster %s", k.ClusterConfig.Metadata.Name)
	logger.Debug("cluster endpoint used by Karpenter: %s", k.ClusterConfig.Status.Endpoint)
	if err := k.HelmInstaller.AddRepo(helmRepo, releaseName); err != nil {
		return fmt.Errorf("failed to add Karpenter repository: %w", err)
	}
	serviceAccountMap := map[string]interface{}{
		create: api.IsEnabled(k.ClusterConfig.Karpenter.CreateServiceAccount),
	}
	if serviceAccountRoleARN != "" {
		serviceAccountMap[serviceAccountAnnotation] = map[string]interface{}{
			api.AnnotationEKSRoleARN: serviceAccountRoleARN,
		}
	}
	values := map[string]interface{}{
		controller: map[string]interface{}{
			clusterName:     k.ClusterConfig.Metadata.Name,
			clusterEndpoint: k.ClusterConfig.Status.Endpoint,
		},
		aws: map[string]interface{}{
			defaultInstanceProfile: instanceProfileName,
		},
		serviceAccount: serviceAccountMap,
	}

	logger.Debug("the following values will be applied to the install: %+v", values)
	if err := k.HelmInstaller.InstallChart(ctx, providers.InstallChartOpts{
		ChartName:       helmChartName,
		CreateNamespace: len(k.ClusterConfig.FargateProfiles) == 0,
		Namespace:       DefaultNamespace,
		ReleaseName:     releaseName,
		Values:          values,
		Version:         k.ClusterConfig.Karpenter.Version,
	}); err != nil {
		return fmt.Errorf("failed to install Karpenter chart: %w", err)
	}
	return nil
}
