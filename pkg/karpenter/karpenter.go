package karpenter

import (
	"context"
	"fmt"

	"github.com/kris-nova/logger"
	"helm.sh/helm/v3/pkg/registry"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/karpenter/providers"
	"github.com/weaveworks/eksctl/pkg/utils"
)

const (
	// DefaultNamespace default namespace for Karpenter
	DefaultNamespace = "karpenter"
	// DefaultServiceAccountName is the name of the service account which is needed for Karpenter
	DefaultServiceAccountName = "karpenter"

	aws                      = "aws"
	clusterEndpoint          = "clusterEndpoint"
	clusterName              = "clusterName"
	create                   = "create"
	defaultInstanceProfile   = "defaultInstanceProfile"
	helmChartName            = "oci://public.ecr.aws/karpenter/karpenter"
	releaseName              = "karpenter"
	serviceAccount           = "serviceAccount"
	serviceAccountAnnotation = "annotations"
	serviceAccountName       = "name"
	settings                 = "settings"
	interruptionQueueName    = "interruptionQueueName"
	interruptionQueue        = "interruptionQueue"
)

// Options contains values which Karpenter uses to configure the installation.
type Options struct {
	HelmInstaller providers.HelmInstaller
	Namespace     string
	ClusterConfig *api.ClusterConfig
}

// ChartInstaller defines a functionality to install Karpenter.
//
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

	serviceAccountMap := map[string]interface{}{
		create: api.IsEnabled(k.ClusterConfig.Karpenter.CreateServiceAccount),
		serviceAccountAnnotation: map[string]interface{}{
			api.AnnotationEKSRoleARN: serviceAccountRoleARN,
		},
		serviceAccountName: DefaultServiceAccountName,
	}

	settingsValues := map[string]interface{}{
		defaultInstanceProfile: instanceProfileName,
		clusterName:            k.ClusterConfig.Metadata.Name,
		clusterEndpoint:        k.ClusterConfig.Status.Endpoint,
	}

	// The Karpenter chart renamed the interruption queue Helm value when it
	// flattened `settings`: charts before v0.33.0 read
	// `settings.aws.interruptionQueueName`, while the flattened layout reads
	// `settings.interruptionQueue`. A flattened chart silently ignores the old
	// spelling, which leaves INTERRUPTION_QUEUE unset on the Karpenter pod and
	// disables spot interruption handling without reporting an error.
	version := k.ClusterConfig.Karpenter.Version
	compareVersions, err := utils.CompareVersions(version, "0.33.0")
	legacyChart := err == nil && compareVersions < 0

	// Only advertise the interruption queue when eksctl actually provisioned it.
	// pkg/cfn/builder creates the SQS queue -- and grants the controller role
	// sqs:ReceiveMessage on it -- only when withSpotInterruptionQueue is enabled,
	// so sending the name unconditionally would point Karpenter at a queue that
	// does not exist and that it has no permission to poll.
	queueEnabled := api.IsEnabled(k.ClusterConfig.Karpenter.WithSpotInterruptionQueue)

	if legacyChart {
		if queueEnabled {
			settingsValues[interruptionQueueName] = k.ClusterConfig.Metadata.Name
		}
		settingsValues = map[string]interface{}{
			aws: settingsValues,
		}
	} else if queueEnabled {
		settingsValues[interruptionQueue] = k.ClusterConfig.Metadata.Name
	}

	values := map[string]interface{}{
		clusterName:     k.ClusterConfig.Metadata.Name,
		clusterEndpoint: k.ClusterConfig.Status.Endpoint,
		aws: map[string]interface{}{
			defaultInstanceProfile: instanceProfileName,
		},
		settings:       settingsValues,
		serviceAccount: serviceAccountMap,
	}

	registryClient, err := registry.NewClient(
		registry.ClientOptEnableCache(true),
	)
	if err != nil {
		return fmt.Errorf("failed to create registry client: %w", err)
	}

	options := providers.InstallChartOpts{
		ChartName:       helmChartName,
		CreateNamespace: true,
		Namespace:       DefaultNamespace,
		ReleaseName:     releaseName,
		Values:          values,
		Version:         version,
		RegistryClient:  registryClient,
	}

	logger.Debug("the following chartOptions will be applied to the install: %+v", options)

	if err := k.HelmInstaller.InstallChart(ctx, options); err != nil {
		return fmt.Errorf("failed to install Karpenter chart: %w", err)
	}
	return nil
}
