package karpenter

import (
	"fmt"
	"time"

	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/karpenter"
	"github.com/weaveworks/eksctl/pkg/karpenter/providers/helm"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/eksctl/pkg/utils/waiters"

	"github.com/aws/aws-sdk-go/aws/request"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
)

// Installer contains all necessary dependencies for the Karpenter Install tasks and others.
type Installer struct {
	stackManager       manager.StackManager
	ctl                *eks.ClusterProvider
	cfg                *api.ClusterConfig
	wait               WaitFunc
	karpenterInstaller karpenter.InstallKarpenter
}

type WaitFunc func(name, msg string, acceptors []request.WaiterAcceptor, newRequest func() *request.Request, waitTimeout time.Duration, troubleshoot func(string) error) error

// NewInstaller creates a new Karpenter installer.
func NewInstaller(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, stackManager manager.StackManager) (*Installer, error) {
	helmInstaller, err := helm.NewInstaller(helm.Options{
		Namespace: karpenter.DefaultKarpenterNamespace,
	})
	if err != nil {
		return nil, err
	}
	karpenterInstaller := karpenter.NewKarpenterInstaller(karpenter.Options{
		HelmInstaller:         helmInstaller,
		Namespace:             karpenter.DefaultKarpenterNamespace,
		ClusterName:           cfg.Metadata.Name,
		AddDefaultProvisioner: api.IsEnabled(cfg.Karpenter.AddDefaultProvisioner),
		CreateServiceAccount:  api.IsEnabled(cfg.Karpenter.CreateServiceAccount),
		ClusterEndpoint:       cfg.Status.Endpoint,
		Version:               cfg.Karpenter.Version,
	})
	return &Installer{
		stackManager:       stackManager,
		ctl:                ctl,
		cfg:                cfg,
		wait:               waiters.Wait,
		karpenterInstaller: karpenterInstaller,
	}, nil
}

func doTasks(taskTree *tasks.TaskTree) error {
	logger.Info(taskTree.Describe())
	if errs := taskTree.DoAllSync(); len(errs) > 0 {
		logger.Info("%d error(s) occurred while installing Karpenter, you may wish to check your Cluster for further information", len(errs))
		for _, err := range errs {
			logger.Critical("%s\n", err.Error())
		}
		return fmt.Errorf("failed to install Karpenter on cluster")
	}
	return nil
}
